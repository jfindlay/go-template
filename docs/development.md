# Development

## Toolchain

| Tool | Role | Docs |
|---|---|---|
| go | Build, dependency management, test runner | https://go.dev/doc/ |
| gofmt | Canonical code formatter (`gofmt -s`) | https://pkg.go.dev/cmd/gofmt |
| go vet | Stdlib correctness/suspicious-construct checks | https://pkg.go.dev/cmd/vet |
| staticcheck | The standard external linter | https://staticcheck.dev/ |
| go test | Unit and integration tests, coverage | https://pkg.go.dev/cmd/go#hdr-Test_packages |
| govulncheck | Official vulnerability scanner (supply-chain audit) | https://go.dev/blog/govulncheck |
| make | Task runner for the check suite | https://www.gnu.org/software/make/ |

## Setup

```
git clone https://github.com/jfindlay/go-template
cd go-template
go build ./...
```

Two external dev tools are required for the full check suite:

```
go install honnef.co/go/tools/cmd/staticcheck@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
```

Ensure `$(go env GOPATH)/bin` is on your `PATH` so the installed tools are found.

## Formatting

```
make fmt          # auto-format (gofmt -s -w .)
make fmt-check    # check only; fails if any file is not gofmt-clean
```

Go has no configurable line length — `gofmt` is the single source of truth for layout. This is the
deliberate divergence from python-template and rust-template, both of which wrap at 100 characters.

## Check suite

```
make check
```

Runs, in order: `fmt-check`, `vet`, `lint` (staticcheck), `test`, `cover`, `vuln` (govulncheck).
Individual targets are available (`make vet`, `make lint`, `make cover`, `make vuln`).

## Versioning

The version is derived from git tags at build time and injected via linker flags
(`-ldflags "-X main.version=$(git describe --tags --always --dirty)"`), wired into the Makefile
`build` target. `go-template --version` prints it. To cut a release:

```
git tag -a v0.1.0 -m "v0.1.0"
```

This mirrors python-template's git-tag-derived versioning (hatch-vcs) rather than rust-template's
explicit field. `go install` of a tagged ref records the version in build info automatically.

## Code conventions

- `gofmt -s` owns layout; no manual line wrapping of code.
- `NewGreeter` is the validation boundary (Pydantic-validator / `Greeter::new` analog).
- Library errors: sentinel `ErrEmptyName` + `%w`-wrapped I/O errors; callers use `errors.Is`.
- `log/slog` for logging; configured once in `main.configureLogging`.
- One third-party dep by design: `github.com/fatih/color`. Stdlib for CLI, logging, errors.

## Testing conventions

- Unit tests are table-driven and live in `_test.go` files next to the code they test (idiomatic
  Go placement — the analog of rust's `#[cfg(test)] mod tests`).
- Integration tests live in `cmd/go-template/main_test.go`; `TestMain` builds the binary once and
  the tests run it as a subprocess (the analog of rust's `CARGO_BIN_EXE_` integration test).
- Coverage is measured via `go test -coverprofile` and gated in `make cover`. See docs/NOTES.md for
  why the Go threshold may differ from python/rust's 100%.

## Project layout

```
greeter.go                    package gotemplate: Greeter, NewGreeter, Greet, GreetToFile, PrintGreeting
greeter_test.go               unit tests (table-driven)
cmd/go-template/main.go       CLI entry point: flag parsing and slog setup
cmd/go-template/main_test.go  CLI subprocess integration tests
Makefile                      maintenance task runner
docs/development.md           this file
docs/NOTES.md                 design notes and rationale
```
