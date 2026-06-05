# Go template agent guide

See [docs/development.md](docs/development.md) for toolchain, setup, and layout.

## Commands

```
make build        # compile (version stamped from git via -ldflags)
make test         # unit + integration tests
make fmt          # auto-format (gofmt -s -w)
make check        # fmt-check + vet + lint + test + cover + vuln
```

Individual targets: `make fmt-check`, `make vet`, `make lint` (staticcheck), `make cover`,
`make vuln` (govulncheck).

## Code conventions

- `gofmt -s` is canonical and owns all code layout. Go has no line-length wrap convention; do not
  impose the python/rust 100-char rule on Go source. (Markdown and comment prose may still wrap.)
- `NewGreeter` is the validation boundary — the Pydantic-validator / `Greeter::new` analog.
  Construct externally sourced data through validating constructors; keep business logic out of them.
- Errors: library returns sentinel errors (`ErrEmptyName`) and wraps I/O failures with
  `fmt.Errorf("...: %w", err)`. Callers use `errors.Is` / `errors.As`. No panics in library code.
- `log/slog` for all structured logging; configured once in `main.configureLogging`, never
  reconfigured.
- Unit tests are table-driven and live in `_test.go` next to the code (idiomatic Go). Integration
  tests build the binary in `TestMain` and run it as a subprocess.
- One third-party dependency by design: `github.com/fatih/color`. Prefer the stdlib for everything
  else (`flag`, `log/slog`, `errors`). `govulncheck` audits the dependency graph.
