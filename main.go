package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/term"
	"os/exec"
)

type FieldType int

const (
	NoField FieldType = iota
	TextField
	NumberField
	FloatField
	DateField
	DateTimeField
	FileField
	URLField
	PasswordField
	SelectField
	BooleanField
	ListField
	HiddenField
	SkipField
	ButtonField
	CheckboxField
	ColorField
	EmailField
	ImageField
	MonthField
	RadioField
	RangeField
	ResetField
	TelField
	TimeField
	WeekField
	ReadonlyField
)

type FieldInfo struct {
	Type     FieldType
	Label    string
	Options  []string
	Note     string
	Pattern  string // regex pattern for validation
	Required bool
	Readonly bool
	Default  string
	Action   string // shell script to run
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <form-file> [-o output-file]\n", os.Args[0])
		os.Exit(1)
	}

	formFile := os.Args[1]
	outputFile := ""
	// Parse -o param if present
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "-o" && i+1 < len(os.Args) {
			outputFile = os.Args[i+1]
			break
		}
	}
	if outputFile == "" {
		outputFile = determineOutputFile(formFile)
	}

	if _, err := os.Stat(formFile); os.IsNotExist(err) {
		fmt.Printf("Error: Form file '%s' not found.\n", formFile)
		os.Exit(1)
	}

	// Check if output file exists and ask for confirmation
	if _, err := os.Stat(outputFile); err == nil {
		fmt.Printf("Output file '%s' already exists. Overwrite? [y/N]: ", outputFile)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			fmt.Println("Operation cancelled.")
			os.Exit(0)
		}
	}

	// Try to create/truncate the output file before wizard begins
	f, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error: Cannot create or write to output file '%s': %v\n", outputFile, err)
		os.Exit(1)
	}
	f.Close()

	// Check if file has any DSL directives
	hasDSL := checkForDSL(formFile)

	fmt.Printf("Interactive configuration for %s\n", filepath.Base(formFile))
	fmt.Println(strings.Repeat("-", 50))

	outputLines, err := processFormFile(formFile, outputFile, hasDSL)
	if err != nil {
		fmt.Printf("Error processing form file: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputFile, []byte(strings.Join(outputLines, "\n")+"\n"), 0644)
	if err != nil {
		fmt.Printf("Error writing to output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Configuration saved to %s\n", outputFile)
}

func determineOutputFile(formFile string) string {
	dir := filepath.Dir(formFile)
	base := filepath.Base(formFile)

	// Strip known extensions
	knownExtensions := []string{".xenv", ".template", ".example"}
	for _, ext := range knownExtensions {
		if strings.HasSuffix(base, ext) {
			base = strings.TrimSuffix(base, ext)
			break
		}
	}

	// Add dot prefix if not already present
	if !strings.HasPrefix(base, ".") {
		base = "." + base
	}

	return filepath.Join(dir, base)
}

func checkForDSL(formFile string) bool {
	file, err := os.Open(formFile)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "@text") ||
			strings.Contains(line, "@number") ||
			strings.Contains(line, "@float") ||
			strings.Contains(line, "@date") ||
			strings.Contains(line, "@file") ||
			strings.Contains(line, "@url") ||
			strings.Contains(line, "@password") ||
			strings.Contains(line, "@select") ||
			strings.Contains(line, "@boolean") ||
			strings.Contains(line, "@list") ||
			strings.Contains(line, "@hidden") {
			return true
		}
	}
	return false
}

