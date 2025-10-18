package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEchoScript(t *testing.T) {
	cmd := exec.Command("sh", "-c", "echo 'Hello World'")
	output, err := cmd.Output()
	
	require.NoError(t, err)
	assert.Contains(t, string(output), "Hello")
}

func TestDateScript(t *testing.T) {
	cmd := exec.Command("date")
	output, err := cmd.Output()
	
	require.NoError(t, err)
	assert.Greater(t, len(output), 0)
}

func TestVariableInterpolation(t *testing.T) {
	cmd := exec.Command("sh", "-c", "echo \"Hello $NAME\"")
	cmd.Env = append(os.Environ(), "NAME=TestUser")
	
	output, err := cmd.Output()
	require.NoError(t, err)
	assert.Contains(t, string(output), "TestUser")
}

func TestScriptWithPipes(t *testing.T) {
	cmd := exec.Command("sh", "-c", "echo 'test' | tr a-z A-Z")
	output, err := cmd.Output()
	
	require.NoError(t, err)
	assert.Contains(t, string(output), "TEST")
}

func TestScriptErrorHandling(t *testing.T) {
	cmd := exec.Command("nonexistentcommand123")
	err := cmd.Run()
	
	assert.Error(t, err, "Invalid command should fail")
}

func TestMultipleCommands(t *testing.T) {
	commands := []string{
		"echo 'Line 1'",
		"echo 'Line 2'",
		"echo 'Line 3'",
	}
	
	outputs := make([]string, 0, len(commands))
	for _, cmdStr := range commands {
		cmd := exec.Command("sh", "-c", cmdStr)
		output, err := cmd.Output()
		require.NoError(t, err)
		outputs = append(outputs, strings.TrimSpace(string(output)))
	}
	
	assert.Equal(t, 3, len(outputs))
	assert.Equal(t, "Line 1", outputs[0])
	assert.Equal(t, "Line 3", outputs[2])
}

func TestScriptWithEnvVars(t *testing.T) {
	cmd := exec.Command("sh", "-c", "echo $TEST_VAR")
	cmd.Env = append(os.Environ(), "TEST_VAR=HelloFromEnv")
	
	output, err := cmd.Output()
	require.NoError(t, err)
	assert.Contains(t, string(output), "HelloFromEnv")
}

func TestScriptTimeout(t *testing.T) {
	// Quick script should complete
	cmd := exec.Command("sh", "-c", "echo 'quick'")
	err := cmd.Run()
	assert.NoError(t, err)
}

func TestScriptStdoutCapture(t *testing.T) {
	cmd := exec.Command("sh", "-c", "echo 'stdout message'")
	output, err := cmd.Output()
	
	require.NoError(t, err)
	assert.Equal(t, "stdout message\n", string(output))
}

func TestScriptStderrCapture(t *testing.T) {
	cmd := exec.Command("sh", "-c", "echo 'error message' >&2")
	stderr, err := cmd.CombinedOutput()
	
	require.NoError(t, err)
	assert.Contains(t, string(stderr), "error message")
}

func TestComplexScript(t *testing.T) {
	script := `
		NAME="Rahul"
		echo "Hello $NAME"
	`
	
	cmd := exec.Command("sh", "-c", script)
	output, err := cmd.Output()
	
	require.NoError(t, err)
	assert.Contains(t, string(output), "Hello Rahul")
}
