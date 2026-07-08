package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func RegisterFolderTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_list_folders",
			mcp.WithDescription("List the folders in a ClickUp space."),
			mcp.WithString("space_id", mcp.Required(), mcp.Description("Space ID")),
			mcp.WithBoolean("archived", mcp.Description("Include archived folders")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			spaceID, err := req.RequireString("space_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListFolders(ctx, spaceID, req.GetBool("archived", false), hasArg(req, "archived"))
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_get_folder",
			mcp.WithDescription("Get a single ClickUp folder by ID."),
			mcp.WithString("folder_id", mcp.Required(), mcp.Description("Folder ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			folderID, err := req.RequireString("folder_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.GetFolder(ctx, folderID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_folder",
			mcp.WithDescription("Create a folder in a ClickUp space."),
			mcp.WithString("space_id", mcp.Required(), mcp.Description("Space ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Folder name")),
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
			out, err := c.CreateFolder(ctx, spaceID, map[string]any{"name": name})
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_folder",
			mcp.WithDescription("Rename a ClickUp folder."),
			mcp.WithString("folder_id", mcp.Required(), mcp.Description("Folder ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("New folder name")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			folderID, err := req.RequireString("folder_id")
			if err != nil {
				return ErrorResult(err)
			}
			name, err := req.RequireString("name")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.UpdateFolder(ctx, folderID, map[string]any{"name": name})
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_delete_folder",
			mcp.WithDescription("Delete a ClickUp folder. This is irreversible."),
			mcp.WithString("folder_id", mcp.Required(), mcp.Description("Folder ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			folderID, err := req.RequireString("folder_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.DeleteFolder(ctx, folderID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"deleted": true, "folder_id": folderID})
		},
	)
}
