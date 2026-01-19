package task

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// TaskType defines the type of task
type TaskType string

const (
	TaskTypeCommand TaskType = "command" // Execute shell command
	TaskTypeScript  TaskType = "script"  // Execute script file
	TaskTypeHTTP    TaskType = "http"    // Make HTTP request
)

// TaskStatus represents the execution status
type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
	StatusSkipped   TaskStatus = "skipped"
)

// Task represents a single automation task
type Task struct {
	ID          string            `yaml:"id" json:"id"`
	Name        string            `yaml:"name" json:"name"`
	Description string            `yaml:"description" json:"description"`
	Type        TaskType          `yaml:"type" json:"type"`
	Command     string            `yaml:"command" json:"command"`
	Args        []string          `yaml:"args" json:"args"`
	WorkDir     string            `yaml:"workdir" json:"workdir"`
	Env         map[string]string `yaml:"env" json:"env"`
	Timeout     time.Duration     `yaml:"timeout" json:"timeout"`
	RetryCount  int               `yaml:"retry_count" json:"retry_count"`
	DependsOn   []string          `yaml:"depends_on" json:"depends_on"`

	// Runtime fields
	Status    TaskStatus `yaml:"-" json:"status"`
	StartTime time.Time  `yaml:"-" json:"start_time"`
	EndTime   time.Time  `yaml:"-" json:"end_time"`
	Output    string     `yaml:"-" json:"output"`
	Error     string     `yaml:"-" json:"error"`
}

// TaskResult represents the result of task execution
type TaskResult struct {
	Task      *Task
	Success   bool
	Output    string
	Error     error
	Duration  time.Duration
	ExitCode  int
}

// Execute runs the task
func (t *Task) Execute(ctx context.Context) *TaskResult {
	result := &TaskResult{
		Task:    t,
		Success: false,
	}

	t.Status = StatusRunning
	t.StartTime = time.Now()
	defer func() {
		t.EndTime = time.Now()
		result.Duration = t.EndTime.Sub(t.StartTime)
	}()

	// Set timeout if specified
	if t.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, t.Timeout)
		defer cancel()
	}

	// Execute based on task type
	switch t.Type {
	case TaskTypeCommand:
		return t.executeCommand(ctx, result)
	case TaskTypeScript:
		return t.executeScript(ctx, result)
	case TaskTypeHTTP:
		return t.executeHTTP(ctx, result)
	default:
		result.Error = fmt.Errorf("unknown task type: %s", t.Type)
		t.Status = StatusFailed
		t.Error = result.Error.Error()
		return result
	}
}

// executeCommand executes a shell command
func (t *Task) executeCommand(ctx context.Context, result *TaskResult) *TaskResult {
	var cmd *exec.Cmd

	// Combine command and args
	if len(t.Args) > 0 {
		// #nosec G204 -- Command and args are from user-controlled task configuration files
		cmd = exec.CommandContext(ctx, t.Command, t.Args...)
	} else {
		// Parse command string
		parts := strings.Fields(t.Command)
		if len(parts) == 0 {
			result.Error = fmt.Errorf("empty command")
			t.Status = StatusFailed
			t.Error = result.Error.Error()
			return result
		}
		// #nosec G204 -- Command is from user-controlled task configuration files
		cmd = exec.CommandContext(ctx, parts[0], parts[1:]...)
	}

	// Set working directory
	if t.WorkDir != "" {
		cmd.Dir = t.WorkDir
	}

	// Set environment variables
	if len(t.Env) > 0 {
		env := cmd.Environ()
		for k, v := range t.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = env
	}

	// Execute command
	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	t.Output = result.Output

	if err != nil {
		result.Error = err
		t.Status = StatusFailed
		t.Error = err.Error()
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
		return result
	}

	result.Success = true
	result.ExitCode = 0
	t.Status = StatusCompleted
	return result
}

// executeScript executes a script file
func (t *Task) executeScript(ctx context.Context, result *TaskResult) *TaskResult {
	// For now, treat script as command execution
	// Can be extended to support different script interpreters
	return t.executeCommand(ctx, result)
}

// executeHTTP makes an HTTP request
func (t *Task) executeHTTP(ctx context.Context, result *TaskResult) *TaskResult {
	// HTTP execution implementation
	// To be implemented based on requirements
	result.Error = fmt.Errorf("HTTP task type not yet implemented")
	t.Status = StatusFailed
	t.Error = result.Error.Error()
	return result
}

// Validate checks if the task configuration is valid
func (t *Task) Validate() error {
	if t.ID == "" {
		return fmt.Errorf("task ID is required")
	}
	if t.Name == "" {
		return fmt.Errorf("task name is required")
	}
	if t.Type == "" {
		return fmt.Errorf("task type is required")
	}
	if t.Command == "" {
		return fmt.Errorf("task command is required")
	}
	return nil
}

// String returns a string representation of the task
func (t *Task) String() string {
	return fmt.Sprintf("Task[%s: %s]", t.ID, t.Name)
}
