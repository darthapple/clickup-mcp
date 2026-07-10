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
			mcp.WithDescription("Request an audit log export for a ClickUp workspace. Enterprise plans only; non-Enterprise workspaces will get an expected 4xx error. ClickUp retains audit log history for approximately 30 days; a start_date older than that likely will not return data for the missing period."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("start_date", mcp.Description("Range start, UTC \"YYYY-MM-DD HH:MM:SS\" or bare \"YYYY-MM-DD\" (midnight UTC)")),
			mcp.WithString("end_date", mcp.Description("Range end, UTC \"YYYY-MM-DD HH:MM:SS\" or bare \"YYYY-MM-DD\" (midnight UTC)")),
			mcp.WithArray("event_types", mcp.WithStringItems(), mcp.Description("Event type names to include")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			body := map[string]any{}
			if err := setDateTime(body, req, "start_date"); err != nil {
				return ErrorResult(err)
			}
			if err := setDateTime(body, req, "end_date"); err != nil {
				return ErrorResult(err)
			}
			setStringSlice(body, req, "event_types")
			out, err := c.CreateAuditLogReport(ctx, teamIDOrDefault(req, c), body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out, "start_date")
		},
	)
}
