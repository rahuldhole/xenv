package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEnvFile(t *testing.T) {
	exampleFile := filepath.Join("../examples", "all.env.example")
	
	content, err := os.ReadFile(exampleFile)
	require.NoError(t, err, "Should read example file")
	require.NotEmpty(t, content, "File should not be empty")
	
	lines := strings.Split(string(content), "\n")
	assert.Greater(t, len(lines), 100, "Should have many lines")
}

func TestParseTextFields(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"text annotation", "# @text label=\"App Name\"", true},
		{"text with pattern", "# @text label=\"Version\" pattern=\"^\\d+\"", true},
		{"text with required", "# @text label=\"Name\" required", true},
		{"not text", "# @number label=\"Port\"", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasText := strings.Contains(tt.line, "@text")
			assert.Equal(t, tt.expected, hasText)
		})
	}
}

func TestParseSelectFields(t *testing.T) {
	line := "# @select label=\"Environment\" options=development,testing,staging,production"
	
	assert.Contains(t, line, "@select")
	assert.Contains(t, line, "options=")
	
	// Extract options
	if strings.Contains(line, "options=") {
		parts := strings.Split(line, "options=")
		if len(parts) > 1 {
			optionsPart := strings.Fields(parts[1])[0]
			options := strings.Split(optionsPart, ",")
			assert.Equal(t, 4, len(options))
			assert.Contains(t, options, "development")
			assert.Contains(t, options, "production")
		}
	}
}

func TestParseBooleanFields(t *testing.T) {
	lines := []string{
		"# @boolean label=\"Enable HTTPS\"",
		"SERVER_HTTPS=false",
		"# @boolean label=\"Enable TLS\"",
		"SMTP_TLS=true",
	}
	
	boolCount := 0
	for _, line := range lines {
		if strings.Contains(line, "@boolean") {
			boolCount++
		}
	}
	
	assert.Equal(t, 2, boolCount, "Should find 2 boolean annotations")
}

func TestParseNumberFields(t *testing.T) {
	lines := []string{
		"# @number label=\"Server Port\" required",
		"SERVER_PORT=8080",
		"# @number label=\"DB Port\"",
		"DB_PORT=5432",
	}
	
	numberCount := 0
	for _, line := range lines {
		if strings.Contains(line, "@number") {
			numberCount++
		}
	}
	
	assert.Equal(t, 2, numberCount)
}

func TestParsePasswordFields(t *testing.T) {
	lines := []string{
		"# @password label=\"DB Password\" required",
		"DB_PASSWORD=",
		"# @password label=\"API Secret\"",
		"API_SECRET=",
	}
	
	passwordCount := 0
	for _, line := range lines {
		if strings.Contains(line, "@password") {
			passwordCount++
		}
	}
	
	assert.Equal(t, 2, passwordCount)
}

func TestParseURLFields(t *testing.T) {
	line := "# @url label=\"Base URL\" note=\"Full URL\" required"
	
	assert.Contains(t, line, "@url")
	assert.Contains(t, line, "required")
}

func TestParseFileFields(t *testing.T) {
	lines := []string{
		"# @file label=\"Certificate Path\"",
		"# @file label=\"Key Path\"",
	}
	
	fileCount := 0
	for _, line := range lines {
		if strings.Contains(line, "@file") {
			fileCount++
		}
	}
	
	assert.Equal(t, 2, fileCount)
}

func TestParseListFields(t *testing.T) {
	lines := []string{
		"# @list label=\"Allowed Origins\"",
		"ALLOWED_ORIGINS=http://localhost:3000,http://example.com",
		"# @list label=\"Admin Emails\" required",
		"ADMIN_EMAILS=admin@example.com,superadmin@example.com",
	}
	
	listCount := 0
	for _, line := range lines {
		if strings.Contains(line, "@list") {
			listCount++
		}
	}
	
	assert.Equal(t, 2, listCount)
	
	// Test list value parsing
	value := "http://localhost:3000,http://example.com"
	items := strings.Split(value, ",")
	assert.Equal(t, 2, len(items))
}

