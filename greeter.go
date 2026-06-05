// Package gotemplate provides the core Greeter type and the Greet, GreetToFile, and
// PrintGreeting functions used by the go-template CLI entry point.
package gotemplate

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/fatih/color"
)

// ErrEmptyName is returned by NewGreeter when the supplied name is empty or whitespace-only.
// It is a sentinel error: callers may test for it with errors.Is.
var ErrEmptyName = errors.New("name must not be empty")

// Greeter is a validated configuration for a greeting operation.
//
// Construct it via NewGreeter, which enforces the non-empty-name invariant — the analog of the
// python-template Pydantic field validator and the rust-template Greeter::new boundary. The struct
// tags make it round-trip cleanly through encoding/json, mirroring the Pydantic BaseModel role.
type Greeter struct {
	// Name is the recipient's name.
	Name string `json:"name"`
	// Loud reports whether the greeting is upper-cased.
	Loud bool `json:"loud"`
}

// NewGreeter constructs a validated Greeter.
//
// It returns ErrEmptyName when name is empty or whitespace-only.
func NewGreeter(name string, loud bool) (Greeter, error) {
	if strings.TrimSpace(name) == "" {
		return Greeter{}, ErrEmptyName
	}
	return Greeter{Name: name, Loud: loud}, nil
}

// Greet formats a greeting string from a Greeter.
func Greet(g Greeter) string {
	message := fmt.Sprintf("Hello, %s!", g.Name)
	if g.Loud {
		message = strings.ToUpper(message)
	}
	slog.Debug("greeting_formatted", "name", g.Name, "loud", g.Loud)
	return message
}

// GreetToFile writes a greeting to path, creating or overwriting the file.
//
// It wraps any I/O failure with %w so callers can inspect the underlying error with errors.Is /
// errors.As.
func GreetToFile(g Greeter, path string) error {
	if err := os.WriteFile(path, []byte(Greet(g)), 0o644); err != nil {
		return fmt.Errorf("failed to write greeting to %s: %w", path, err)
	}
	return nil
}

// PrintGreeting prints a colorized greeting to stdout.
//
// Loud greetings render bold red; quiet greetings render bold green — mirroring the
// python-template rich styling and the rust-template owo-colors styling.
func PrintGreeting(g Greeter) {
	message := Greet(g)
	var c *color.Color
	if g.Loud {
		c = color.New(color.FgRed, color.Bold)
	} else {
		c = color.New(color.FgGreen, color.Bold)
	}
	if _, err := c.Println(message); err != nil {
		// Printing is best-effort; fall back to a plain write so output is never lost.
		// This branch is not reachable in tests (color.Println only fails if stdout is
		// broken), so it appears as uncovered in the coverage profile. See docs/NOTES.md.
		fmt.Println(message)
	}
	slog.Debug("greeting_printed", "name", g.Name, "loud", g.Loud)
}
