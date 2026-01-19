# Contributing to Go CLI Tool

First off, thank you for considering contributing to Go CLI Tool! 🎉

## Code of Conduct

By participating in this project, you are expected to uphold our Code of Conduct:
- Be respectful and inclusive
- Welcome newcomers
- Accept constructive criticism gracefully
- Focus on what is best for the community

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the issue
- **Expected behavior**
- **Actual behavior**
- **Environment details** (OS, Go version, etc.)
- **Code samples** or error messages

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear and descriptive title**
- **Provide detailed description** of the suggested enhancement
- **Explain why this enhancement would be useful**
- **Include examples** of how it would work

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Follow the coding style** of the project
3. **Write tests** for your changes
4. **Update documentation** as needed
5. **Ensure tests pass** before submitting

#### Pull Request Process

1. Update the README.md with details of changes if applicable
2. Add tests for any new functionality
3. Ensure all tests pass: `make test`
4. Run linters: `make lint`
5. Format code: `make fmt`
6. Update CHANGELOG.md with your changes
7. The PR will be merged once you have sign-off from maintainers

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional but recommended)

### Setup Steps

```bash
# Clone your fork
git clone https://github.com/yourusername/go-cli-tool.git
cd go-cli-tool

# Add upstream remote
git remote add upstream https://github.com/originalowner/go-cli-tool.git

# Install dependencies
go mod download

# Run tests to verify setup
make test
```

## Coding Guidelines

### Go Style

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` and `goimports` for formatting
- Run `go vet` to catch common errors
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Testing

- Write unit tests for all new functions
- Aim for >80% code coverage
- Use table-driven tests where appropriate
- Include integration tests for new commands

Example test structure:

```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "test",
            expected: "result",
            wantErr:  false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := MyFunction(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expected, got)
            }
        })
    }
}
```

### Documentation

- Add comments to exported functions and types
- Use godoc format
- Update README.md for user-facing changes
- Add examples where helpful

### Commit Messages

Follow conventional commits format:

```
type(scope): subject

body

footer
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `style`: Code style changes
- `chore`: Maintenance tasks

Example:
```
feat(hello): add uppercase flag

Add --upper flag to hello command to convert
output to uppercase.

Closes #123
```

## Running Tests

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run integration tests
make test-integration

# Run with coverage
make test-coverage

# Run benchmarks
make test-bench
```

## Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install locally
make install
```

## Questions?

Feel free to open an issue with your question or reach out to maintainers.

Thank you for contributing! 🙏
