package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func RegisterAttachmentTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_create_task_attachment",
			mcp.WithDescription("Upload a local file as an attachment on a ClickUp task. file_path must be readable by the server process."),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
			mcp.WithString("file_path", mcp.Required(), mcp.Description("Path to the file to upload")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			filePath, err := req.RequireString("file_path")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.CreateTaskAttachment(ctx, taskID, filePath)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)
}
