package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func RegisterGoalTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_list_goals",
			mcp.WithDescription("List the goals in a ClickUp workspace."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			out, err := c.ListGoals(ctx, teamIDOrDefault(req, c))
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_goal",
			mcp.WithDescription("Create a goal in a ClickUp workspace."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Goal name")),
			mcp.WithNumber("due_date", mcp.Description("Due date, Unix ms timestamp")),
			mcp.WithString("description", mcp.Description("Goal description")),
			mcp.WithBoolean("multiple_owners", mcp.Description("Allow multiple owners")),
			mcp.WithArray("owners", mcp.WithStringItems(), mcp.Description("Owner user IDs")),
			mcp.WithString("color", mcp.Description("Goal color")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name, err := req.RequireString("name")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{"name": name}
			setFloat(body, req, "due_date")
			setString(body, req, "description")
			setBool(body, req, "multiple_owners")
			setStringSlice(body, req, "owners")
			setString(body, req, "color")
			out, err := c.CreateGoal(ctx, teamIDOrDefault(req, c), body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_get_goal",
			mcp.WithDescription("Get a single ClickUp goal by ID."),
			mcp.WithString("goal_id", mcp.Required(), mcp.Description("Goal ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			goalID, err := req.RequireString("goal_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.GetGoal(ctx, goalID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_goal",
			mcp.WithDescription("Update a ClickUp goal."),
			mcp.WithString("goal_id", mcp.Required(), mcp.Description("Goal ID")),
			mcp.WithString("name", mcp.Description("Goal name")),
			mcp.WithNumber("due_date", mcp.Description("Due date, Unix ms timestamp")),
			mcp.WithString("description", mcp.Description("Goal description")),
			mcp.WithArray("add_owners", mcp.WithStringItems(), mcp.Description("Owner user IDs to add")),
			mcp.WithArray("rem_owners", mcp.WithStringItems(), mcp.Description("Owner user IDs to remove")),
			mcp.WithString("color", mcp.Description("Goal color")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			goalID, err := req.RequireString("goal_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			setString(body, req, "name")
			setFloat(body, req, "due_date")
			setString(body, req, "description")
			setStringSlice(body, req, "add_owners")
			setStringSlice(body, req, "rem_owners")
			setString(body, req, "color")
			out, err := c.UpdateGoal(ctx, goalID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_delete_goal",
			mcp.WithDescription("Delete a ClickUp goal."),
			mcp.WithString("goal_id", mcp.Required(), mcp.Description("Goal ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			goalID, err := req.RequireString("goal_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.DeleteGoal(ctx, goalID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"deleted": true, "goal_id": goalID})
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_key_result",
			mcp.WithDescription("Create a key result (target) on a ClickUp goal."),
			mcp.WithString("goal_id", mcp.Required(), mcp.Description("Goal ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Key result name")),
			mcp.WithString("type", mcp.Required(), mcp.Description("number, currency, boolean, percentage, or automatic")),
			mcp.WithArray("owners", mcp.WithStringItems(), mcp.Description("Owner user IDs")),
			mcp.WithNumber("steps_start", mcp.Description("Starting value")),
			mcp.WithNumber("steps_end", mcp.Description("Target value")),
			mcp.WithString("unit", mcp.Description("Unit label")),
			mcp.WithArray("task_ids", mcp.WithStringItems(), mcp.Description("Task IDs to track automatically")),
			mcp.WithArray("list_ids", mcp.WithStringItems(), mcp.Description("List IDs to track automatically")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			goalID, err := req.RequireString("goal_id")
			if err != nil {
				return ErrorResult(err)
			}
			name, err := req.RequireString("name")
			if err != nil {
				return ErrorResult(err)
			}
			krType, err := req.RequireString("type")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{"name": name, "type": krType}
			setStringSlice(body, req, "owners")
			setFloat(body, req, "steps_start")
			setFloat(body, req, "steps_end")
			setString(body, req, "unit")
			setStringSlice(body, req, "task_ids")
			setStringSlice(body, req, "list_ids")
			out, err := c.CreateKeyResult(ctx, goalID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_key_result",
			mcp.WithDescription("Update a ClickUp key result's progress."),
			mcp.WithString("key_result_id", mcp.Required(), mcp.Description("Key result ID")),
			mcp.WithNumber("steps_current", mcp.Description("Current value")),
			mcp.WithString("note", mcp.Description("Progress note")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			keyResultID, err := req.RequireString("key_result_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			setFloat(body, req, "steps_current")
			setString(body, req, "note")
			out, err := c.UpdateKeyResult(ctx, keyResultID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_delete_key_result",
			mcp.WithDescription("Delete a ClickUp key result."),
			mcp.WithString("key_result_id", mcp.Required(), mcp.Description("Key result ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			keyResultID, err := req.RequireString("key_result_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.DeleteKeyResult(ctx, keyResultID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"deleted": true, "key_result_id": keyResultID})
		},
	)
}
