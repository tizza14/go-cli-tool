package task

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Executor manages and executes tasks
type Executor struct {
	tasks       map[string]*Task
	results     map[string]*TaskResult
	mu          sync.RWMutex
	concurrency int
	verbose     bool
}

// NewExecutor creates a new task executor
func NewExecutor(concurrency int, verbose bool) *Executor {
	if concurrency <= 0 {
		concurrency = 1
	}
	return &Executor{
		tasks:       make(map[string]*Task),
		results:     make(map[string]*TaskResult),
		concurrency: concurrency,
		verbose:     verbose,
	}
}

// AddTask adds a task to the executor
func (e *Executor) AddTask(task *Task) error {
	if err := task.Validate(); err != nil {
		return fmt.Errorf("invalid task: %w", err)
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.tasks[task.ID]; exists {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}

	e.tasks[task.ID] = task
	return nil
}

// AddTasks adds multiple tasks
func (e *Executor) AddTasks(tasks []*Task) error {
	for _, task := range tasks {
		if err := e.AddTask(task); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteAll executes all tasks respecting dependencies
func (e *Executor) ExecuteAll(ctx context.Context) error {
	e.mu.RLock()
	taskCount := len(e.tasks)
	e.mu.RUnlock()

	if taskCount == 0 {
		return fmt.Errorf("no tasks to execute")
	}

	// Build execution order based on dependencies
	executionOrder, err := e.buildExecutionOrder()
	if err != nil {
		return fmt.Errorf("failed to build execution order: %w", err)
	}

	// Execute tasks in order
	for _, taskID := range executionOrder {
		e.mu.RLock()
		task := e.tasks[taskID]
		e.mu.RUnlock()

		// Check if dependencies succeeded
		if !e.checkDependencies(task) {
			task.Status = StatusSkipped
			e.mu.Lock()
			e.results[taskID] = &TaskResult{
				Task:    task,
				Success: false,
				Error:   fmt.Errorf("dependencies failed"),
			}
			e.mu.Unlock()
			continue
		}

		// Execute task with retry
		result := e.executeWithRetry(ctx, task)
		
		e.mu.Lock()
		e.results[taskID] = result
		e.mu.Unlock()

		// Stop if task failed and no retry
		if !result.Success && task.RetryCount == 0 {
			return fmt.Errorf("task %s failed: %w", taskID, result.Error)
		}
	}

	return nil
}

// ExecuteTask executes a specific task by ID
func (e *Executor) ExecuteTask(ctx context.Context, taskID string) (*TaskResult, error) {
	e.mu.RLock()
	task, exists := e.tasks[taskID]
	e.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("task %s not found", taskID)
	}

	result := e.executeWithRetry(ctx, task)
	
	e.mu.Lock()
	e.results[taskID] = result
	e.mu.Unlock()

	return result, nil
}

// executeWithRetry executes a task with retry logic
func (e *Executor) executeWithRetry(ctx context.Context, task *Task) *TaskResult {
	var result *TaskResult
	maxAttempts := task.RetryCount + 1

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if e.verbose {
			fmt.Printf("[%s] Executing task %s (attempt %d/%d)...\n", 
				time.Now().Format("15:04:05"), task.Name, attempt, maxAttempts)
		}

		result = task.Execute(ctx)

		if result.Success {
			if e.verbose {
				fmt.Printf("[%s] Task %s completed successfully (%.2fs)\n", 
					time.Now().Format("15:04:05"), task.Name, result.Duration.Seconds())
			}
			return result
		}

		if attempt < maxAttempts {
			if e.verbose {
				fmt.Printf("[%s] Task %s failed, retrying... (%v)\n", 
					time.Now().Format("15:04:05"), task.Name, result.Error)
			}
			time.Sleep(time.Second * 2) // Wait before retry
		}
	}

	if e.verbose {
		fmt.Printf("[%s] Task %s failed after %d attempts\n", 
			time.Now().Format("15:04:05"), task.Name, maxAttempts)
	}

	return result
}

// checkDependencies checks if all dependencies completed successfully
func (e *Executor) checkDependencies(task *Task) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, depID := range task.DependsOn {
		result, exists := e.results[depID]
		if !exists || !result.Success {
			return false
		}
	}
	return true
}

// buildExecutionOrder builds task execution order based on dependencies
func (e *Executor) buildExecutionOrder() ([]string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	visited := make(map[string]bool)
	order := []string{}

	var visit func(taskID string) error
	visit = func(taskID string) error {
		if visited[taskID] {
			return nil
		}

		task, exists := e.tasks[taskID]
		if !exists {
			return fmt.Errorf("task %s not found", taskID)
		}

		// Visit dependencies first
		for _, depID := range task.DependsOn {
			if _, exists := e.tasks[depID]; !exists {
				return fmt.Errorf("dependency %s not found for task %s", depID, taskID)
			}
			if err := visit(depID); err != nil {
				return err
			}
		}

		visited[taskID] = true
		order = append(order, taskID)
		return nil
	}

	// Visit all tasks
	for taskID := range e.tasks {
		if err := visit(taskID); err != nil {
			return nil, err
		}
	}

	return order, nil
}

// GetResults returns all task results
func (e *Executor) GetResults() map[string]*TaskResult {
	e.mu.RLock()
	defer e.mu.RUnlock()

	results := make(map[string]*TaskResult)
	for k, v := range e.results {
		results[k] = v
	}
	return results
}

// GetResult returns a specific task result
func (e *Executor) GetResult(taskID string) (*TaskResult, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	result, exists := e.results[taskID]
	return result, exists
}

// Clear clears all tasks and results
func (e *Executor) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.tasks = make(map[string]*Task)
	e.results = make(map[string]*TaskResult)
}
