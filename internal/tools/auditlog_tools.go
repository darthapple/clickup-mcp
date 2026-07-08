package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func RegisterAuditLogTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_create_audit_log_report",
			mcp.WithDescription("Request an audit log export for a ClickUp workspace. Enterprise plans only; non-Enterprise workspaces will get an expected 4xx error."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithNumber("start_date", mcp.Description("Range start, Unix ms timestamp")),
			mcp.WithNumber("end_date", mcp.Description("Range end, Unix ms timestamp")),
			mcp.WithArray("event_types", mcp.WithStringItems(), mcp.Description("Event type names to include")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			body := map[string]any{}
			setFloat(body, req, "start_date")
			setFloat(body, req, "end_date")
			setStringSlice(body, req, "event_types")
			out, err := c.CreateAuditLogReport(ctx, teamIDOrDefault(req, c), body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)
}
