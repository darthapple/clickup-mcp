package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func RegisterTagTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_add_task_tag",
			mcp.WithDescription("Add an existing space tag to a task."),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
			mcp.WithString("tag_name", mcp.Required(), mcp.Description("Tag name")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			tagName, err := req.RequireString("tag_name")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.AddTaskTag(ctx, taskID, tagName); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"added": true, "task_id": taskID, "tag_name": tagName})
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_remove_task_tag",
			mcp.WithDescription("Remove a tag from a task."),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
			mcp.WithString("tag_name", mcp.Required(), mcp.Description("Tag name")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			tagName, err := req.RequireString("tag_name")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.RemoveTaskTag(ctx, taskID, tagName); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"removed": true, "task_id": taskID, "tag_name": tagName})
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_space_tags",
			mcp.WithDescription("List the tags defined in a ClickUp space."),
			mcp.WithString("space_id", mcp.Required(), mcp.Description("Space ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			spaceID, err := req.RequireString("space_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListSpaceTags(ctx, spaceID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_space_tag",
			mcp.WithDescription("Create a tag in a ClickUp space."),
			mcp.WithString("space_id", mcp.Required(), mcp.Description("Space ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Tag name")),
			mcp.WithString("tag_fg", mcp.Description("Foreground color, e.g. #000000")),
			mcp.WithString("tag_bg", mcp.Description("Background color, e.g. #ffffff")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			spaceID, err := req.RequireString("space_id")
			if err != nil {
				return ErrorResult(err)
			}
			name, err := req.RequireString("name")
			if err != nil {
				return ErrorResult(err)
			}
			tag := map[string]any{"name": name}
			setString(tag, req, "tag_fg")
			setString(tag, req, "tag_bg")
			out, err := c.CreateSpaceTag(ctx, spaceID, map[string]any{"tag": tag})
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_space_tag",
			mcp.WithDescription("Update a tag's name or colors in a ClickUp space."),
			mcp.WithString("space_id", mcp.Required(), mcp.Description("Space ID")),
			mcp.WithString("tag_name", mcp.Required(), mcp.Description("Current tag name")),
			mcp.WithString("name", mcp.Description("New tag name")),
			mcp.WithString("tag_fg", mcp.Description("Foreground color, e.g. #000000")),
			mcp.WithString("tag_bg", mcp.Description("Background color, e.g. #ffffff")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			spaceID, err := req.RequireString("space_id")
			if err != nil {
				return ErrorResult(err)
			}
			tagName, err := req.RequireString("tag_name")
			if err != nil {
				return ErrorResult(err)
			}
			tag := map[string]any{}
			setString(tag, req, "name")
			setString(tag, req, "tag_fg")
			setString(tag, req, "tag_bg")
			out, err := c.UpdateSpaceTag(ctx, spaceID, tagName, map[string]any{"tag": tag})
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_delete_space_tag",
			mcp.WithDescription("Delete a tag from a ClickUp space."),
			mcp.WithString("space_id", mcp.Required(), mcp.Description("Space ID")),
			mcp.WithString("tag_name", mcp.Required(), mcp.Description("Tag name")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			spaceID, err := req.RequireString("space_id")
			if err != nil {
				return ErrorResult(err)
			}
			tagName, err := req.RequireString("tag_name")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.DeleteSpaceTag(ctx, spaceID, tagName); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"deleted": true, "space_id": spaceID, "tag_name": tagName})
		},
	)
}
