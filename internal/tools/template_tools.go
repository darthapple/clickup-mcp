package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func RegisterTemplateTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_create_folder_from_template",
			mcp.WithDescription("Create a folder in a ClickUp space from a folder template."),
			mcp.WithString("space_id", mcp.Required(), mcp.Description("Space ID")),
			mcp.WithString("template_id", mcp.Required(), mcp.Description("Folder template ID. Cannot be discovered through this MCP server — there is no template-listing tool. Find it in the ClickUp web app (Space/Folder/List settings -> Templates) before calling this tool.")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Name for the new folder")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			spaceID, err := req.RequireString("space_id")
			if err != nil {
				return ErrorResult(err)
			}
			templateID, err := req.RequireString("template_id")
			if err != nil {
				return ErrorResult(err)
			}
			name, err := req.RequireString("name")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.CreateFolderFromTemplate(ctx, spaceID, templateID, map[string]any{"name": name})
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_list_from_template_in_folder",
			mcp.WithDescription("Create a list in a ClickUp folder from a list template."),
			mcp.WithString("folder_id", mcp.Required(), mcp.Description("Folder ID")),
			mcp.WithString("template_id", mcp.Required(), mcp.Description("List template ID. Cannot be discovered through this MCP server — there is no template-listing tool. Find it in the ClickUp web app (Space/Folder/List settings -> Templates) before calling this tool.")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Name for the new list")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			folderID, err := req.RequireString("folder_id")
			if err != nil {
				return ErrorResult(err)
			}
			templateID, err := req.RequireString("template_id")
			if err != nil {
				return ErrorResult(err)
			}
			name, err := req.RequireString("name")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.CreateListFromTemplateInFolder(ctx, folderID, templateID, map[string]any{"name": name})
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_list_from_template_in_space",
			mcp.WithDescription("Create a folderless list in a ClickUp space from a list template."),
			mcp.WithString("space_id", mcp.Required(), mcp.Description("Space ID")),
			mcp.WithString("template_id", mcp.Required(), mcp.Description("List template ID. Cannot be discovered through this MCP server — there is no template-listing tool. Find it in the ClickUp web app (Space/Folder/List settings -> Templates) before calling this tool.")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Name for the new list")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			spaceID, err := req.RequireString("space_id")
			if err != nil {
				return ErrorResult(err)
			}
			templateID, err := req.RequireString("template_id")
			if err != nil {
				return ErrorResult(err)
			}
			name, err := req.RequireString("name")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.CreateListFromTemplateInSpace(ctx, spaceID, templateID, map[string]any{"name": name})
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_task_from_template",
			mcp.WithDescription("Create a task in a ClickUp list from a task template."),
			mcp.WithString("list_id", mcp.Required(), mcp.Description("List ID")),
			mcp.WithString("template_id", mcp.Required(), mcp.Description("Task template ID. Cannot be discovered through this MCP server — there is no template-listing tool. Find it in the ClickUp web app (Space/Folder/List settings -> Templates) before calling this tool.")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Name for the new task")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			listID, err := req.RequireString("list_id")
			if err != nil {
				return ErrorResult(err)
			}
			templateID, err := req.RequireString("template_id")
			if err != nil {
				return ErrorResult(err)
			}
			name, err := req.RequireString("name")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.CreateTaskFromTemplate(ctx, listID, templateID, map[string]any{"name": name})
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)
}
