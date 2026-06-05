# go-template maintenance tasks. Run `make check` for the full suite.
#
# Two external tools are required for the full suite; install them with:
#   go install honnef.co/go/tools/cmd/staticcheck@latest
#   go install golang.org/x/vuln/cmd/govulncheck@latest

BINARY  := go-template
PKG     := ./...
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X main.version=$(VERSION)
COVER_THRESHOLD := 90

.PHONY: all build fmt fmt-check vet lint test cover vuln check clean

all: check

build: ## Compile the binary with the version stamped in.
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/go-template

fmt: ## Auto-format all Go source (gofmt -s).
	gofmt -s -w .

fmt-check: ## Fail if any file is not gofmt-clean.
	@unformatted=$$(gofmt -s -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "Not gofmt-clean:"; echo "$$unformatted"; exit 1; \
	fi

vet: ## Run go vet (stdlib correctness checks).
	go vet $(PKG)

lint: ## Run staticcheck (the standard external linter).
	staticcheck $(PKG)

test: ## Run unit and integration tests.
	go test $(PKG)

cover: ## Run tests with coverage and enforce the threshold.
	# Cover the library package only: cmd/go-template runs as a subprocess in its integration
	# tests, so its lines are never attributed to the parent coverage profile. See docs/NOTES.md.
	go test -coverprofile=coverage.out -covermode=atomic .
	@total=$$(go tool cover -func=coverage.out | awk '/^total:/ {print $$3}' | tr -d '%'); \
	echo "total coverage: $$total%"; \
	awk -v t="$$total" -v min=$(COVER_THRESHOLD) 'BEGIN { exit (t+0 < min) ? 1 : 0 }' \
		|| { echo "coverage $$total% is below threshold $(COVER_THRESHOLD)%"; exit 1; }

vuln: ## Scan dependencies for known vulnerabilities.
	govulncheck $(PKG)

check: fmt-check vet lint test cover vuln ## Run the full check suite.

clean: ## Remove build and coverage artifacts.
	rm -f $(BINARY) coverage.out
