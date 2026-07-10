# clickup-mcp

[![CI](https://github.com/darthapple/clickup-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/darthapple/clickup-mcp/actions/workflows/ci.yml)
[![Latest release](https://img.shields.io/github/v/release/darthapple/clickup-mcp)](https://github.com/darthapple/clickup-mcp/releases/latest)
[![Go version](https://img.shields.io/github/go-mod/go-version/darthapple/clickup-mcp)](go.mod)
[![License: MIT](https://img.shields.io/github/license/darthapple/clickup-mcp)](LICENSE)

A local MCP server (stdio transport) exposing the ClickUp REST API v2 (plus
the v3 Chat and Docs APIs) as MCP tools, built on
[`mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go).

## Known behavior notes for callers

Every tool's full behavior is documented in its own `description` field
(and each parameter's own `description`) — that's the source of truth an
MCP client/agent actually sees. The items below are cross-cutting patterns
worth knowing before you start, since they each affect several tools at
once and are easy to miss by reading one tool's schema in isolation:

- **Dates/times are plain UTC strings**, not epoch milliseconds:
  `"YYYY-MM-DD HH:MM:SS"` for a precise moment, or bare `"YYYY-MM-DD"`
  where only a calendar date applies (`due_date`/`start_date`). Duration
  fields (`duration`, `duration_ms`) are the one exception and stay in
  milliseconds, since they're an elapsed length, not a point in time.
- **Custom task IDs (`"CT-123"`-style) only work with `clickup_get_task`.**
  Every other task-scoped tool (update, delete, comments, checklists,
  custom fields, dependencies, links, list membership) only accepts
  ClickUp's internal task ID and 404s on a custom one with no hint why —
  resolve it via `clickup_get_task` first.
- **`clickup_update_doc_page` defaults `content_edit_mode` to `"replace"`**
  when omitted, silently overwriting the entire page instead of appending.
  Always pass `content_edit_mode` explicitly when adding to existing
  content.
- **Several tools silently scope down when a filter is omitted**, not
  error: `clickup_list_time_entries` defaults to the last 30 days AND the
  calling token's own user; for a complete or aggregated view use
  `clickup_get_list_time_report`/`clickup_get_user_time_report` instead of
  assembling one by hand.
- **Some list-style tools cap out at one server-side page** with no
  signal more data exists beyond it: the comment-listing tools
  (`clickup_list_task_comments` and friends) and `clickup_search_docs`
  (50-doc cap). `clickup_list_tasks`/`clickup_search_tasks`/
  `clickup_get_view_tasks` do support pagination, but one page per call —
  the caller must re-call with an incrementing `page` itself.
- **Guest management and workspace-user-admin tools are Enterprise-plan
  only** at the ClickUp API level (invite/update/remove workspace user,
  all guest invite/permission tools) — non-Enterprise workspaces get an
  expected 4xx, independent of what the ClickUp web UI otherwise allows.
- **Template IDs (`template_id`, used by the `clickup_create_*_from_template`
  tools) can't be discovered through this server** — there's no
  template-listing tool; find the ID in the ClickUp web app first.

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

Every MCP client below needs the same two things: a path to the `clickup-mcp`
binary (download from [Releases](#releases), or `go build` it yourself), and
`CLICKUP_API_TOKEN`/`CLICKUP_TEAM_ID` passed as environment variables — MCP
clients start this as a subprocess and inject the env vars themselves, so
they don't need to already be set in your shell.

### Claude Code

```sh
claude mcp add clickup --scope user \
  --env CLICKUP_API_TOKEN=pk_xxx \
  --env CLICKUP_TEAM_ID=xxxxxxxx \
  -- /path/to/clickup-mcp
```

`--scope user` makes it available in every project; use `--scope project` to
write it into `.mcp.json` and share it with a team via git instead. Manage
with `claude mcp list` / `claude mcp get clickup` / `claude mcp remove clickup`.

### Claude Desktop

Edit the config file directly — `~/Library/Application Support/Claude/claude_desktop_config.json`
on macOS, `%APPDATA%\Claude\claude_desktop_config.json` on Windows:

```json
{
  "mcpServers": {
    "clickup": {
      "command": "/path/to/clickup-mcp",
      "env": {
        "CLICKUP_API_TOKEN": "pk_xxx",
        "CLICKUP_TEAM_ID": "xxxxxxxx"
      }
    }
  }
}
```

Restart Claude Desktop after editing.

### Cursor / Windsurf / any client using an `mcp.json`-style config

Most other MCP clients (Cursor's `.cursor/mcp.json`, Windsurf's
`~/.codeium/windsurf/mcp_config.json`, etc.) accept the same shape:

```json
{
  "mcpServers": {
    "clickup": {
      "command": "/path/to/clickup-mcp",
      "args": [],
      "env": {
        "CLICKUP_API_TOKEN": "pk_xxx",
        "CLICKUP_TEAM_ID": "xxxxxxxx"
      }
    }
  }
}
```

Check that specific client's docs for the exact config file location and key
name (`mcpServers` vs `servers`, etc.) if it isn't listed here.

### Always running the latest release

MCP clients need a literal filesystem path, not a URL, so "always latest"
means re-fetching the binary to a fixed path rather than changing config:

```sh
curl -sL -o ~/.local/bin/clickup-mcp \
  https://github.com/darthapple/clickup-mcp/releases/latest/download/clickup-mcp-<os>-<arch>
chmod +x ~/.local/bin/clickup-mcp
```

(`<os>`/`<arch>` per the [Releases](#releases) section above, e.g.
`darwin-arm64`, `linux-amd64`.) Point every client above at
`~/.local/bin/clickup-mcp`; re-running that `curl` line whenever you want to
update is all that's needed since the config never has to change.
