# API Documentation

This document describes the internal API of the Go CLI Tool.

## Package: `cmd`

The `cmd` package contains all CLI command definitions.

### Root Command

```go
var rootCmd = &cobra.Command{
    Use:   "go-cli-tool",
    Short: "A powerful CLI tool built with Go",
    Long:  "A comprehensive CLI tool...",
}
```

**Global Flags:**
- `--config string`: Configuration file path
- `--verbose, -v`: Enable verbose output
- `--version`: Show version information

### Hello Command

```go
var helloCmd = &cobra.Command{
    Use:   "hello",
    Short: "Print a greeting message",
}
```

**Flags:**
- `--name, -n string`: Name to greet (default: "World")
- `--upper, -u`: Convert message to uppercase

**Examples:**
```bash
go-cli-tool hello
go-cli-tool hello --name John
go-cli-tool hello --name John --upper
```

### Version Command

```go
var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print the version information",
}
```

Displays:
- Version number
- Build date
- Git commit hash

## Package: `internal/greeter`

Business logic for greeting functionality.

### Type: Greeter

```go
type Greeter struct {
    prefix string
}
```

### Functions

#### New

```go
func New() *Greeter
```

Creates a new Greeter with default prefix "Hello".

#### NewWithPrefix

```go
func NewWithPrefix(prefix string) *Greeter
```

Creates a new Greeter with custom prefix.

#### Greet

```go
func (g *Greeter) Greet(name string, uppercase bool) string
```

Generates a greeting message.

**Parameters:**
- `name`: Name to greet (defaults to "World" if empty)
- `uppercase`: Convert to uppercase if true

**Returns:**
- Formatted greeting string

#### SetPrefix

```go
func (g *Greeter) SetPrefix(prefix string)
```

Updates the greeting prefix.

#### GetPrefix

```go
func (g *Greeter) GetPrefix() string
```

Returns the current greeting prefix.
