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

## Building

```sh
# Host binary, for local development/testing on this machine only:
go build -o bin/clickup-mcp ./cmd/clickup-mcp
```

`bin/` is gitignored.

## Testing

```sh
go test ./...                                  # unit tests, no live API calls
go vet ./...
```

A manual, live smoke test (`cmd/clickup-mcp/smoke_test.go`, build-tagged
`manual`) exercises the server end to end over stdio against the real
ClickUp API:

```sh
go build -o bin/clickup-mcp ./cmd/clickup-mcp
set -a; source ../../.env; set +a
go test -tags manual ./cmd/clickup-mcp/... -run TestSmoke -v
```

It only performs read-only calls (list/get), so it's safe to run against a
real workspace.

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
