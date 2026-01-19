package greeter

import (
	"fmt"
	"strings"
)

// Greeter handles greeting operations
type Greeter struct {
	prefix string
}

// New creates a new Greeter instance
func New() *Greeter {
	return &Greeter{
		prefix: "Hello",
	}
}

// NewWithPrefix creates a new Greeter with a custom prefix
func NewWithPrefix(prefix string) *Greeter {
	return &Greeter{
		prefix: prefix,
	}
}

// Greet generates a greeting message
func (g *Greeter) Greet(name string, uppercase bool) string {
	if name == "" {
		name = "World"
	}

	message := fmt.Sprintf("%s, %s!", g.prefix, name)

	if uppercase {
		message = strings.ToUpper(message)
	}

	return message
}

// SetPrefix updates the greeting prefix
func (g *Greeter) SetPrefix(prefix string) {
	g.prefix = prefix
}

// GetPrefix returns the current greeting prefix
func (g *Greeter) GetPrefix() string {
	return g.prefix
}
