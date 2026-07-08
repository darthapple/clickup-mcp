package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func RegisterCustomFieldTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_list_list_fields",
			mcp.WithDescription("List the custom fields accessible from a ClickUp list."),
			mcp.WithString("list_id", mcp.Required(), mcp.Description("List ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			listID, err := req.RequireString("list_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListListFields(ctx, listID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_folder_fields",
			mcp.WithDescription("List the custom fields accessible from a ClickUp folder."),
			mcp.WithString("folder_id", mcp.Required(), mcp.Description("Folder ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			folderID, err := req.RequireString("folder_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListFolderFields(ctx, folderID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_space_fields",
			mcp.WithDescription("List the custom fields accessible from a ClickUp space."),
			mcp.WithString("space_id", mcp.Required(), mcp.Description("Space ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			spaceID, err := req.RequireString("space_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListSpaceFields(ctx, spaceID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_workspace_fields",
			mcp.WithDescription("List every custom field defined in a ClickUp workspace."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			out, err := c.ListWorkspaceFields(ctx, teamIDOrDefault(req, c))
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_set_task_custom_field",
			mcp.WithDescription("Set a custom field's value on a ClickUp task. The value's required shape depends on the field's type (text, number, dropdown option index, label UUID array, checkbox bool, etc); pass it JSON-encoded in value_json."),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
			mcp.WithString("field_id", mcp.Required(), mcp.Description("Custom field ID")),
			mcp.WithString("value_json", mcp.Required(), mcp.Description(`JSON-encoded field value, e.g. "42", "\"some text\"", or [\"uuid1\",\"uuid2\"]`)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			fieldID, err := req.RequireString("field_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			if err := setRawJSON(body, req, "value_json", "value"); err != nil {
				return ErrorResult(err)
			}
			out, err := c.SetTaskCustomField(ctx, taskID, fieldID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_remove_task_custom_field",
			mcp.WithDescription("Clear a custom field's value on a ClickUp task."),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
			mcp.WithString("field_id", mcp.Required(), mcp.Description("Custom field ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			fieldID, err := req.RequireString("field_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.RemoveTaskCustomField(ctx, taskID, fieldID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"removed": true, "task_id": taskID, "field_id": fieldID})
		},
	)
}
