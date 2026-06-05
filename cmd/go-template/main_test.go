package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMain builds the binary once and runs the integration tests against it, then cleans up.
// This is the idiomatic Go way to run subprocess CLI tests (analogous to rust-template's
// CARGO_BIN_EXE_ integration test).
var binPath string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "go-template-test")
	if err != nil {
		panic(err)
	}
	binPath = filepath.Join(dir, "go-template")
	build := exec.Command("go", "build", "-o", binPath, ".")
	if out, err := build.CombinedOutput(); err != nil {
		panic("build failed: " + string(out))
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

// runBin executes the test binary with args and returns stdout, stderr, and exit code.
// Named runBin (not run) to avoid shadowing main.go's run function in the same package.
func runBin(t *testing.T, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.Command(binPath, args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	code := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run binary: %v", err)
		}
	}
	return stdout.String(), stderr.String(), code
}

func TestGreetsToStdout(t *testing.T) {
	stdout, _, code := runBin(t, "World")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout, "Hello, World!") {
		t.Errorf("stdout = %q, want it to contain %q", stdout, "Hello, World!")
	}
}

func TestLoudFlagUppercases(t *testing.T) {
	// stdlib flag stops parsing at the first non-flag argument, so --loud must precede the name.
	stdout, _, code := runBin(t, "--loud", "World")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout, "HELLO, WORLD!") {
		t.Errorf("stdout = %q, want it to contain %q", stdout, "HELLO, WORLD!")
	}
}

func TestMissingNameExitsNonzero(t *testing.T) {
	_, _, code := runBin(t)
	if code == 0 {
		t.Error("exit code = 0, want nonzero for missing name")
	}
}
