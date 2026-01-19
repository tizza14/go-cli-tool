# Task Automation Guide

## Overview

The Task Automation feature allows you to define and execute automated tasks using YAML configuration files. Tasks can be shell commands, scripts, or HTTP requests, with support for dependencies, retries, and timeouts.

## Features

- ✅ **Multiple Task Types**: Command, Script, HTTP
- ✅ **Task Dependencies**: Execute tasks in order based on dependencies
- ✅ **Retry Logic**: Automatic retry on failure
- ✅ **Timeout Control**: Set maximum execution time per task
- ✅ **Parallel Execution**: Run independent tasks concurrently
- ✅ **Environment Variables**: Pass custom environment variables to tasks
- ✅ **Working Directory**: Set custom working directory per task
- ✅ **YAML Configuration**: Easy-to-read configuration format

## Quick Start

### 1. Create a Task Configuration File

```bash
# Create a sample configuration
go-cli-tool task init -f tasks.yaml --example
```

### 2. Validate the Configuration

```bash
go-cli-tool task validate -f tasks.yaml
```

### 3. List All Tasks

```bash
go-cli-tool task list -f tasks.yaml
```

### 4. Run Tasks

```bash
# Run all tasks
go-cli-tool task run -f tasks.yaml

# Run a specific task
go-cli-tool task run -f tasks.yaml --id my-task

# Run with verbose output
go-cli-tool task run -f tasks.yaml -v

# Run with concurrency
go-cli-tool task run -f tasks.yaml -c 3
```

## Configuration Format

### Basic Structure

```yaml
version: "1.0"

# Default settings applied to all tasks
defaults:
  timeout: 30s
  retry_count: 1
  workdir: "."

# Task definitions
tasks:
  - id: task-1
    name: "Task Name"
    description: "Task description"
    type: command
    command: echo
    args:
      - "Hello World"
```

### Task Properties

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `id` | string | Yes | Unique task identifier |
| `name` | string | Yes | Human-readable task name |
| `description` | string | No | Task description |
| `type` | string | Yes | Task type: `command`, `script`, `http` |
| `command` | string | Yes | Command or script to execute |
| `args` | []string | No | Command arguments |
| `workdir` | string | No | Working directory |
| `env` | map | No | Environment variables |
| `timeout` | duration | No | Maximum execution time |
| `retry_count` | int | No | Number of retries on failure |
| `depends_on` | []string | No | List of task IDs this task depends on |

### Task Types

#### 1. Command Tasks

Execute shell commands:

```yaml
- id: hello
  name: "Hello World"
  type: command
  command: echo
  args:
    - "Hello World"
```

#### 2. Script Tasks

Execute script files:

```yaml
- id: run-script
  name: "Run PowerShell Script"
  type: script
  command: powershell
  args:
    - "-File"
    - "./scripts/deploy.ps1"
  workdir: "./scripts"
```

#### 3. HTTP Tasks (Coming Soon)

Make HTTP requests:

```yaml
- id: api-call
  name: "API Health Check"
  type: http
  command: "GET https://api.example.com/health"
```

## Examples

### Example 1: Simple Sequential Tasks

```yaml
version: "1.0"

tasks:
  - id: step1
    name: "Step 1"
    type: command
    command: echo
    args: ["Starting process..."]

  - id: step2
    name: "Step 2"
    type: command
    command: echo
    args: ["Processing..."]
    depends_on: [step1]

  - id: step3
    name: "Step 3"
    type: command
    command: echo
    args: ["Complete!"]
    depends_on: [step2]
```

### Example 2: Build Pipeline

```yaml
version: "1.0"

defaults:
  timeout: 300s
  workdir: "."

tasks:
  - id: install
    name: "Install Dependencies"
    type: command
    command: go
    args: [mod, download]

  - id: lint
    name: "Lint Code"
    type: command
    command: golangci-lint
    args: [run]
    depends_on: [install]

  - id: test
    name: "Run Tests"
    type: command
    command: go
    args: [test, ./..., -v]
    depends_on: [lint]
    retry_count: 2

  - id: build
    name: "Build Application"
    type: command
    command: go
    args: [build, -o, app.exe]
    depends_on: [test]
```

### Example 3: File Operations

```yaml
version: "1.0"

tasks:
  - id: create-backup-dir
    name: "Create Backup Directory"
    type: command
    command: powershell
    args:
      - "-Command"
      - "New-Item -ItemType Directory -Path './backup' -Force"

  - id: copy-files
    name: "Copy Files"
    type: command
    command: powershell
    args:
      - "-Command"
      - "Copy-Item -Path './data/*' -Destination './backup/' -Recurse"
    depends_on: [create-backup-dir]

  - id: compress
    name: "Compress Backup"
    type: command
    command: powershell
    args:
      - "-Command"
      - "Compress-Archive -Path './backup/*' -DestinationPath 'backup.zip'"
    depends_on: [copy-files]
```

### Example 4: Environment Variables

```yaml
version: "1.0"

tasks:
  - id: deploy
    name: "Deploy Application"
    type: command
    command: deploy.bat
    env:
      ENVIRONMENT: "production"
      API_KEY: "your-api-key"
      DEBUG: "false"
    timeout: 600s
```

## Advanced Features

### Task Dependencies

Tasks can depend on other tasks. Dependencies are executed first:

```yaml
tasks:
  - id: A
    name: "Task A"
    type: command
    command: echo
    args: ["A"]

  - id: B
    name: "Task B"
    type: command
    command: echo
    args: ["B"]
    depends_on: [A]  # B runs after A

  - id: C
    name: "Task C"
    type: command
    command: echo
    args: ["C"]
    depends_on: [A, B]  # C runs after both A and B
```

### Retry Logic

Tasks can automatically retry on failure:

```yaml
- id: unstable-task
  name: "Unstable Task"
  type: command
  command: flaky-script.bat
  retry_count: 3  # Will retry up to 3 times
```

### Timeout Control

Prevent tasks from running too long:

```yaml
- id: long-task
  name: "Long Running Task"
  type: command
  command: long-process.exe
  timeout: 5m  # Will stop after 5 minutes
```

### Concurrent Execution

Run independent tasks in parallel:

```bash
# Run up to 5 tasks concurrently
go-cli-tool task run -f tasks.yaml -c 5
```

## Best Practices

1. **Use Descriptive IDs**: Make task IDs clear and meaningful
2. **Add Descriptions**: Document what each task does
3. **Set Reasonable Timeouts**: Prevent tasks from hanging
4. **Use Dependencies**: Ensure tasks run in the correct order
5. **Test Individually**: Test each task before running the full pipeline
6. **Handle Errors**: Use retry logic for unreliable tasks
7. **Keep Commands Simple**: Break complex operations into multiple tasks
8. **Version Control**: Keep task configurations in git

## Troubleshooting

### Task Fails Immediately

- Check command path and syntax
- Verify working directory exists
- Ensure dependencies are satisfied

### Task Times Out

- Increase timeout value
- Check if process is hanging
- Optimize the command/script

### Dependency Issues

- Validate task configuration
- Check for circular dependencies
- Ensure all dependent tasks exist

### Environment Variables Not Working

- Check variable names (case-sensitive)
- Verify values don't contain special characters
- Test command manually with same environment

## See Also

- [Examples Directory](../examples/)
- [API Documentation](API.md)
- [Architecture Guide](ARCHITECTURE.md)
