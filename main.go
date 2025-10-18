package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/term"
)

type FieldType int

const (
	NoField FieldType = iota
	InputField
	PasswordField
	SelectField
	CheckboxField
)

type FieldInfo struct {
	Type    FieldType
	Label   string
	Options []string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <form-file>\n", os.Args[0])
		os.Exit(1)
	}

	formFile := os.Args[1]

	var outputFile string
	if strings.HasSuffix(formFile, ".xenv") {
		outputFile = filepath.Join(filepath.Dir(formFile), "."+strings.TrimSuffix(filepath.Base(formFile), ".xenv"))
	} else {
		outputFile = formFile
	}

	if _, err := os.Stat(formFile); os.IsNotExist(err) {
		fmt.Printf("Error: Form file '%s' not found.\n", formFile)
		os.Exit(1)
	}

	fmt.Printf("Interactive configuration for %s\n", filepath.Base(formFile))
	fmt.Println(strings.Repeat("-", 50))

	outputLines, err := processFormFile(formFile, outputFile)
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

func processFormFile(formFile, outputFile string) ([]string, error) {
	file, err := os.Open(formFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	keyValueRegex := regexp.MustCompile(`^([^=]+)=(.*)$`)
	var outputLines []string
	var currentFieldInfo *FieldInfo

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			fieldInfo := parseFieldDirective(line)
			if fieldInfo != nil {
				currentFieldInfo = fieldInfo
			}
			outputLines = append(outputLines, line)
			continue
		}

		if matches := keyValueRegex.FindStringSubmatch(line); len(matches) > 0 {
			key := matches[1]
			defaultValue := matches[2]
			existingValue := getExistingValue(key, defaultValue, outputFile)

			var newValue string
			if currentFieldInfo != nil && currentFieldInfo.Type != NoField {
				label := currentFieldInfo.Label
				if label == "" {
					label = key
				}
				newValue, err = promptForValue(currentFieldInfo, label, existingValue)
				if err != nil {
					return nil, err
				}
				currentFieldInfo = nil
			} else {
				newValue = existingValue
			}

			outputLines = append(outputLines, key+"="+newValue)
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

	switch {
	case strings.Contains(line, "@input"):
		return &FieldInfo{Type: InputField, Label: getLabel(line)}
	case strings.Contains(line, "@password"):
		return &FieldInfo{Type: PasswordField, Label: getLabel(line)}
	case strings.Contains(line, "@select"):
		label := getLabel(line)
		options := []string{}
		if m := regexp.MustCompile(`options=([^\s]*)`).FindStringSubmatch(line); len(m) > 1 {
			options = strings.Split(m[1], ",")
		}
		return &FieldInfo{Type: SelectField, Label: label, Options: options}
	case strings.Contains(line, "@checkbox"):
		return &FieldInfo{Type: CheckboxField, Label: getLabel(line)}
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
	switch f.Type {
	case InputField:
		return promptInput(label, defaultValue)
	case PasswordField:
		return promptPassword(label, defaultValue)
	case SelectField:
		return promptSelect(label, f.Options, defaultValue)
	case CheckboxField:
		return promptCheckbox(label, defaultValue)
	default:
		return defaultValue, nil
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

func promptCheckbox(label, defaultValue string) (string, error) {
	defaultBool := strings.EqualFold(defaultValue, "true") || strings.EqualFold(defaultValue, "yes") || defaultValue == "1"
	prompt := label
	if defaultBool {
		prompt += " [Y/n]: "
	} else {
		prompt += " [y/N]: "
	}
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
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
		return "false", nil
	}
}
