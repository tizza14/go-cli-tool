package task

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the task configuration file
type Config struct {
	Version  string       `yaml:"version"`
	Tasks    []*Task      `yaml:"tasks"`
	Defaults TaskDefaults `yaml:"defaults"`
}

// TaskDefaults contains default values for tasks
type TaskDefaults struct {
	Timeout    time.Duration `yaml:"timeout"`
	RetryCount int           `yaml:"retry_count"`
	WorkDir    string        `yaml:"workdir"`
}

// LoadConfig loads task configuration from a YAML file
func LoadConfig(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults to tasks
	for _, task := range config.Tasks {
		if task.Timeout == 0 && config.Defaults.Timeout > 0 {
			task.Timeout = config.Defaults.Timeout
		}
		if task.RetryCount == 0 && config.Defaults.RetryCount > 0 {
			task.RetryCount = config.Defaults.RetryCount
		}
		if task.WorkDir == "" && config.Defaults.WorkDir != "" {
			task.WorkDir = config.Defaults.WorkDir
		}
		if task.Status == "" {
			task.Status = StatusPending
		}
	}

	return &config, nil
}

// SaveConfig saves task configuration to a YAML file
func SaveConfig(filepath string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Version == "" {
		return fmt.Errorf("config version is required")
	}
	if len(c.Tasks) == 0 {
		return fmt.Errorf("at least one task is required")
	}

	// Validate each task
	taskIDs := make(map[string]bool)
	for _, task := range c.Tasks {
		if err := task.Validate(); err != nil {
			return fmt.Errorf("task %s: %w", task.ID, err)
		}
		if taskIDs[task.ID] {
			return fmt.Errorf("duplicate task ID: %s", task.ID)
		}
		taskIDs[task.ID] = true
	}

	// Validate dependencies
	for _, task := range c.Tasks {
		for _, depID := range task.DependsOn {
			if !taskIDs[depID] {
				return fmt.Errorf("task %s depends on non-existent task: %s", task.ID, depID)
			}
		}
	}

	return nil
}
