# Examples

This document provides practical examples of using the Go CLI Tool.

## Basic Usage

### Getting Help

```bash
# Show general help
go-cli-tool --help

# Show help for a specific command
go-cli-tool hello --help
```

### Version Information

```bash
# Show version
go-cli-tool version

# Show version in root command
go-cli-tool --version
```

## Hello Command Examples

### Simple Greeting

```bash
# Default greeting
$ go-cli-tool hello
Hello, World!

# Custom name
$ go-cli-tool hello --name Alice
Hello, Alice!

# Using short flag
$ go-cli-tool hello -n Bob
Hello, Bob!
```

### Uppercase Output

```bash
# Uppercase greeting
$ go-cli-tool hello --name Charlie --upper
HELLO, CHARLIE!

# Using short flags
$ go-cli-tool hello -n David -u
HELLO, DAVID!
```

### Combining with Other Flags

```bash
# Verbose output with greeting
$ go-cli-tool hello --name Eve --verbose
Using config file: /home/user/.go-cli-tool.yaml
Hello, Eve!
```

## Configuration Examples

### Using Configuration File

Create `~/.go-cli-tool.yaml`:

```yaml
verbose: true
```

Then run:

```bash
$ go-cli-tool hello --name Frank
Using config file: /home/user/.go-cli-tool.yaml
Hello, Frank!
```

### Using Environment Variables

```bash
# Set environment variable (prefix: GOCLITOOL_)
export GOCLITOOL_VERBOSE=true

# Run command
$ go-cli-tool hello --name Grace
Using config file: /home/user/.go-cli-tool.yaml
Hello, Grace!
```

### Custom Configuration File

```bash
# Create custom config
echo "verbose: true" > my-config.yaml

# Use custom config
$ go-cli-tool --config my-config.yaml hello --name Henry
Using config file: my-config.yaml
Hello, Henry!
```

## Scripting Examples

### Bash Script

```bash
#!/bin/bash

# Generate greetings for multiple names
names=("Alice" "Bob" "Charlie")

for name in "${names[@]}"; do
    go-cli-tool hello --name "$name"
done
```

Output:
```
Hello, Alice!
Hello, Bob!
Hello, Charlie!
```

### Error Handling in Scripts

```bash
#!/bin/bash

if go-cli-tool hello --name "$1"; then
    echo "Success!"
else
    echo "Command failed with exit code $?"
    exit 1
fi
```

### Capturing Output

```bash
#!/bin/bash

# Capture output in variable
greeting=$(go-cli-tool hello --name World)
echo "The greeting is: $greeting"

# Process output
go-cli-tool hello --name Alice | tr '[:lower:]' '[:upper:]'
```

## PowerShell Examples

### Basic Usage

```powershell
# Simple greeting
go-cli-tool hello -name Alice

# Multiple commands
$names = @("Alice", "Bob", "Charlie")
foreach ($name in $names) {
    go-cli-tool hello -name $name
}
```

### Error Handling

```powershell
try {
    go-cli-tool hello -name "Test"
    Write-Host "Command succeeded"
} catch {
    Write-Host "Command failed: $_"
    exit 1
}
```

## Advanced Examples

### Creating Aliases

Bash:
```bash
# Add to ~/.bashrc or ~/.bash_profile
alias greet='go-cli-tool hello --name'

# Usage
$ greet Alice
Hello, Alice!
```

PowerShell:
```powershell
# Add to $PROFILE
Set-Alias greet 'go-cli-tool hello -name'

# Usage
PS> greet Alice
Hello, Alice!
```

### Custom Output Formatting

```bash
# JSON-like output (example for future enhancement)
go-cli-tool hello --name Alice --format json

# Table output
go-cli-tool hello --name Alice --format table
```

### Piping and Redirection

```bash
# Save to file
go-cli-tool hello --name Alice > greeting.txt

# Append to file
go-cli-tool hello --name Bob >> greeting.txt

# Pipe to other commands
go-cli-tool hello --name Charlie | grep -i charlie
```

## Integration Examples

### CI/CD Pipeline

GitHub Actions:
```yaml
- name: Run CLI tool
  run: |
    ./go-cli-tool version
    ./go-cli-tool hello --name "CI"
```

GitLab CI:
```yaml
script:
  - ./go-cli-tool version
  - ./go-cli-tool hello --name "CI"
```

### Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o go-cli-tool

FROM alpine:latest
COPY --from=builder /app/go-cli-tool /usr/local/bin/
ENTRYPOINT ["go-cli-tool"]
CMD ["hello"]
```

Run:
```bash
docker run my-cli-tool hello --name Docker
```

### Kubernetes Job

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: cli-greeting
spec:
  template:
    spec:
      containers:
      - name: cli-tool
        image: my-cli-tool:latest
        args: ["hello", "--name", "Kubernetes"]
      restartPolicy: Never
```

## Testing Examples

### Manual Testing

```bash
# Test default behavior
go-cli-tool hello

# Test with various inputs
go-cli-tool hello --name ""
go-cli-tool hello --name "Test User"
go-cli-tool hello --name "测试" # Unicode

# Test flags
go-cli-tool hello -n Test -u -v

# Test error cases
go-cli-tool invalid-command
```

### Automated Testing

```bash
# Run in test script
assert_equals() {
    expected="$1"
    actual="$2"
    if [ "$expected" != "$actual" ]; then
        echo "FAIL: Expected '$expected', got '$actual'"
        exit 1
    fi
}

result=$(go-cli-tool hello --name Test)
assert_equals "Hello, Test!" "$result"
```

## Real-World Scenarios

### Daily Standup Greeter

```bash
#!/bin/bash
hour=$(date +%H)

if [ $hour -lt 12 ]; then
    greeting="Good morning"
elif [ $hour -lt 18 ]; then
    greeting="Good afternoon"
else
    greeting="Good evening"
fi

echo "$greeting, team!"
go-cli-tool hello --name "Team"
```

### Notification System

```bash
#!/bin/bash
message=$(go-cli-tool hello --name "$USER")
notify-send "Welcome" "$message"
```

### Logging System

```bash
#!/bin/bash
timestamp=$(date '+%Y-%m-%d %H:%M:%S')
message=$(go-cli-tool hello --name "$USER")
echo "[$timestamp] $message" >> ~/greetings.log
```

## Tips and Tricks

### 1. Use Short Flags for Quick Commands
```bash
go-cli-tool hello -n Alice -u
```

### 2. Combine with Shell Features
```bash
# Command substitution
echo "Message: $(go-cli-tool hello -n World)"

# Process substitution
diff <(go-cli-tool hello -n A) <(go-cli-tool hello -n B)
```

### 3. Create Wrapper Scripts
```bash
#!/bin/bash
# greet-team.sh
for member in "$@"; do
    go-cli-tool hello --name "$member"
done
```

### 4. Use Configuration for Defaults
Store common settings in config file to avoid repetitive flags.

## Troubleshooting

### Command Not Found
```bash
# Ensure binary is in PATH
export PATH=$PATH:/path/to/go-cli-tool

# Or use full path
/path/to/go-cli-tool hello
```

### Permission Denied
```bash
# Make executable
chmod +x go-cli-tool
```

### Config File Not Loading
```bash
# Check config file location
go-cli-tool hello --verbose

# Specify config explicitly
go-cli-tool --config ~/.go-cli-tool.yaml hello
```