func processFormFile(formFile, outputFile string, hasDSL bool) ([]string, error) {
	file, err := os.Open(formFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	keyValueRegex := regexp.MustCompile(`^([^=]+)=(.*)$`)
	var outputLines []string
	var currentFieldInfo *FieldInfo
	skipMode := false
	hiddenNext := false

	envVars := make(map[string]string)
	// Preload all key-value pairs for env injection
	file.Seek(0, 0)
	scannerEnv := bufio.NewScanner(file)
	for scannerEnv.Scan() {
		if matches := regexp.MustCompile(`^([^=]+)=(.*)$`).FindStringSubmatch(scannerEnv.Text()); len(matches) > 0 {
			envVars[matches[1]] = matches[2]
		}
	}
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			fieldInfo := parseFieldDirective(line)
			if fieldInfo != nil {
				switch fieldInfo.Type {
				case SkipField:
					skipMode = true
				case HiddenField:
					hiddenNext = true
				default:
					currentFieldInfo = fieldInfo
				}
			}
			outputLines = append(outputLines, line)
			continue
		}

		if matches := keyValueRegex.FindStringSubmatch(line); len(matches) > 0 {
			key := matches[1]
			defaultValue := matches[2]
			existingValue := getExistingValue(key, defaultValue, outputFile)

			if skipMode {
				outputLines = append(outputLines, key+"="+existingValue)
			} else if hiddenNext {
				outputLines = append(outputLines, key+"="+existingValue)
				hiddenNext = false
			} else if currentFieldInfo != nil && currentFieldInfo.Type != NoField {
				if currentFieldInfo.Readonly {
					fmt.Printf("%s (readonly): %s\n", currentFieldInfo.Label, existingValue)
					outputLines = append(outputLines, key+"="+existingValue)
					currentFieldInfo = nil
					continue
				}
				if currentFieldInfo.Action != "" {
					fmt.Printf("This action will run the following shell script:\n")
					fmt.Print("Do you want to run this script? [y/N/v (view)]: ")
					reader := bufio.NewReader(os.Stdin)
					resp, _ := reader.ReadString('\n')
					resp = strings.TrimSpace(strings.ToLower(resp))
					if resp == "view" || resp == "v" {
						fmt.Println("Script content:")
						fmt.Println(currentFieldInfo.Action)
						fmt.Print("Run this script now? [y/N]: ")
						resp, _ = reader.ReadString('\n')
						resp = strings.TrimSpace(strings.ToLower(resp))
					}
					if resp == "y" || resp == "yes" {
						// Prepare environment with all vars from previous answers (outputLines) and file
						env := os.Environ()
						// Add previous answers (outputLines) first, so they override file values
						for _, l := range outputLines {
							if kv := regexp.MustCompile(`^([^=]+)=(.*)$`).FindStringSubmatch(l); len(kv) > 0 {
								env = append(env, fmt.Sprintf("%s=%s", kv[1], kv[2]))
							}
						}
						// Add file values (envVars) if not already present
						for k, v := range envVars {
							found := false
							for _, e := range env {
								if strings.HasPrefix(e, k+"=") {
									found = true
									break
								}
							}
							if !found {
								env = append(env, fmt.Sprintf("%s=%s", k, v))
							}
						}
						cmd := exec.Command("sh", "-c", currentFieldInfo.Action)
						cmd.Env = env
						out, err := cmd.Output()
						if err != nil {
							fmt.Printf("Error running action: %v\n", err)
						}
						val := strings.TrimSpace(string(out))
						outputLines = append(outputLines, key+"="+val)
						envVars[key] = val
						currentFieldInfo = nil
						continue
					} else {
						fmt.Println("Skipped running the script.")
					}
				}
				label := currentFieldInfo.Label
				if label == "" {
					label = key
				}
				val, err := promptForValue(currentFieldInfo, label, existingValue)
				if err != nil {
					return nil, err
				}
				outputLines = append(outputLines, key+"="+val)
				envVars[key] = val
				currentFieldInfo = nil
			} else if !hasDSL || key == "PLAIN_TEXT_VAR" {
				val, err := promptInput(key, existingValue)
				if err != nil {
					return nil, err
				}
				outputLines = append(outputLines, key+"="+val)
				envVars[key] = val
			} else {
				outputLines = append(outputLines, key+"="+existingValue)
				envVars[key] = existingValue
			}
		} else {
			outputLines = append(outputLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return outputLines, nil
}

func parseFieldDirective(line string) *FieldInfo {
	getLabel := func(line string) string {
		labelRegex := regexp.MustCompile(`label="([^"]*)"`)
		if m := labelRegex.FindStringSubmatch(line); len(m) > 1 {
			return m[1]
		}
		return ""
	}

	getNote := func(line string) string {
		noteRegex := regexp.MustCompile(`note="([^"]*)"`)
		if m := noteRegex.FindStringSubmatch(line); len(m) > 1 {
			return m[1]
		}
		return ""
	}

	getPattern := func(line string) string {
		patternRegex := regexp.MustCompile(`pattern="([^"]*)"`)
		if m := patternRegex.FindStringSubmatch(line); len(m) > 1 {
			return m[1]
		}
		return ""
	}

	getRequired := func(line string) bool {
		// Matches required or required=true (case-insensitive)
		requiredRegex := regexp.MustCompile(`(?i)\brequired\b(?:\s*=\s*"?true"?)?`)
		return requiredRegex.MatchString(line)
	}

	getReadonly := func(line string) bool {
		readonlyRegex := regexp.MustCompile(`(?i)\breadonly\b(?:\s*=\s*"?true"?)?`)
		return readonlyRegex.MatchString(line)
	}

	getDefault := func(line string) string {
		defaultRegex := regexp.MustCompile(`default="([^"]*)"`)
		if m := defaultRegex.FindStringSubmatch(line); len(m) > 1 {
			return m[1]
		}
		return ""
	}

	getAction := func(line string) string {
		re := regexp.MustCompile(`action=(?:"([^"]*)"|` + "`" + `([^` + "`" + `]*)` + "`" + `)`)
		m := re.FindStringSubmatch(line)
		if len(m) == 3 {
			if m[1] != "" {
				return m[1]
			}
			return m[2]
		}
		return ""
	}

	switch {
	case strings.Contains(line, "@hidden"):
		return &FieldInfo{Type: HiddenField, Note: getNote(line), Action: getAction(line)}
	case strings.Contains(line, "@skip"):
		return &FieldInfo{Type: SkipField, Note: getNote(line), Action: getAction(line)}
	case strings.Contains(line, "@readonly"):
		return &FieldInfo{Type: ReadonlyField, Label: getLabel(line), Note: getNote(line), Readonly: true, Action: getAction(line)}
	case strings.Contains(line, "@button"):
		return &FieldInfo{Type: ButtonField, Label: getLabel(line), Note: getNote(line), Action: getAction(line)}
	case strings.Contains(line, "@checkbox"):
		return &FieldInfo{Type: CheckboxField, Label: getLabel(line), Note: getNote(line), Action: getAction(line)}
	case strings.Contains(line, "@color"):
		return &FieldInfo{Type: ColorField, Label: getLabel(line), Note: getNote(line), Action: getAction(line)}
	case strings.Contains(line, "@date"):
		return &FieldInfo{Type: DateField, Label: getLabel(line), Note: getNote(line), Pattern: getPattern(line), Required: getRequired(line), Readonly: getReadonly(line), Action: getAction(line)}
	case strings.Contains(line, "@datetime"):
		return &FieldInfo{Type: DateTimeField, Label: getLabel(line), Note: getNote(line), Pattern: getPattern(line), Required: getRequired(line), Readonly: getReadonly(line), Action: getAction(line)}
	case strings.Contains(line, "@email"):
		return &FieldInfo{Type: EmailField, Label: getLabel(line), Note: getNote(line), Pattern: getPattern(line), Required: getRequired(line), Readonly: getReadonly(line), Action: getAction(line)}
	case strings.Contains(line, "@file"):
		return &FieldInfo{Type: FileField, Label: getLabel(line), Note: getNote(line), Required: getRequired(line), Readonly: getReadonly(line), Action: getAction(line)}
	case strings.Contains(line, "@image"):
		return &FieldInfo{Type: ImageField, Label: getLabel(line), Note: getNote(line), Action: getAction(line)}
	case strings.Contains(line, "@month"):
		return &FieldInfo{Type: MonthField, Label: getLabel(line), Note: getNote(line), Action: getAction(line)}
	case strings.Contains(line, "@number"):
		return &FieldInfo{Type: NumberField, Label: getLabel(line), Note: getNote(line), Required: getRequired(line), Readonly: getReadonly(line), Action: getAction(line)}
	case strings.Contains(line, "@password"):
		return &FieldInfo{Type: PasswordField, Label: getLabel(line), Note: getNote(line), Required: getRequired(line), Readonly: getReadonly(line), Action: getAction(line)}
	case strings.Contains(line, "@radio"):
		options := []string{}
		if m := regexp.MustCompile(`options=([^\s]*)`).FindStringSubmatch(line); len(m) > 1 {
			options = strings.Split(m[1], ",")
		}
		return &FieldInfo{Type: RadioField, Label: getLabel(line), Options: options, Note: getNote(line), Action: getAction(line)}
	case strings.Contains(line, "@range"):
		return &FieldInfo{Type: RangeField, Label: getLabel(line), Note: getNote(line), Default: getDefault(line), Action: getAction(line)}
	case strings.Contains(line, "@reset"):
		return &FieldInfo{Type: ResetField, Label: getLabel(line), Note: getNote(line), Default: getDefault(line), Action: getAction(line)}
	case strings.Contains(line, "@tel"):
		return &FieldInfo{Type: TelField, Label: getLabel(line), Note: getNote(line), Pattern: getPattern(line), Required: getRequired(line), Readonly: getReadonly(line), Action: getAction(line)}
	case strings.Contains(line, "@text"):
		return &FieldInfo{Type: TextField, Label: getLabel(line), Note: getNote(line), Pattern: getPattern(line), Required: getRequired(line), Readonly: getReadonly(line), Action: getAction(line)}
	case strings.Contains(line, "@time"):
		return &FieldInfo{Type: TimeField, Label: getLabel(line), Note: getNote(line), Action: getAction(line)}
	case strings.Contains(line, "@url"):
		return &FieldInfo{Type: URLField, Label: getLabel(line), Note: getNote(line), Required: getRequired(line), Readonly: getReadonly(line), Action: getAction(line)}
	case strings.Contains(line, "@week"):
		return &FieldInfo{Type: WeekField, Label: getLabel(line), Note: getNote(line), Action: getAction(line)}
	case strings.Contains(line, "@select"):
		label := getLabel(line)
		options := []string{}
		if m := regexp.MustCompile(`options=([^\s]*)`).FindStringSubmatch(line); len(m) > 1 {
			options = strings.Split(m[1], ",")
		}
		return &FieldInfo{Type: SelectField, Label: label, Options: options, Note: getNote(line), Required: getRequired(line), Readonly: getReadonly(line), Action: getAction(line)}
	case strings.Contains(line, "@boolean"):
		return &FieldInfo{Type: BooleanField, Label: getLabel(line), Note: getNote(line), Required: getRequired(line), Readonly: getReadonly(line), Action: getAction(line)}
	case strings.Contains(line, "@list"):
		return &FieldInfo{Type: ListField, Label: getLabel(line), Note: getNote(line), Required: getRequired(line), Readonly: getReadonly(line), Action: getAction(line)}
	default:
		return nil
	}
}

func getExistingValue(key, defaultValue, outputFile string) string {
	file, err := os.Open(outputFile)
	if err != nil {
		return defaultValue
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	keyRegex := regexp.MustCompile(fmt.Sprintf(`^%s=(.*)$`, regexp.QuoteMeta(key)))

	for scanner.Scan() {
		line := scanner.Text()
		if m := keyRegex.FindStringSubmatch(line); len(m) > 1 {
			return m[1]
		}
	}

	return defaultValue
}

func promptForValue(f *FieldInfo, label, defaultValue string) (string, error) {
	var value string
	var err error
	for {
		switch f.Type {
		case TextField:
			value, err = promptText(label, defaultValue, f.Pattern)
		case NumberField:
			value, err = promptNumber(label, defaultValue)
		case FloatField:
			value, err = promptFloat(label, defaultValue)
		case DateField:
			value, err = promptDate(label, defaultValue, f.Pattern)
		case DateTimeField:
			value, err = promptText(label+" (YYYY-MM-DD HH:MM)", defaultValue, `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$`)
		case FileField:
			value, err = promptFile(label, defaultValue)
		case URLField:
			value, err = promptURL(label, defaultValue)
		case PasswordField:
			value, err = promptPassword(label, defaultValue)
		case SelectField:
			value, err = promptSelect(label, f.Options, defaultValue)
		case BooleanField:
			value, err = promptBoolean(label, defaultValue)
		case ListField:
			value, err = promptList(label, defaultValue)
		case ButtonField:
			value, err = promptButton(label, defaultValue)
		case CheckboxField:
			value, err = promptBoolean(label, defaultValue)
		case ColorField:
			value, err = promptColor(label, defaultValue)
		case EmailField:
			value, err = promptEmail(label, defaultValue)
		case ImageField:
			value, err = promptImage(label, defaultValue)
		case MonthField:
			value, err = promptMonth(label, defaultValue)
		case RadioField:
			value, err = promptRadio(label, f.Options, defaultValue)
		case RangeField:
			value, err = promptRange(label, defaultValue)
		case ResetField:
			value, err = promptReset(label, f.Default)
		case TelField:
			value, err = promptTel(label, defaultValue)
		case TimeField:
			value, err = promptTime(label, defaultValue)
		case WeekField:
			value, err = promptWeek(label, defaultValue)
		default:
			value = defaultValue
			err = nil
		}
		if err != nil {
			return "", err
		}
		if !f.Required || strings.TrimSpace(value) != "" {
			return value, nil
		}
		fmt.Println("This field is required. Please enter a value.")
		// Next loop will prompt again
	}
}

func promptText(label, defaultValue, pattern string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		prompt := label
		if defaultValue != "" {
			prompt += fmt.Sprintf(" [%s]", defaultValue)
		}
		prompt += ": "
		fmt.Print(prompt)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			return defaultValue, nil
		}

		// Validate with custom pattern if provided
		if pattern != "" {
			matched, err := regexp.MatchString(pattern, input)
			if err != nil {
				fmt.Printf("Invalid pattern: %v\n", err)
				continue
			}
			if !matched {
				fmt.Printf("Invalid input. Must match pattern: %s\n", pattern)
				continue
			}
		}

		return input, nil
	}
}

func promptNumber(label, defaultValue string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		prompt := label
		if defaultValue != "" {
			prompt += fmt.Sprintf(" [%s]", defaultValue)
		}
		prompt += ": "
		fmt.Print(prompt)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			return defaultValue, nil
		}

		// Validate it's a valid integer
		matched, _ := regexp.MatchString(`^-?\d+$`, input)
		if !matched {
			fmt.Println("Invalid input. Must be a valid integer number.")
			continue
		}

		return input, nil
	}
}

