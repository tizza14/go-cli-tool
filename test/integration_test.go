// +build integration

package test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const binaryName = "go-cli-tool"

func TestMain(m *testing.M) {
	// Build the binary before running integration tests
	err := exec.Command("go", "build", "-o", binaryName, "../main.go").Run()
	if err != nil {
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	os.Remove(binaryName)
	os.Exit(code)
}

func TestIntegration_HelloCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  bool
	}{
		{
			name:           "hello with default name",
			args:           []string{"hello"},
			expectedOutput: "Hello, World!",
			expectedError:  false,
		},
		{
			name:           "hello with custom name",
			args:           []string{"hello", "--name", "Integration"},
			expectedOutput: "Hello, Integration!",
			expectedError:  false,
		},
		{
			name:           "hello with uppercase",
			args:           []string{"hello", "-n", "Test", "-u"},
			expectedOutput: "HELLO, TEST!",
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./"+binaryName, tt.args...)
			output, err := cmd.CombinedOutput()

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Contains(t, string(output), tt.expectedOutput)
			}
		})
	}
}

func TestIntegration_VersionCommand(t *testing.T) {
	cmd := exec.Command("./"+binaryName, "version")
	output, err := cmd.CombinedOutput()

	require.NoError(t, err)
	outputStr := string(output)
	assert.Contains(t, outputStr, "Version:")
	assert.Contains(t, outputStr, "Build Date:")
	assert.Contains(t, outputStr, "Git Commit:")
}

func TestIntegration_HelpCommand(t *testing.T) {
	cmd := exec.Command("./"+binaryName, "--help")
	output, err := cmd.CombinedOutput()

	require.NoError(t, err)
	outputStr := string(output)
	assert.Contains(t, outputStr, "A powerful CLI tool built with Go")
	assert.Contains(t, outputStr, "Available Commands:")
	assert.Contains(t, outputStr, "hello")
	assert.Contains(t, outputStr, "version")
}

func TestIntegration_InvalidCommand(t *testing.T) {
	cmd := exec.Command("./"+binaryName, "invalid-command")
	output, err := cmd.CombinedOutput()

	require.Error(t, err)
	outputStr := strings.ToLower(string(output))
	assert.Contains(t, outputStr, "unknown command")
}
