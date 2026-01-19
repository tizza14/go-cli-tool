# Go CLI Tool

A powerful and well-structured CLI tool built with Go, demonstrating best practices in project organization, testing, and documentation.

## Features

- 🚀 Fast and efficient CLI built with [Cobra](https://github.com/spf13/cobra)
- ⚙️ Configuration management with [Viper](https://github.com/spf13/viper)
- ✅ Comprehensive unit and integration tests
- 📝 Well-documented code and API
- 🔧 Easy to extend with new commands

## Installation

### From Source

```bash
git clone https://github.com/yourusername/go-cli-tool.git
cd go-cli-tool
go build -o go-cli-tool
```

### Using Go Install

```bash
go install github.com/yourusername/go-cli-tool@latest
```

## Usage

### Basic Commands

```bash
# Display help
go-cli-tool --help

# Show version information
go-cli-tool version

# Print a greeting
go-cli-tool hello

# Print a greeting with a custom name
go-cli-tool hello --name John

# Print a greeting in uppercase
go-cli-tool hello --name John --upper
```

### Configuration

The tool supports configuration via:
- Configuration file: `~/.go-cli-tool.yaml` or specify with `--config` flag
- Environment variables (prefixed with `GOCLITOOL_`)
- Command-line flags

Example configuration file:

```yaml
verbose: true
```

### Global Flags

- `-v, --verbose`: Enable verbose output
- `--config`: Specify configuration file path
- `--help`: Display help information
- `--version`: Display version information

## Development

### Prerequisites

- Go 1.21 or higher
- Make (optional, for using Makefile commands)

### Building

```bash
go build -o go-cli-tool
```

### Running Tests

#### Unit Tests

```bash
go test ./... -v
```

#### With Coverage

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

#### Integration Tests

```bash
go test ./test -tags=integration -v
```

#### Run All Tests

```bash
go test ./... -tags=integration -v -coverprofile=coverage.out
```

### Project Structure

```
.
├── cmd/                  # Command implementations
│   ├── root.go          # Root command
│   ├── hello.go         # Hello command
│   ├── version.go       # Version command
│   └── *_test.go        # Command tests
├── internal/            # Internal packages
│   └── greeter/         # Greeter logic
│       ├── greeter.go
│       └── greeter_test.go
├── test/                # Integration tests
│   └── integration_test.go
├── main.go              # Application entry point
├── go.mod               # Go module definition
└── README.md            # This file
```

## Adding New Commands

1. Create a new file in the `cmd` directory (e.g., `cmd/mycommand.go`)
2. Define your command using Cobra:

```go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Brief description",
    Long:  `Detailed description`,
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("My command executed!")
        return nil
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

3. Add tests in `cmd/mycommand_test.go`
4. Update documentation

## Testing

This project emphasizes testing with:

- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test the compiled binary with real commands
- **Benchmarks**: Performance testing for critical paths
- **Table-Driven Tests**: Comprehensive test coverage with multiple scenarios

Test coverage goal: **>80%**

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes
4. Ensure all tests pass (`go test ./... -v`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Testify](https://github.com/stretchr/testify) - Testing toolkit