func promptFloat(label, defaultValue string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		prompt := label
		if defaultValue != "" {
			prompt += fmt.Sprintf(" [%s]", defaultValue)
		}
		prompt += ": "
		fmt.Print(prompt)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			return defaultValue, nil
		}

		// Validate it's a valid float
		matched, _ := regexp.MatchString(`^-?\d+\.?\d*$`, input)
		if !matched {
			fmt.Println("Invalid input. Must be a valid floating point number.")
			continue
		}

		return input, nil
	}
}

func promptDate(label, defaultValue, pattern string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	// Default date pattern if not provided
	if pattern == "" {
		pattern = `^\d{4}-\d{2}-\d{2}$` // YYYY-MM-DD
	}

	for {
		prompt := label
		if defaultValue != "" {
			prompt += fmt.Sprintf(" [%s]", defaultValue)
		}
		prompt += ": "
		fmt.Print(prompt)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			return defaultValue, nil
		}

		// Validate with date pattern
		matched, err := regexp.MatchString(pattern, input)
		if err != nil {
			fmt.Printf("Invalid pattern: %v\n", err)
			continue
		}
		if !matched {
			fmt.Printf("Invalid date format. Expected pattern: %s (e.g., YYYY-MM-DD)\n", pattern)
			continue
		}

		return input, nil
	}
}

