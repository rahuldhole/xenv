package prompt

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/term"

	"github.com/rahuldhole/xenv/app/field"
)

func ForValue(f *field.Info, label, defaultValue string) (string, error) {
	var value string
	var err error
	for {
		switch f.Type {
		case field.TextField:
			value, err = Text(label, defaultValue, f.Pattern)
		case field.NumberField:
			value, err = Number(label, defaultValue)
		case field.FloatField:
			value, err = Float(label, defaultValue)
		case field.DateField:
			value, err = Date(label, defaultValue, f.Pattern)
		case field.DateTimeField:
			value, err = Text(label+" (YYYY-MM-DD HH:MM)", defaultValue, `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$`)
		case field.FileField:
			value, err = File(label, defaultValue)
		case field.URLField:
			value, err = URL(label, defaultValue)
		case field.PasswordField:
			value, err = Password(label, defaultValue)
		case field.SelectField:
			value, err = Select(label, f.Options, defaultValue)
		case field.BooleanField:
			value, err = Boolean(label, defaultValue)
		case field.ListField:
			value, err = List(label, defaultValue)
		case field.ButtonField:
			value, err = Button(label, defaultValue)
		case field.CheckboxField:
			value, err = Boolean(label, defaultValue)
		case field.ColorField:
			value, err = Color(label, defaultValue)
		case field.EmailField:
			value, err = Email(label, defaultValue)
		case field.ImageField:
			value, err = Image(label, defaultValue)
		case field.MonthField:
			value, err = Month(label, defaultValue)
		case field.RadioField:
			value, err = Radio(label, f.Options, defaultValue)
		case field.RangeField:
			value, err = Range(label, defaultValue)
		case field.ResetField:
			value, err = Reset(label, f.Default)
		case field.TelField:
			value, err = Tel(label, defaultValue)
		case field.TimeField:
			value, err = Time(label, defaultValue)
		case field.WeekField:
			value, err = Week(label, defaultValue)
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
	}
}

func Text(label, defaultValue, pattern string) (string, error) {
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

func Number(label, defaultValue string) (string, error) {
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

		matched, _ := regexp.MatchString(`^-?\d+$`, input)
		if !matched {
			fmt.Println("Invalid input. Must be a valid integer number.")
			continue
		}

		return input, nil
	}
}

func Float(label, defaultValue string) (string, error) {
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

		matched, _ := regexp.MatchString(`^-?\d+\.?\d*$`, input)
		if !matched {
			fmt.Println("Invalid input. Must be a valid floating point number.")
			continue
		}

		return input, nil
	}
}

func Date(label, defaultValue, pattern string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	if pattern == "" {
		pattern = `^\d{4}-\d{2}-\d{2}$`
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

func File(label, defaultValue string) (string, error) {
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

		matched, _ := regexp.MatchString(`^[\/~\.]?[\w\-\.\/]+$`, input)
		if !matched {
			fmt.Println("Invalid input. Must be a valid file path.")
			continue
		}

		return input, nil
	}
}

func URL(label, defaultValue string) (string, error) {
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

		matched, _ := regexp.MatchString(`^https?:\/\/[^\s]+$`, input)
		if !matched {
			fmt.Println("Invalid input. Must be a valid URL (http:// or https://).")
			continue
		}

		return input, nil
	}
}

func Input(label, defaultValue string) (string, error) {
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

func Password(label, defaultValue string) (string, error) {
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

func Select(label string, options []string, defaultValue string) (string, error) {
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

func Boolean(label, defaultValue string) (string, error) {
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

func List(label, defaultValue string) (string, error) {
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

		matched, _ := regexp.MatchString(`^[^,]+(,[^,]+)*$`, input)
		if !matched {
			fmt.Println("Invalid input. Must be comma-separated values (e.g., value1,value2,value3).")
			continue
		}

		return input, nil
	}
}

func Button(label, defaultValue string) (string, error) {
	fmt.Printf("%s [button, no input]: %s\n", label, defaultValue)
	return defaultValue, nil
}

func Color(label, defaultValue string) (string, error) {
	return Text(label+" (color hex or name)", defaultValue, "")
}

func Email(label, defaultValue string) (string, error) {
	return Text(label+" (email)", defaultValue, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
}

func Image(label, defaultValue string) (string, error) {
	return File(label+" (image path)", defaultValue)
}

func Month(label, defaultValue string) (string, error) {
	return Text(label+" (YYYY-MM)", defaultValue, `^\d{4}-\d{2}$`)
}

func Radio(label string, options []string, defaultValue string) (string, error) {
	return Select(label+" (radio)", options, defaultValue)
}

func Range(label, defaultValue string) (string, error) {
	return Number(label+" (range)", defaultValue)
}

func Reset(label, defaultValue string) (string, error) {
	fmt.Printf("%s [reset, default=%s]\n", label, defaultValue)
	return defaultValue, nil
}

func Tel(label, defaultValue string) (string, error) {
	return Text(label+" (telephone)", defaultValue, `^\+?[0-9\- ]+$`)
}

func Time(label, defaultValue string) (string, error) {
	return Text(label+" (HH:MM)", defaultValue, `^\d{2}:\d{2}$`)
}

func Week(label, defaultValue string) (string, error) {
	return Text(label+" (YYYY-Www)", defaultValue, `^\d{4}-W\d{2}$`)
}