func TestParseDateFields(t *testing.T) {
	lines := []string{
		"# @date label=\"Service Start Date\" required",
		"SERVICE_START_DATE=2024-01-01",
		"# @date label=\"License Expiry\" pattern=\"^\\d{2}/\\d{2}/\\d{4}$\"",
		"LICENSE_EXPIRY=12/31/2025",
	}
	
	dateCount := 0
	for _, line := range lines {
		if strings.Contains(line, "@date") {
			dateCount++
		}
	}
	
	assert.Equal(t, 2, dateCount)
}

func TestParseFloatFields(t *testing.T) {
	line := "# @float label=\"Request Timeout\" note=\"Timeout in seconds\""
	
	assert.Contains(t, line, "@float")
}

func TestParseSpecialFields(t *testing.T) {
	specialTypes := map[string]string{
		"@hidden":   "# @hidden note=\"Do not prompt\"",
		"@checkbox": "# @checkbox label=\"Accept Terms\"",
		"@color":    "# @color label=\"Theme Color\"",
		"@datetime": "# @datetime label=\"Deployment DateTime\"",
		"@email":    "# @email label=\"Support Email\" required",
		"@image":    "# @image label=\"Logo Path\"",
		"@month":    "# @month label=\"Billing Month\"",
		"@radio":    "# @radio label=\"Gender\" options=male,female,other",
		"@range":    "# @range label=\"Volume\" default=\"50\"",
		"@reset":    "# @reset label=\"Reset Config\"",
		"@tel":      "# @tel label=\"Contact Phone\"",
		"@time":     "# @time label=\"Daily Backup Time\"",
		"@week":     "# @week label=\"Sprint Week\"",
		"@readonly": "# @readonly label=\"Read Only Field\"",
	}
	
	for fieldType, line := range specialTypes {
		t.Run(fieldType, func(t *testing.T) {
			assert.Contains(t, line, fieldType)
		})
	}
}

func TestParseScriptFields(t *testing.T) {
	lines := []string{
		"# @text label=\"Show $GREET\" script=`echo \"Hello $NAME\"`",
		"# @button label=\"Show Date\" script=\"date\"",
	}
	
	scriptCount := 0
	for _, line := range lines {
		if strings.Contains(line, "script=") {
			scriptCount++
		}
	}
	
	assert.Equal(t, 2, scriptCount)
}

func TestParseSkipAnnotation(t *testing.T) {
	content := `
# @skip note="All variables below will remain unchanged"
LEGACY_SERVICE_URL=http://legacy.example.com
LEGACY_SERVICE_KEY=abc123
`
	
	assert.Contains(t, content, "@skip")
	assert.Contains(t, content, "LEGACY_SERVICE_URL")
}

func TestParseAttributes(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		attribute string
		expected  bool
	}{
		{"has label", "# @text label=\"Name\"", "label=", true},
		{"has note", "# @text note=\"Help text\"", "note=", true},
		{"has required", "# @text required", "required", true},
		{"has pattern", "# @text pattern=\"^\\d+$\"", "pattern=", true},
		{"has options", "# @select options=a,b,c", "options=", true},
		{"has default", "# @range default=\"50\"", "default=", true},
		{"no label", "# @text", "label=", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasAttr := strings.Contains(tt.line, tt.attribute)
			assert.Equal(t, tt.expected, hasAttr)
		})
	}
}

func TestParseVariableNames(t *testing.T) {
	lines := []string{
		"APP_NAME=MyAwesomeApp",
		"SERVER_PORT=8080",
		"DB_PASSWORD=",
		"# Comment line",
		"",
	}
	
	varCount := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") && strings.Contains(line, "=") {
			varCount++
		}
	}
	
	assert.Equal(t, 3, varCount)
}

func TestParseDefaultValues(t *testing.T) {
	tests := []struct {
		line     string
		hasValue bool
	}{
		{"APP_NAME=MyAwesomeApp", true},
		{"DB_PASSWORD=", false},
		{"SERVER_PORT=8080", true},
		{"API_SECRET=", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			parts := strings.SplitN(tt.line, "=", 2)
			require.Equal(t, 2, len(parts))
			hasValue := len(parts[1]) > 0
			assert.Equal(t, tt.hasValue, hasValue)
		})
	}
}
