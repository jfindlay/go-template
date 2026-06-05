# Go template

A minimal Go project template using the standard toolchain (**go**, **gofmt**, **go vet**) with
**staticcheck** linting, **govulncheck** supply-chain scanning, stdlib **flag** CLI parsing,
**log/slog** structured logging, and **fatih/color** terminal styling. The full maintenance suite
is driven by a **Makefile**.

## Installation

```
go install github.com/jfindlay/go-template/cmd/go-template@latest
```

## Usage

```
go-template [--loud] [-v] <name>
```

| Argument | Description |
|---|---|
| `name` | Name of the person to greet (required). |
| `--loud` | Upper-case the greeting. |
| `-v`, `--verbose` | Enable debug logging to stderr. |
| `--version` | Print the version and exit. |

### Examples

```
$ go-template Alice
Hello, Alice!

$ go-template --loud Bob
HELLO, BOB!
```

## Development

See [docs/development.md](docs/development.md).

## License

GPL-3.0-or-later — see [LICENSE](LICENSE).
