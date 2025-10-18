package field

import (
	"regexp"
	"strings"
)

func ParseDirective(line string) *Info {
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

	getScript := func(line string) string {
		re := regexp.MustCompile(`script=(?:"([^"]*)"|` + "`" + `([^` + "`" + `]*)` + "`" + `)`)
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
		return &Info{Type: HiddenField, Note: getNote(line), Script: getScript(line)}
	case strings.Contains(line, "@skip"):
		return &Info{Type: SkipField, Note: getNote(line), Script: getScript(line)}
	case strings.Contains(line, "@readonly"):
		return &Info{Type: ReadonlyField, Label: getLabel(line), Note: getNote(line), Readonly: true, Script: getScript(line)}
	case strings.Contains(line, "@button"):
		return &Info{Type: ButtonField, Label: getLabel(line), Note: getNote(line), Script: getScript(line)}
	case strings.Contains(line, "@checkbox"):
		return &Info{Type: CheckboxField, Label: getLabel(line), Note: getNote(line), Script: getScript(line)}
	case strings.Contains(line, "@color"):
		return &Info{Type: ColorField, Label: getLabel(line), Note: getNote(line), Script: getScript(line)}
	case strings.Contains(line, "@date"):
		return &Info{Type: DateField, Label: getLabel(line), Note: getNote(line), Pattern: getPattern(line), Required: getRequired(line), Readonly: getReadonly(line), Script: getScript(line)}
	case strings.Contains(line, "@datetime"):
		return &Info{Type: DateTimeField, Label: getLabel(line), Note: getNote(line), Pattern: getPattern(line), Required: getRequired(line), Readonly: getReadonly(line), Script: getScript(line)}
	case strings.Contains(line, "@email"):
		return &Info{Type: EmailField, Label: getLabel(line), Note: getNote(line), Pattern: getPattern(line), Required: getRequired(line), Readonly: getReadonly(line), Script: getScript(line)}
	case strings.Contains(line, "@file"):
		return &Info{Type: FileField, Label: getLabel(line), Note: getNote(line), Required: getRequired(line), Readonly: getReadonly(line), Script: getScript(line)}
	case strings.Contains(line, "@image"):
		return &Info{Type: ImageField, Label: getLabel(line), Note: getNote(line), Script: getScript(line)}
	case strings.Contains(line, "@month"):
		return &Info{Type: MonthField, Label: getLabel(line), Note: getNote(line), Script: getScript(line)}
	case strings.Contains(line, "@number"):
		return &Info{Type: NumberField, Label: getLabel(line), Note: getNote(line), Required: getRequired(line), Readonly: getReadonly(line), Script: getScript(line)}
	case strings.Contains(line, "@password"):
		return &Info{Type: PasswordField, Label: getLabel(line), Note: getNote(line), Required: getRequired(line), Readonly: getReadonly(line), Script: getScript(line)}
	case strings.Contains(line, "@radio"):
		options := []string{}
		if m := regexp.MustCompile(`options=([^\s]*)`).FindStringSubmatch(line); len(m) > 1 {
			options = strings.Split(m[1], ",")
		}
		return &Info{Type: RadioField, Label: getLabel(line), Options: options, Note: getNote(line), Script: getScript(line)}
	case strings.Contains(line, "@range"):
		return &Info{Type: RangeField, Label: getLabel(line), Note: getNote(line), Default: getDefault(line), Script: getScript(line)}
	case strings.Contains(line, "@reset"):
		return &Info{Type: ResetField, Label: getLabel(line), Note: getNote(line), Default: getDefault(line), Script: getScript(line)}
	case strings.Contains(line, "@tel"):
		return &Info{Type: TelField, Label: getLabel(line), Note: getNote(line), Pattern: getPattern(line), Required: getRequired(line), Readonly: getReadonly(line), Script: getScript(line)}
	case strings.Contains(line, "@text"):
		return &Info{Type: TextField, Label: getLabel(line), Note: getNote(line), Pattern: getPattern(line), Required: getRequired(line), Readonly: getReadonly(line), Script: getScript(line)}
	case strings.Contains(line, "@time"):
		return &Info{Type: TimeField, Label: getLabel(line), Note: getNote(line), Script: getScript(line)}
	case strings.Contains(line, "@url"):
		return &Info{Type: URLField, Label: getLabel(line), Note: getNote(line), Required: getRequired(line), Readonly: getReadonly(line), Script: getScript(line)}
	case strings.Contains(line, "@week"):
		return &Info{Type: WeekField, Label: getLabel(line), Note: getNote(line), Script: getScript(line)}
	case strings.Contains(line, "@select"):
		label := getLabel(line)
		options := []string{}
		if m := regexp.MustCompile(`options=([^\s]*)`).FindStringSubmatch(line); len(m) > 1 {
			options = strings.Split(m[1], ",")
		}
		return &Info{Type: SelectField, Label: label, Options: options, Note: getNote(line), Required: getRequired(line), Readonly: getReadonly(line), Script: getScript(line)}
	case strings.Contains(line, "@boolean"):
		return &Info{Type: BooleanField, Label: getLabel(line), Note: getNote(line), Required: getRequired(line), Readonly: getReadonly(line), Script: getScript(line)}
	case strings.Contains(line, "@list"):
		return &Info{Type: ListField, Label: getLabel(line), Note: getNote(line), Required: getRequired(line), Readonly: getReadonly(line), Script: getScript(line)}
	default:
		return nil
	}
}
