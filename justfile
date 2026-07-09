# List available recipes.
default:
    @just --list

# Build the server binary into bin/.
build:
    go build -o bin/clickup-mcp ./cmd/clickup-mcp

# Run every unit + integration test (no live ClickUp API calls, no secrets).
test:
    go test ./...

# Same as `test`, but with verbose per-test output.
test-verbose:
    go test ./... -v

# Same as `test`, plus the race detector and coverage percentages (what CI runs).
test-race:
    go test ./... -race -cover

# Run only the internal/clickup REST client unit tests.
test-clickup:
    go test ./internal/clickup/...

# Run only the internal/tools MCP tool handler unit tests.
test-tools:
    go test ./internal/tools/...

# Run the integration suite: full internal pipeline, fake ClickUp server.
test-integration:
    go test ./cmd/clickup-mcp/... -run TestIntegration

# Run a single test by name (regex), e.g. `just test-run TestClickupCreateTask`.
test-run name:
    go test ./... -run '{{name}}' -v

# go vet static analysis.
vet:
    go vet ./...

# Check formatting; fails and lists files if any are not gofmt-clean.
fmt-check:
    #!/usr/bin/env bash
    set -euo pipefail
    unformatted=$(gofmt -l .)
    if [ -n "$unformatted" ]; then
        echo "The following files are not gofmt-formatted:"
        echo "$unformatted"
        exit 1
    fi

# Reformat all files with gofmt.
fmt:
    gofmt -w .

# Run everything CI runs: build, vet, fmt-check, race tests.
ci: build vet fmt-check test-race

# Real e2e suite against the live ClickUp API; needs CLICKUP_E2E_SPACE_ID set (see README).
test-e2e: build
    set -a && source .env && set +a && go test -tags e2e ./cmd/clickup-mcp/... -run TestE2E -v

# Fast, read-only sanity check against the real ClickUp API; safe on any workspace.
smoke: build
    set -a && source .env && set +a && go test -tags smoke ./cmd/clickup-mcp/... -run TestSmoke -v