func promptFile(label, defaultValue string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		prompt := label
		if defaultValue != "" {
			prompt += fmt.Sprintf(" [%s]", defaultValue)
		}
		prompt += ": "
		fmt.Print(prompt)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			return defaultValue, nil
		}

		// Validate it looks like a file path
		matched, _ := regexp.MatchString(`^[\/~\.]?[\w\-\.\/]+$`, input)
		if !matched {
			fmt.Println("Invalid input. Must be a valid file path.")
			continue
		}

		return input, nil
	}
}

func promptURL(label, defaultValue string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		prompt := label
		if defaultValue != "" {
			prompt += fmt.Sprintf(" [%s]", defaultValue)
		}
		prompt += ": "
		fmt.Print(prompt)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			return defaultValue, nil
		}

		// Validate it looks like a URL
		matched, _ := regexp.MatchString(`^https?:\/\/[^\s]+$`, input)
		if !matched {
			fmt.Println("Invalid input. Must be a valid URL (http:// or https://).")
			continue
		}

		return input, nil
	}
}

func promptInput(label, defaultValue string) (string, error) {
	prompt := label
	if defaultValue != "" {
		prompt += fmt.Sprintf(" [%s]", defaultValue)
	}
	prompt += ": "
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue, nil
	}
	return input, nil
}

