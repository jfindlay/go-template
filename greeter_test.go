package gotemplate

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestGreet(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		loud     bool
		expected string
	}{
		{"quiet", "World", false, "Hello, World!"},
		{"loud", "World", true, "HELLO, WORLD!"},
		{"other-name", "Alice", false, "Hello, Alice!"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := NewGreeter(tc.input, tc.loud)
			if err != nil {
				t.Fatalf("NewGreeter(%q, %v) returned error: %v", tc.input, tc.loud, err)
			}
			if got := Greet(g); got != tc.expected {
				t.Errorf("Greet() = %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestNewGreeterLoudDefault(t *testing.T) {
	g, err := NewGreeter("X", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.Loud {
		t.Error("Loud = true, want false")
	}
}

func TestNewGreeterRejectsEmptyName(t *testing.T) {
	cases := []string{"", "   ", "\t\n"}
	for _, name := range cases {
		if _, err := NewGreeter(name, false); !errors.Is(err, ErrEmptyName) {
			t.Errorf("NewGreeter(%q) error = %v, want ErrEmptyName", name, err)
		}
	}
}

func TestGreetToFileWritesGreeting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "greeting.txt")
	g, err := NewGreeter("World", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := GreetToFile(g, path); err != nil {
		t.Fatalf("GreetToFile returned error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if string(data) != "Hello, World!" {
		t.Errorf("file contents = %q, want %q", data, "Hello, World!")
	}
}

func TestGreetToFileErrorsOnBadPath(t *testing.T) {
	g, err := NewGreeter("World", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// A path whose parent directory does not exist forces a write error.
	bad := filepath.Join(t.TempDir(), "nonexistent-dir", "greeting.txt")
	if err := GreetToFile(g, bad); err == nil {
		t.Error("GreetToFile to bad path returned nil, want error")
	}
}

func TestPrintGreeting(t *testing.T) {
	// PrintGreeting writes to stdout; this test exercises both branches for coverage and asserts
	// it does not panic. Output content is asserted in the integration test.
	for _, loud := range []bool{false, true} {
		g, err := NewGreeter("World", loud)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		PrintGreeting(g)
	}
}
