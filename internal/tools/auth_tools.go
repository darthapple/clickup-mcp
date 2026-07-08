package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func RegisterAuthTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_get_user",
			mcp.WithDescription("Get the ClickUp user associated with the configured API token."),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			out, err := c.GetUser(ctx)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_workspaces",
			mcp.WithDescription("List the ClickUp workspaces (teams) the configured API token can access."),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			out, err := c.ListWorkspaces(ctx)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)
}
