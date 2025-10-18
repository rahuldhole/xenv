package processor

import (
	"os"
	"testing"
)

func TestDetermineOutputFile(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"config.xenv", ".config"},
		{"env.template", ".env"},
		{".env.example", ".env"},
		{"myapp.xenv", ".myapp"},
		{"/path/to/config.xenv", "/path/to/.config"},
	}

	for _, tt := range tests {
		result := DetermineOutputFile(tt.input)
		if result != tt.expected {
			t.Errorf("DetermineOutputFile(%s) = %s; want %s", tt.input, result, tt.expected)
		}
	}
}

func TestCheckForDSL(t *testing.T) {
	// Create a temporary test file
	tmpFile := "test_check_dsl.tmp"
	content := `# @text label="Test"
VAR=value
`
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile)

	hasDSL := CheckForDSL(tmpFile)
	if !hasDSL {
		t.Error("Expected CheckForDSL to return true for file with @text directive")
	}
}
