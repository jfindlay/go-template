// Command go-template greets someone from the command line.
//
// It parses arguments with the stdlib flag package, configures structured logging via log/slog,
// and delegates to the gotemplate library.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	gotemplate "github.com/jfindlay/go-template"
)

// version is injected at build time via -ldflags "-X main.version=...". It defaults to "dev" for
// `go run` / `go build` without ldflags. See the Makefile `build` target.
var version = "dev"

// configureLogging sets up slog with a text handler writing to stderr.
//
// Verbose mode sets the level to Debug; otherwise Warn. Configured once here, never reconfigured —
// mirroring python-template _configure_logging and rust-template configure_logging.
func configureLogging(verbose bool) {
	level := slog.LevelWarn
	if verbose {
		level = slog.LevelDebug
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))
}

// run parses args, configures logging, and prints the greeting. It returns an exit code so main
// stays a thin shell and the logic is testable.
func run(args []string, stderr *os.File) int {
	fs := flag.NewFlagSet("go-template", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: go-template [--loud] [-v] <name>\n\n")
		fmt.Fprintf(stderr, "Greet someone from the command line.\n\n")
		fs.PrintDefaults()
	}
	loud := fs.Bool("loud", false, "Upper-case the greeting.")
	var verbose bool
	fs.BoolVar(&verbose, "v", false, "Enable debug logging to stderr.")
	fs.BoolVar(&verbose, "verbose", false, "Enable debug logging to stderr.")
	showVersion := fs.Bool("version", false, "Print version and exit.")

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *showVersion {
		fmt.Println(version)
		return 0
	}

	configureLogging(verbose)

	rest := fs.Args()
	if len(rest) != 1 {
		fs.Usage()
		return 2
	}

	g, err := gotemplate.NewGreeter(rest[0], *loud)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	gotemplate.PrintGreeting(g)
	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stderr))
}
