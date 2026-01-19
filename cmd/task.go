package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/go-cli-tool/internal/task"
)

var (
	taskFile     string
	taskID       string
	taskList     bool
	concurrency  int
	noColor      bool
)

// taskCmd represents the task command
var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Execute automation tasks",
	Long: `Execute automation tasks defined in a configuration file.

The task command allows you to run automated tasks such as:
- Running shell commands
- Executing scripts
- Making HTTP requests
- Batch processing

Tasks can have dependencies, retry logic, and timeouts.`,
	Example: `  # Execute all tasks from a config file
  go-cli-tool task run --file tasks.yaml

  # Execute a specific task
  go-cli-tool task run --file tasks.yaml --id my-task

  # List all tasks in a config file
  go-cli-tool task list --file tasks.yaml

  # Validate task configuration
  go-cli-tool task validate --file tasks.yaml`,
}

// taskRunCmd executes tasks
var taskRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tasks from configuration file",
	Long:  `Execute one or more tasks defined in a YAML configuration file.`,
	RunE:  runTasks,
}

// taskListCmd lists tasks
var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks in configuration file",
	Long:  `Display all tasks defined in the configuration file.`,
	RunE:  listTasks,
}

// taskValidateCmd validates task configuration
var taskValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate task configuration file",
	Long:  `Check if the task configuration file is valid.`,
	RunE:  validateTasks,
}

// taskInitCmd creates a sample task configuration
var taskInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a sample task configuration file",
	Long:  `Generate a sample task configuration file to get started.`,
	RunE:  initTaskConfig,
}

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.AddCommand(taskRunCmd, taskListCmd, taskValidateCmd, taskInitCmd)

	// Flags for task command
	taskCmd.PersistentFlags().StringVarP(&taskFile, "file", "f", "tasks.yaml", "task configuration file")

	// Flags for run command
	taskRunCmd.Flags().StringVar(&taskID, "id", "", "run specific task by ID")
	taskRunCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "number of concurrent tasks")
	taskRunCmd.Flags().BoolVar(&noColor, "no-color", false, "disable colored output")

	// Flags for init command
	taskInitCmd.Flags().BoolVar(&taskList, "example", false, "create file with example tasks")
}

func runTasks(cmd *cobra.Command, args []string) error {
	// Load configuration
	config, err := task.LoadConfig(taskFile)
	if err != nil {
		return fmt.Errorf("❌ Failed to load config: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return fmt.Errorf("❌ Invalid config: %w", err)
	}

	// Create executor
	executor := task.NewExecutor(concurrency, verbose)

	// Add tasks
	if err := executor.AddTasks(config.Tasks); err != nil {
		return fmt.Errorf("failed to add tasks: %w", err)
	}

	fmt.Printf("📋 Loaded %d task(s) from %s\n\n", len(config.Tasks), taskFile)

	// Execute tasks
	ctx := context.Background()
	startTime := time.Now()

	var execErr error
	if taskID != "" {
		// Execute specific task
		fmt.Printf("▶️  Executing task: %s\n\n", taskID)
		_, execErr = executor.ExecuteTask(ctx, taskID)
	} else {
		// Execute all tasks
		fmt.Println("▶️  Executing all tasks...")
		execErr = executor.ExecuteAll(ctx)
	}

	duration := time.Since(startTime)

	// Display results
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 Execution Summary")
	fmt.Println(strings.Repeat("=", 60))

	results := executor.GetResults()
	successCount := 0
	failCount := 0

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Task ID\tStatus\tDuration\tMessage")
	fmt.Fprintln(w, "-------\t------\t--------\t-------")

	for id, result := range results {
		status := "✅ Success"
		message := "Completed"
		if !result.Success {
			status = "❌ Failed"
			failCount++
			if result.Error != nil {
				message = result.Error.Error()
			}
		} else {
			successCount++
		}

		durationStr := fmt.Sprintf("%.2fs", result.Duration.Seconds())
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", id, status, durationStr, message)
	}
	w.Flush()

	fmt.Printf("\nTotal Duration: %.2fs\n", duration.Seconds())
	fmt.Printf("Success: %d\n", successCount)
	fmt.Printf("Failed: %d\n", failCount)

	if execErr != nil {
		return fmt.Errorf("\n⚠️  Execution completed with errors: %w", execErr)
	}

	return nil
}

func listTasks(cmd *cobra.Command, args []string) error {
	config, err := task.LoadConfig(taskFile)
	if err != nil {
		return fmt.Errorf("❌ Failed to load config: %w", err)
	}

	fmt.Printf("📋 Tasks in %s:\n\n", taskFile)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tType\tCommand\tDependencies")
	fmt.Fprintln(w, "--\t----\t----\t-------\t------------")

	for _, t := range config.Tasks {
		deps := "-"
		if len(t.DependsOn) > 0 {
			deps = fmt.Sprintf("%v", t.DependsOn)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", t.ID, t.Name, t.Type, t.Command, deps)
	}
	w.Flush()

	fmt.Printf("\nTotal: %d task(s)\n", len(config.Tasks))
	return nil
}

func validateTasks(cmd *cobra.Command, args []string) error {
	config, err := task.LoadConfig(taskFile)
	if err != nil {
		return fmt.Errorf("❌ Failed to load config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("❌ Invalid config: %w", err)
	}

	fmt.Printf("✅ Configuration file '%s' is valid!\n", taskFile)
	fmt.Printf("   Version: %s\n", config.Version)
	fmt.Printf("   Tasks: %d\n", len(config.Tasks))
	return nil
}

func initTaskConfig(cmd *cobra.Command, args []string) error {
	// Check if file already exists
	if _, err := os.Stat(taskFile); err == nil {
		return fmt.Errorf("file '%s' already exists", taskFile)
	}

	// Create sample configuration
	config := &task.Config{
		Version: "1.0",
		Defaults: task.TaskDefaults{
			Timeout:    30 * time.Second,
			RetryCount: 0,
		},
		Tasks: []*task.Task{},
	}

	if taskList {
		// Add example tasks
		config.Tasks = []*task.Task{
			{
				ID:          "hello",
				Name:        "Hello World",
				Description: "Print a simple greeting",
				Type:        task.TaskTypeCommand,
				Command:     "echo",
				Args:        []string{"Hello from automation!"},
			},
			{
				ID:          "date",
				Name:        "Show Date",
				Description: "Display current date and time",
				Type:        task.TaskTypeCommand,
				Command:     "date",
			},
			{
				ID:          "list-files",
				Name:        "List Files",
				Description: "List files in current directory",
				Type:        task.TaskTypeCommand,
				Command:     "dir",
				DependsOn:   []string{"hello"},
			},
		}
	}

	// Ensure directory exists
	dir := filepath.Dir(taskFile)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Save configuration
	if err := task.SaveConfig(taskFile, config); err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	fmt.Printf("✅ Created task configuration file: %s\n", taskFile)
	if taskList {
		fmt.Println("   Generated with example tasks")
	}
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Edit the file to add your tasks")
	fmt.Println("  2. Run: go-cli-tool task validate -f", taskFile)
	fmt.Println("  3. Run: go-cli-tool task run -f", taskFile)

	return nil
}

