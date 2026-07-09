# clickup-mcp

A local MCP server (stdio transport) exposing the ClickUp REST API v2 (plus
the v3 Chat and Docs APIs) as MCP tools, built on
[`mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go).

## Environment variables

Required:

- `CLICKUP_API_TOKEN` — a ClickUp personal API token (starts with `pk_`).
- `CLICKUP_TEAM_ID` — your ClickUp workspace ID, used as the default
  `team_id` for tools that need one.

Optional:

- `CLICKUP_API_BASE_URL` (default `https://api.clickup.com/api/v2`)
- `CLICKUP_API_BASE_URL_V3` (default `https://api.clickup.com/api/v3`)
- `CLICKUP_HTTP_TIMEOUT` (default `30s`)
- `CLICKUP_MAX_RETRIES` (default `4`, applies to 429/5xx responses)
- `CLICKUP_E2E_SPACE_ID` — only needed to run the real e2e suite (see
  Testing below). Must point at a ClickUp Space dedicated to test fixtures:
  the suite creates and deletes real Lists/Tasks/Comments/Checklists inside
  it on every run.

## Building

```sh
# Host binary, for local development/testing on this machine only:
go build -o bin/clickup-mcp ./cmd/clickup-mcp
```

`bin/` is gitignored.

## Testing

There are four tiers, from fastest/cheapest to most realistic. A `justfile`
wraps all of them — run `just` (no arguments) to list every recipe.

| Tier | Command | Hits real ClickUp? | What it proves |
| --- | --- | --- | --- |
| Unit | `just test-clickup`, `just test-tools` | No | Each REST client method / MCP tool handler does the right thing in isolation, against a fake `httptest` server whose fixtures we write ourselves. |
| Integration | `just test-integration` | No | The full `RegisterAll -> server.MCPServer -> JSON-RPC -> handler` pipeline is wired together correctly — still against a fake server, so it can't catch a wrong assumption about ClickUp's actual behavior. |
| **E2E** | `just test-e2e` | **Yes** | The real thing: create/get/update/delete lifecycles run against the live API in a disposable sandbox List (created and torn down by the suite itself), with assertions against ClickUp's actual responses — not fixtures we invented. This is the tier that answers "how do we know our assumptions about ClickUp are correct." |
| Smoke | `just smoke` | Yes | Fast, read-only pre-flight check ("can we auth and list things at all") safe to run against *any* real workspace with zero setup. Shallower than e2e by design. |

Plain `go test ./...` (or `go vet ./...`, `gofmt -l .`) runs unit +
integration only — no secrets, no network, and what CI
(`.github/workflows/ci.yml`) runs on every push/PR.

**E2E** (`cmd/clickup-mcp/e2e_test.go`, build tag `e2e`) and **smoke**
(`cmd/clickup-mcp/smoke_test.go`, build tag `smoke`) both spawn the compiled
binary over real stdio and call it against the real ClickUp API, so they're
excluded from CI (no ClickUp secrets are stored there) and run locally only:

```sh
just test-e2e   # needs CLICKUP_API_TOKEN, CLICKUP_TEAM_ID, CLICKUP_E2E_SPACE_ID
just smoke      # needs CLICKUP_API_TOKEN, CLICKUP_TEAM_ID only
```

E2E mutates real data (inside its own disposable sandbox List only) — never
point `CLICKUP_E2E_SPACE_ID` at a Space you use for real work. Smoke never
creates/updates/deletes anything, so it's safe against your main workspace.

## Releases

Every push to `main` that passes CI is checked for
[Conventional Commits](https://www.conventionalcommits.org/) since the last
release tag (`.github/workflows/release.yml`). If it finds any `feat:`/`fix:`/
breaking-change (`!:` or a `BREAKING CHANGE:` footer) commit subjects, it:

1. Computes the next semver tag (breaking → major, `feat:` → minor, `fix:` →
   patch) and pushes it.
2. Creates a GitHub Release with auto-generated notes.
3. Cross-compiles the binary for `linux/amd64`, `linux/arm64`,
   `darwin/amd64`, `darwin/arm64`, and `windows/amd64`, and uploads each as a
   raw (uncompressed) asset — no Go toolchain needed to consume it.

No matching commits since the last tag → no release; the workflow just exits.

Each binary is uploaded under two names:

- `clickup-mcp-vX.Y.Z-<os>-<arch>` — pinned to that exact release.
- `clickup-mcp-<os>-<arch>` — re-uploaded (overwritten) on every release, so
  it always reflects whatever GitHub currently marks as the **latest**
  release.

To always fetch the newest build for a platform, without needing to know the
current version:

```sh
curl -L -o clickup-mcp \
  https://github.com/darthapple/clickup-mcp/releases/latest/download/clickup-mcp-linux-amd64
chmod +x clickup-mcp
```

To pin to an exact version instead, use the versioned asset name under that
release's tag:

```sh
curl -L -o clickup-mcp \
  https://github.com/darthapple/clickup-mcp/releases/download/v1.2.3/clickup-mcp-v1.2.3-linux-amd64
chmod +x clickup-mcp
```

(Windows assets have a `.exe` suffix, e.g. `clickup-mcp-windows-amd64.exe`.)

`clickup-mcp --version` prints the running binary's version — set at build
time via `-ldflags "-X main.version=vX.Y.Z"`; local `go build` without that
flag reports `dev`.

## Wiring into an MCP client

Point any MCP-compatible client at the compiled binary as a local stdio
subprocess, e.g. in an `mcp.json`:

```json
"clickup": {
  "command": "/path/to/clickup-mcp",
  "args": []
}
```

`CLICKUP_API_TOKEN`/`CLICKUP_TEAM_ID` reach the subprocess via inherited
environment variables, so make sure they're set in the parent process's
environment (e.g. via `--env-file .env` if the client runs in a container).
