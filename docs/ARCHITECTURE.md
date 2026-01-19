# Architecture

This document describes the architecture and design decisions of the Go CLI Tool.

## Overview

The Go CLI Tool follows a clean architecture pattern with clear separation of concerns:

```
┌─────────────────────────────────────────┐
│           CLI Interface (Cobra)          │
│                                          │
│  ┌──────────┐  ┌──────────┐  ┌────────┐│
│  │  hello   │  │ version  │  │  root  ││
│  └──────────┘  └──────────┘  └────────┘│
└─────────────┬───────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│         Business Logic Layer             │
│                                          │
│  ┌──────────────────────────────────┐   │
│  │   internal/greeter               │   │
│  │   - Core business logic          │   │
│  │   - No external dependencies     │   │
│  └──────────────────────────────────┘   │
└─────────────────────────────────────────┘
```

## Project Structure

### `/cmd`
Contains all CLI command definitions using the Cobra framework.

- **`root.go`**: Defines the root command and global flags
- **`hello.go`**: Hello command implementation
- **`version.go`**: Version information command
- **`*_test.go`**: Unit tests for commands

**Design Principles:**
- Commands are thin wrappers around business logic
- Commands handle input/output and flag parsing
- Business logic is delegated to internal packages

### `/internal`
Contains private application code that cannot be imported by other projects.

- **`greeter/`**: Core greeting functionality
  - Encapsulates greeting logic
  - Fully tested with unit tests
  - No external dependencies

**Design Principles:**
- Pure business logic
- Framework-agnostic
- Highly testable
- Single responsibility

### `/test`
Contains integration tests that test the entire application as a black box.

**Design Principles:**
- Test compiled binary
- No access to internal code
- Real-world usage scenarios

## Design Patterns

### Command Pattern
Using Cobra's command pattern for CLI structure:
- Each command is self-contained
- Easy to add new commands
- Clear command hierarchy

### Dependency Injection
Business logic components are injected:
```go
g := greeter.New()
message := g.Greet(name, uppercase)
```

### Configuration Management
Using Viper for flexible configuration:
- Configuration files
- Environment variables
- Command-line flags
- Default values

Priority order (highest to lowest):
1. Command-line flags
2. Environment variables
3. Configuration file
4. Default values

## Testing Strategy

### Unit Tests
- Test individual functions and methods
- Mock external dependencies
- Use table-driven tests
- Target: >80% coverage

### Integration Tests
- Test compiled binary
- Real command execution
- End-to-end scenarios
- Verify actual behavior

### Benchmarks
- Performance testing
- Identify bottlenecks
- Track performance over time

## Error Handling

Errors are handled at multiple levels:

1. **Business Logic**: Returns errors to callers
2. **Commands**: Handle errors and provide user-friendly messages
3. **Root Command**: Catches unhandled errors and exits with status code

```go
RunE: func(cmd *cobra.Command, args []string) error {
    if err := doSomething(); err != nil {
        return fmt.Errorf("failed to do something: %w", err)
    }
    return nil
}
```

## Configuration

Configuration is managed through Viper with the following sources:

1. **Config File** (`~/.go-cli-tool.yaml`)
2. **Environment Variables** (prefix: `GOCLITOOL_`)
3. **Command Flags**

Example:
```yaml
# ~/.go-cli-tool.yaml
verbose: true
```

## Extension Points

### Adding New Commands

1. Create new file in `/cmd`:
```go
package cmd

import "github.com/spf13/cobra"

var newCmd = &cobra.Command{
    Use:   "newcommand",
    Short: "Description",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

func init() {
    rootCmd.AddCommand(newCmd)
}
```

2. Add tests
3. Update documentation

### Adding Business Logic

1. Create new package in `/internal`
2. Implement logic with tests
3. Use in commands

## Build and Release

### Version Information
Version info is injected at build time:
```bash
go build -ldflags "-X github.com/yourusername/go-cli-tool/cmd.Version=1.0.0"
```

### Cross-Platform Builds
Supports multiple platforms:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

## Dependencies

### Core Dependencies
- **Cobra**: CLI framework
- **Viper**: Configuration management
- **Testify**: Testing assertions

### Dependency Management
- Use Go modules
- Pin versions in `go.mod`
- Regular dependency updates
- Security scanning

## Performance Considerations

- Minimal startup time
- Efficient string operations
- Lazy initialization where appropriate
- Benchmarks for critical paths

## Security

- No sensitive data in logs
- Secure configuration handling
- Input validation
- Regular security updates

## Future Enhancements

Potential areas for expansion:
- Plugin system
- Interactive mode
- Shell completion
- Config validation
- Logging framework
- Metrics/telemetry
