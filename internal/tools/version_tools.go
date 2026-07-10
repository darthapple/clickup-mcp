package tools

import (
	"context"
	"runtime"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

// ServerVersion is this binary's own version. main assigns it from
// main.version (itself set via -ldflags "-X main.version=vX.Y.Z" at release
// build time, "dev" for a local go build) before calling RegisterAll, so
// clickup_get_server_version always reports the real running build rather
// than a hardcoded string.
var ServerVersion = "dev"

// RegisterVersionTools registers a tool that reports this server's own
// build version and the Go runtime it's compiled with. c is unused — this
// isn't a ClickUp API call — but kept for signature parity with every other
// Register*Tools function, since register_test.go's registerFuncs relies on
// a uniform signature to duplicate-check every resource area in one pass.
func RegisterVersionTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_get_server_version",
			mcp.WithDescription("Get this clickup-mcp server's own version and the Go runtime "+
				"it was built with — metadata about the MCP server binary itself, matching "+
				"what `clickup-mcp --version` prints on the command line. Not a ClickUp API "+
				"call and not the ClickUp API's own version."),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return JSONResult(map[string]any{
				"version":    ServerVersion,
				"go_version": runtime.Version(),
			})
		},
	)
}
