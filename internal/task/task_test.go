package task

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTask_Validate(t *testing.T) {
	tests := []struct {
		name    string
		task    *Task
		wantErr bool
	}{
		{
			name: "valid task",
			task: &Task{
				ID:      "test-1",
				Name:    "Test Task",
				Type:    TaskTypeCommand,
				Command: "echo hello",
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			task: &Task{
				Name:    "Test Task",
				Type:    TaskTypeCommand,
				Command: "echo hello",
			},
			wantErr: true,
		},
		{
			name: "missing name",
			task: &Task{
				ID:      "test-1",
				Type:    TaskTypeCommand,
				Command: "echo hello",
			},
			wantErr: true,
		},
		{
			name: "missing type",
			task: &Task{
				ID:      "test-1",
				Name:    "Test Task",
				Command: "echo hello",
			},
			wantErr: true,
		},
		{
			name: "missing command",
			task: &Task{
				ID:   "test-1",
				Name: "Test Task",
				Type: TaskTypeCommand,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTask_Execute(t *testing.T) {
	tests := []struct {
		name        string
		task        *Task
		wantSuccess bool
		wantOutput  string
	}{
		{
			name: "successful command",
			task: &Task{
				ID:      "test-1",
				Name:    "Echo Test",
				Type:    TaskTypeCommand,
				Command: "powershell",
				Args:    []string{"-Command", "Write-Host hello"},
			},
			wantSuccess: true,
			wantOutput:  "hello",
		},
		{
			name: "failed command",
			task: &Task{
				ID:      "test-2",
				Name:    "Invalid Command",
				Type:    TaskTypeCommand,
				Command: "nonexistent-command-12345",
			},
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result := tt.task.Execute(ctx)

			assert.Equal(t, tt.wantSuccess, result.Success)
			if tt.wantOutput != "" {
				assert.Contains(t, result.Output, tt.wantOutput)
			}
			assert.NotZero(t, result.Duration)
		})
	}
}

func TestTask_ExecuteWithTimeout(t *testing.T) {
	task := &Task{
		ID:      "test-timeout",
		Name:    "Timeout Test",
		Type:    TaskTypeCommand,
		Command: "ping",
		Args:    []string{"127.0.0.1", "-n", "10"},
		Timeout: 500 * time.Millisecond,
	}

	ctx := context.Background()
	result := task.Execute(ctx)

	assert.False(t, result.Success)
	assert.NotNil(t, result.Error)
}

func TestExecutor_AddTask(t *testing.T) {
	executor := NewExecutor(1, false)

	task := &Task{
		ID:      "test-1",
		Name:    "Test Task",
		Type:    TaskTypeCommand,
		Command: "echo hello",
	}

	err := executor.AddTask(task)
	require.NoError(t, err)

	// Try adding duplicate
	err = executor.AddTask(task)
	assert.Error(t, err)
}

func TestExecutor_ExecuteTask(t *testing.T) {
	executor := NewExecutor(1, false)

	task := &Task{
		ID:      "test-1",
		Name:    "Echo Test",
		Type:    TaskTypeCommand,
		Command: "powershell",
		Args:    []string{"-Command", "Write-Host hello"},
	}

	err := executor.AddTask(task)
	require.NoError(t, err)

	ctx := context.Background()
	result, err := executor.ExecuteTask(ctx, "test-1")

	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Contains(t, result.Output, "hello")
}

func TestExecutor_ExecuteAll(t *testing.T) {
	executor := NewExecutor(1, false)

	tasks := []*Task{
		{
			ID:      "task-1",
			Name:    "First Task",
			Type:    TaskTypeCommand,
			Command: "powershell",
			Args:    []string{"-Command", "Write-Host first"},
		},
		{
			ID:        "task-2",
			Name:      "Second Task",
			Type:      TaskTypeCommand,
			Command:   "powershell",
			Args:      []string{"-Command", "Write-Host second"},
			DependsOn: []string{"task-1"},
		},
	}

	err := executor.AddTasks(tasks)
	require.NoError(t, err)

	ctx := context.Background()
	err = executor.ExecuteAll(ctx)
	require.NoError(t, err)

	results := executor.GetResults()
	assert.Len(t, results, 2)
	assert.True(t, results["task-1"].Success)
	assert.True(t, results["task-2"].Success)
}

func TestExecutor_ExecuteWithRetry(t *testing.T) {
	executor := NewExecutor(1, false)

	task := &Task{
		ID:         "test-retry",
		Name:       "Retry Test",
		Type:       TaskTypeCommand,
		Command:    "nonexistent-command",
		RetryCount: 2,
	}

	err := executor.AddTask(task)
	require.NoError(t, err)

	ctx := context.Background()
	result, _ := executor.ExecuteTask(ctx, "test-retry")

	assert.False(t, result.Success)
	assert.NotNil(t, result.Error)
}