func promptPassword(label, defaultValue string) (string, error) {
	prompt := label
	if defaultValue != "" {
		prompt += " [press enter to keep current]: "
	} else {
		prompt += ": "
	}
	fmt.Print(prompt)
	passBytes, _ := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	pass := strings.TrimSpace(string(passBytes))
	if pass == "" {
		return defaultValue, nil
	}
	return pass, nil
}

func promptSelect(label string, options []string, defaultValue string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(label + ":")
		for i, opt := range options {
			if opt == defaultValue {
				fmt.Printf("  %d) %s (current)\n", i+1, opt)
			} else {
				fmt.Printf("  %d) %s\n", i+1, opt)
			}
		}
		fmt.Printf("Select [1-%d]: ", len(options))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			return defaultValue, nil
		}
		var sel int
		if _, err := fmt.Sscanf(input, "%d", &sel); err == nil && sel >= 1 && sel <= len(options) {
			return options[sel-1], nil
		}
		fmt.Println("Invalid selection, try again.")
	}
}

func promptBoolean(label, defaultValue string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		defaultBool := strings.EqualFold(defaultValue, "true") || strings.EqualFold(defaultValue, "yes") || defaultValue == "1"
		prompt := label
		if defaultBool {
			prompt += " [Y/n]: "
		} else {
			prompt += " [y/N]: "
		}
		fmt.Print(prompt)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "y", "yes", "true":
			return "true", nil
		case "n", "no", "false":
			return "false", nil
		case "":
			if defaultBool {
				return "true", nil
			}
			return "false", nil
		default:
			fmt.Println("Invalid input. Must be true/false, yes/no, or y/n.")
		}
	}
}

