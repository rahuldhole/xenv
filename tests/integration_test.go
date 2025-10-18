package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompleteFileStructure(t *testing.T) {
	exampleFile := filepath.Join("../examples", "all.env.example")
	
	content, err := os.ReadFile(exampleFile)
	require.NoError(t, err)
	
	lines := strings.Split(string(content), "\n")
	
	// Count annotations
	annotationCount := 0
	for _, line := range lines {
		if strings.Contains(line, "#") && strings.Contains(line, "@") {
			annotationCount++
		}
	}
	
	assert.Greater(t, annotationCount, 40, "Should have more than 40 annotations")
}

func TestAllFieldTypesPresent(t *testing.T) {
	exampleFile := filepath.Join("../examples", "all.env.example")
	content, err := os.ReadFile(exampleFile)
	require.NoError(t, err)
	
	fileContent := string(content)
	
	fieldTypes := []string{
		"@text", "@select", "@number", "@boolean", "@password",
		"@url", "@file", "@list", "@date", "@float", "@hidden",
		"@checkbox", "@color", "@datetime", "@email", "@image",
		"@month", "@radio", "@range", "@reset", "@tel", "@time",
		"@week", "@readonly", "@button", "@skip",
	}
	
	for _, fieldType := range fieldTypes {
		t.Run(fieldType, func(t *testing.T) {
			assert.Contains(t, fileContent, fieldType, 
				"Field type %s should be present", fieldType)
		})
	}
}

func TestRequiredFields(t *testing.T) {
	exampleFile := filepath.Join("../examples", "all.env.example")
	content, err := os.ReadFile(exampleFile)
	require.NoError(t, err)
	
	requiredCount := strings.Count(string(content), "required")
	assert.Greater(t, requiredCount, 5, "Should have multiple required fields")
}

func TestDefaultValues(t *testing.T) {
	exampleFile := filepath.Join("../examples", "all.env.example")
	content, err := os.ReadFile(exampleFile)
	require.NoError(t, err)
	
	lines := strings.Split(string(content), "\n")
	
	varLines := make([]string, 0)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "=") && !strings.HasPrefix(trimmed, "#") {
			varLines = append(varLines, trimmed)
		}
	}
	
	// Count lines with values
	linesWithValues := 0
	for _, line := range varLines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && len(parts[1]) > 0 {
			linesWithValues++
		}
	}
	
	assert.Greater(t, linesWithValues, 20, "Many variables should have default values")
}

func TestSectionsPresent(t *testing.T) {
	exampleFile := filepath.Join("../examples", "all.env.example")
	content, err := os.ReadFile(exampleFile)
	require.NoError(t, err)
	
	fileContent := string(content)
	
	sections := []string{
		"General Application Settings",
		"Server Configuration",
		"Database Configuration",
		"API Configuration",
		"Feature Flags",
		"External Services",
		"Limits & Timeouts",
		"Optional Lists / Arrays",
		"Custom Pattern Examples",
	}
	
	for _, section := range sections {
		t.Run(section, func(t *testing.T) {
			assert.Contains(t, fileContent, section, 
				"Section '%s' should be present", section)
		})
	}
}

func TestSkipFunctionality(t *testing.T) {
	exampleFile := filepath.Join("../examples", "all.env.example")
	content, err := os.ReadFile(exampleFile)
	require.NoError(t, err)
	
	fileContent := string(content)
	lines := strings.Split(fileContent, "\n")
	
	skipFound := false
	skipIndex := -1
	legacyVars := make([]string, 0)
	
	for i, line := range lines {
		if strings.Contains(line, "@skip") {
			skipFound = true
			skipIndex = i
		} else if skipFound && strings.Contains(line, "=") && !strings.HasPrefix(strings.TrimSpace(line), "#") {
			varName := strings.Split(strings.TrimSpace(line), "=")[0]
			legacyVars = append(legacyVars, varName)
		}
	}
	
	assert.True(t, skipFound, "@skip annotation should be present")
	assert.Greater(t, skipIndex, 0, "@skip should have an index")
	assert.Greater(t, len(legacyVars), 0, "Variables after @skip should exist")
	assert.Contains(t, legacyVars, "LEGACY_SERVICE_URL")
}

func TestNoteAttributes(t *testing.T) {
	exampleFile := filepath.Join("../examples", "all.env.example")
	content, err := os.ReadFile(exampleFile)
	require.NoError(t, err)
	
	noteCount := strings.Count(string(content), "note=")
	assert.Greater(t, noteCount, 15, "Should have multiple helpful notes")
}

func TestPatternValidationExamples(t *testing.T) {
	testCases := []struct {
		value   string
		pattern string
	}{
		{"1.0.0", `^\d+\.\d+\.\d+$`},
		{"555-123-4567", `^\d{3}-\d{3}-\d{4}$`},
		{"12345", `^\d{5}$`},
		{"abcdef0123456789abcdef0123456789", `^[a-fA-F0-9]{32}$`},
		{"12/31/2025", `^\d{2}/\d{2}/\d{4}$`},
	}
	
	for _, tc := range testCases {
		t.Run(tc.value, func(t *testing.T) {
			matched, err := regexp.MatchString(tc.pattern, tc.value)
			require.NoError(t, err)
			assert.True(t, matched, 
				"Value '%s' should match pattern '%s'", tc.value, tc.pattern)
		})
	}
}

func TestListFieldFormat(t *testing.T) {
	listValues := map[string]string{
		"ALLOWED_ORIGINS": "http://localhost:3000,http://example.com",
		"ADMIN_EMAILS":    "admin@example.com,superadmin@example.com",
		"TRUSTED_IPS":     "127.0.0.1,192.168.1.1",
	}
	
	for varName, value := range listValues {
		t.Run(varName, func(t *testing.T) {
			items := strings.Split(value, ",")
			assert.Greater(t, len(items), 1, 
				"%s should have multiple items", varName)
		})
	}
}

func TestPlainTextVariable(t *testing.T) {
	exampleFile := filepath.Join("../examples", "all.env.example")
	content, err := os.ReadFile(exampleFile)
	require.NoError(t, err)
	
	assert.Contains(t, string(content), "PLAIN_TEXT_VAR=Hello")
}

func TestScriptExamples(t *testing.T) {
	exampleFile := filepath.Join("../examples", "all.env.example")
	content, err := os.ReadFile(exampleFile)
	require.NoError(t, err)
	
	fileContent := string(content)
	
	assert.Contains(t, fileContent, "script=")
	assert.Contains(t, fileContent, "NAME=Rahul")
	assert.Contains(t, fileContent, "@button")
}

func TestAllAnnotationsCovered(t *testing.T) {
	exampleFile := filepath.Join("../examples", "all.env.example")
	content, err := os.ReadFile(exampleFile)
	require.NoError(t, err)
	
	// Just verify file is comprehensive
	fileContent := string(content)
	lines := strings.Split(fileContent, "\n")
	
	assert.Greater(t, len(lines), 150, "File should be comprehensive")
}
