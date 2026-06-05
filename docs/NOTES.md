# Notes

## Why one dependency (fatih/color) and otherwise stdlib

Go's standard library covers CLI parsing (`flag`), structured logging (`log/slog`, since 1.21), and
error wrapping (`errors`/`fmt`) — the affordances that python-template and rust-template reach for
third-party libraries to provide (argparse-vs-clap, structlog-vs-tracing, pydantic-vs-serde). Go
culture prizes a minimal dependency graph, so the template uses the stdlib for all of those and
keeps exactly one third-party dependency, `github.com/fatih/color`, for terminal styling — the
analog of python's rich and rust's owo-colors. Keeping one real dependency (rather than zero) gives
`govulncheck` and `go.sum` a non-trivial graph to audit, demonstrating the supply-chain tooling.

## Supply-chain posture

Two halves, mirroring rust-template's pinning + auditing story:
1. **Pinning** — `go.sum` is committed and records the cryptographic hash of every module version
   in the build graph (direct and transitive). The Go module proxy + checksum database
   (`sum.golang.org`) provide tamper-evidence by default.
2. **Auditing** — `govulncheck` cross-references the build against the official Go vulnerability
   database (`vuln.go.dev`) and, crucially, reports only vulnerabilities in code paths actually
   reachable from this program — lower-noise than a pure version-match advisory check.

## Coverage threshold divergence

python-template and rust-template enforce 100% coverage. Go attributes coverage differently: lines
executed only inside the subprocess spawned by the integration test are not credited to the parent
test run. The `make cover` target therefore measures coverage on the library package only (`.`),
not `./...`. The `cmd/go-template` package is compiled and exercised via the integration tests but
its lines don't appear in the profile.

The library (`greeter.go`) achieves 94.7%. The sole uncovered path is the `fmt.Println` fallback
in `PrintGreeting`, which fires only when `color.Color.Println` returns an error — an essentially
unreachable condition in normal use. `COVER_THRESHOLD` is set to 90 to give honest headroom without
rewarding contorted tests.

## gofmt vs the 100-char rule

The sibling templates wrap at 100 characters. Go deliberately has no line-length setting: `gofmt`
is the canonical, non-configurable formatter, and fighting it is non-idiomatic. The go-template
therefore drops the 100-char convention for Go source (it still applies to Markdown/prose). This is
called out in AGENTS.md and development.md so the divergence is intentional and documented, not an
oversight.