func promptList(label, defaultValue string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		prompt := label + " (comma-separated)"
		if defaultValue != "" {
			prompt += fmt.Sprintf(" [%s]", defaultValue)
		}
		prompt += ": "
		fmt.Print(prompt)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			return defaultValue, nil
		}

		// Validate it's a comma-separated list
		matched, _ := regexp.MatchString(`^[^,]+(,[^,]+)*$`, input)
		if !matched {
			fmt.Println("Invalid input. Must be comma-separated values (e.g., value1,value2,value3).")
			continue
		}

		return input, nil
	}
}

// For new types, you can implement simple prompt stubs like below.
// For brevity, only the main prompt types are implemented. You can expand as needed.

func promptButton(label, defaultValue string) (string, error) {
	fmt.Printf("%s [button, no input]: %s\n", label, defaultValue)
	return defaultValue, nil
}

func promptColor(label, defaultValue string) (string, error) {
	return promptText(label+" (color hex or name)", defaultValue, "")
}

func promptEmail(label, defaultValue string) (string, error) {
	return promptText(label+" (email)", defaultValue, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
}

func promptImage(label, defaultValue string) (string, error) {
	return promptFile(label+" (image path)", defaultValue)
}

func promptMonth(label, defaultValue string) (string, error) {
	return promptText(label+" (YYYY-MM)", defaultValue, `^\d{4}-\d{2}$`)
}

func promptRadio(label string, options []string, defaultValue string) (string, error) {
	return promptSelect(label+" (radio)", options, defaultValue)
}

func promptRange(label, defaultValue string) (string, error) {
	return promptNumber(label+" (range)", defaultValue)
}

func promptReset(label, defaultValue string) (string, error) {
	fmt.Printf("%s [reset, default=%s]\n", label, defaultValue)
	return defaultValue, nil
}

func promptTel(label, defaultValue string) (string, error) {
	return promptText(label+" (telephone)", defaultValue, `^\+?[0-9\- ]+$`)
}

func promptTime(label, defaultValue string) (string, error) {
	return promptText(label+" (HH:MM)", defaultValue, `^\d{2}:\d{2}$`)
}

func promptWeek(label, defaultValue string) (string, error) {
	return promptText(label+" (YYYY-Www)", defaultValue, `^\d{4}-W\d{2}$`)
}
