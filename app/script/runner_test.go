package script

import (
	"strings"
	"testing"
)

func TestRunAuto(t *testing.T) {
	script := "echo hello"
	outputLines := []string{"VAR1=value1"}
	envVars := map[string]string{"VAR2": "value2"}
	resolved := map[string]string{}
	
	result := RunAuto(script, outputLines, envVars, resolved, "TEST_VAR")
	
	if result != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result)
	}
}

func TestRunAuto_WithEnvVars(t *testing.T) {
	script := "echo $VAR1"
	outputLines := []string{"VAR1=testvalue"}
	envVars := map[string]string{}
	resolved := map[string]string{}
	
	result := RunAuto(script, outputLines, envVars, resolved, "TEST_VAR")
	
	if !strings.Contains(result, "testvalue") {
		t.Errorf("Expected result to contain 'testvalue', got '%s'", result)
	}
}
