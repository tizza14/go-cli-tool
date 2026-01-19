package greeter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	g := New()
	require.NotNil(t, g)
	assert.Equal(t, "Hello", g.GetPrefix())
}

func TestNewWithPrefix(t *testing.T) {
	g := NewWithPrefix("Hi")
	require.NotNil(t, g)
	assert.Equal(t, "Hi", g.GetPrefix())
}

func TestGreeter_Greet(t *testing.T) {
	tests := []struct {
		name      string
		greeter   *Greeter
		inputName string
		uppercase bool
		want      string
	}{
		{
			name:      "default greeting",
			greeter:   New(),
			inputName: "World",
			uppercase: false,
			want:      "Hello, World!",
		},
		{
			name:      "custom name",
			greeter:   New(),
			inputName: "John",
			uppercase: false,
			want:      "Hello, John!",
		},
		{
			name:      "uppercase greeting",
			greeter:   New(),
			inputName: "John",
			uppercase: true,
			want:      "HELLO, JOHN!",
		},
		{
			name:      "empty name defaults to World",
			greeter:   New(),
			inputName: "",
			uppercase: false,
			want:      "Hello, World!",
		},
		{
			name:      "custom prefix",
			greeter:   NewWithPrefix("Hey"),
			inputName: "Alice",
			uppercase: false,
			want:      "Hey, Alice!",
		},
		{
			name:      "custom prefix uppercase",
			greeter:   NewWithPrefix("Greetings"),
			inputName: "Bob",
			uppercase: true,
			want:      "GREETINGS, BOB!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.greeter.Greet(tt.inputName, tt.uppercase)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGreeter_SetPrefix(t *testing.T) {
	g := New()
	assert.Equal(t, "Hello", g.GetPrefix())

	g.SetPrefix("Hi")
	assert.Equal(t, "Hi", g.GetPrefix())

	message := g.Greet("Test", false)
	assert.Equal(t, "Hi, Test!", message)
}

func BenchmarkGreeter_Greet(b *testing.B) {
	g := New()
	for i := 0; i < b.N; i++ {
		g.Greet("Benchmark", false)
	}
}

func BenchmarkGreeter_GreetUppercase(b *testing.B) {
	g := New()
	for i := 0; i < b.N; i++ {
		g.Greet("Benchmark", true)
	}
}
