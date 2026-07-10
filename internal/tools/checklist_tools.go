package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func RegisterChecklistTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_create_checklist",
			mcp.WithDescription("Create a checklist on a ClickUp task."),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID"+taskIDCaveat)),
			mcp.WithString("name", mcp.Required(), mcp.Description("Checklist name")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			name, err := req.RequireString("name")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.CreateChecklist(ctx, taskID, map[string]any{"name": name})
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_checklist",
			mcp.WithDescription("Rename a ClickUp checklist or change its position."),
			mcp.WithString("checklist_id", mcp.Required(), mcp.Description("Checklist ID")),
			mcp.WithString("name", mcp.Description("New checklist name")),
			mcp.WithNumber("position", mcp.Description("New position/order index")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			checklistID, err := req.RequireString("checklist_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			setString(body, req, "name")
			setFloat(body, req, "position")
			out, err := c.UpdateChecklist(ctx, checklistID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_delete_checklist",
			mcp.WithDescription("Delete a ClickUp checklist."),
			mcp.WithString("checklist_id", mcp.Required(), mcp.Description("Checklist ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			checklistID, err := req.RequireString("checklist_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.DeleteChecklist(ctx, checklistID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"deleted": true, "checklist_id": checklistID})
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_checklist_item",
			mcp.WithDescription("Add an item to a ClickUp checklist."),
			mcp.WithString("checklist_id", mcp.Required(), mcp.Description("Checklist ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Item text")),
			mcp.WithString("assignee", mcp.Description("User ID to assign the item to")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			checklistID, err := req.RequireString("checklist_id")
			if err != nil {
				return ErrorResult(err)
			}
			name, err := req.RequireString("name")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{"name": name}
			setString(body, req, "assignee")
			out, err := c.CreateChecklistItem(ctx, checklistID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_checklist_item",
			mcp.WithDescription("Update a ClickUp checklist item."),
			mcp.WithString("checklist_id", mcp.Required(), mcp.Description("Checklist ID")),
			mcp.WithString("checklist_item_id", mcp.Required(), mcp.Description("Checklist item ID")),
			mcp.WithString("name", mcp.Description("New item text")),
			mcp.WithBoolean("resolved", mcp.Description("Mark the item resolved/unresolved")),
			mcp.WithString("assignee", mcp.Description("User ID to assign the item to")),
			mcp.WithString("parent", mcp.Description("Parent checklist item ID, to nest this item")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			checklistID, err := req.RequireString("checklist_id")
			if err != nil {
				return ErrorResult(err)
			}
			itemID, err := req.RequireString("checklist_item_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			setString(body, req, "name")
			setBool(body, req, "resolved")
			setString(body, req, "assignee")
			setString(body, req, "parent")
			out, err := c.UpdateChecklistItem(ctx, checklistID, itemID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_delete_checklist_item",
			mcp.WithDescription("Delete a ClickUp checklist item."),
			mcp.WithString("checklist_id", mcp.Required(), mcp.Description("Checklist ID")),
			mcp.WithString("checklist_item_id", mcp.Required(), mcp.Description("Checklist item ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			checklistID, err := req.RequireString("checklist_id")
			if err != nil {
				return ErrorResult(err)
			}
			itemID, err := req.RequireString("checklist_item_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.DeleteChecklistItem(ctx, checklistID, itemID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"deleted": true, "checklist_id": checklistID, "checklist_item_id": itemID})
		},
	)
}
