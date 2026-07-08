package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func RegisterDependencyTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_add_task_dependency",
			mcp.WithDescription("Add a dependency relationship between two ClickUp tasks. Provide exactly one of depends_on or dependency_of."),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
			mcp.WithString("depends_on", mcp.Description("ID of the task this one depends on (waiting-on relationship)")),
			mcp.WithString("dependency_of", mcp.Description("ID of the task that depends on this one (blocking relationship)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			setString(body, req, "depends_on")
			setString(body, req, "dependency_of")
			if err := c.AddTaskDependency(ctx, taskID, body); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"added": true, "task_id": taskID})
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_remove_task_dependency",
			mcp.WithDescription("Remove a dependency relationship between two ClickUp tasks. Provide the same depends_on/dependency_of pair used to create it."),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
			mcp.WithString("depends_on", mcp.Description("ID of the task this one depends on")),
			mcp.WithString("dependency_of", mcp.Description("ID of the task that depends on this one")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			dependsOn := req.GetString("depends_on", "")
			dependencyOf := req.GetString("dependency_of", "")
			if err := c.RemoveTaskDependency(ctx, taskID, dependsOn, dependencyOf); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"removed": true, "task_id": taskID})
		},
	)
}
