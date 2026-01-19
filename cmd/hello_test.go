package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
	}{
		{
			name:           "default hello",
			args:           []string{"hello"},
			expectedOutput: "Hello, World!",
		},
		{
			name:           "hello with name",
			args:           []string{"hello", "--name", "John"},
			expectedOutput: "Hello, John!",
		},
		{
			name:           "hello with name short flag",
			args:           []string{"hello", "-n", "Alice"},
			expectedOutput: "Hello, Alice!",
		},
		{
			name:           "hello uppercase",
			args:           []string{"hello", "--name", "Bob", "--upper"},
			expectedOutput: "HELLO, BOB!",
		},
		{
			name:           "hello uppercase short flags",
			args:           []string{"hello", "-n", "Charlie", "-u"},
			expectedOutput: "HELLO, CHARLIE!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for each test to avoid state leakage
			cmd := rootCmd
			cmd.SetArgs(tt.args)

			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			err := cmd.Execute()
			assert.NoError(t, err)

			output := buf.String()
			assert.Contains(t, output, tt.expectedOutput)
		})
	}
}

func TestHelloCommandHelp(t *testing.T) {
	rootCmd.SetArgs([]string{"hello", "--help"})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	err := rootCmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Print a greeting message")
	assert.Contains(t, output, "--name")
	assert.Contains(t, output, "--upper")
}
