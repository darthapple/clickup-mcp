package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func buildViewBody(req mcp.CallToolRequest) (map[string]any, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return nil, err
	}
	body := map[string]any{"name": name}
	setString(body, req, "type")
	return body, nil
}

var viewBodyOptions = []mcp.ToolOption{
	mcp.WithString("name", mcp.Required(), mcp.Description("View name")),
	mcp.WithString("type", mcp.Description("View type: list, board, calendar, gantt, table, timeline, etc.")),
}

func RegisterViewTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_list_team_views",
			mcp.WithDescription("List the workspace-level (Everything) views in a ClickUp workspace."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			out, err := c.ListTeamViews(ctx, teamIDOrDefault(req, c))
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_space_views",
			mcp.WithDescription("List the views defined on a ClickUp space."),
			mcp.WithString("space_id", mcp.Required(), mcp.Description("Space ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			spaceID, err := req.RequireString("space_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListSpaceViews(ctx, spaceID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_space_view",
			append([]mcp.ToolOption{
				mcp.WithDescription("Create a view on a ClickUp space."),
				mcp.WithString("space_id", mcp.Required(), mcp.Description("Space ID")),
			}, viewBodyOptions...)...,
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			spaceID, err := req.RequireString("space_id")
			if err != nil {
				return ErrorResult(err)
			}
			body, err := buildViewBody(req)
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.CreateSpaceView(ctx, spaceID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_folder_views",
			mcp.WithDescription("List the views defined on a ClickUp folder."),
			mcp.WithString("folder_id", mcp.Required(), mcp.Description("Folder ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			folderID, err := req.RequireString("folder_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListFolderViews(ctx, folderID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_folder_view",
			append([]mcp.ToolOption{
				mcp.WithDescription("Create a view on a ClickUp folder."),
				mcp.WithString("folder_id", mcp.Required(), mcp.Description("Folder ID")),
			}, viewBodyOptions...)...,
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			folderID, err := req.RequireString("folder_id")
			if err != nil {
				return ErrorResult(err)
			}
			body, err := buildViewBody(req)
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.CreateFolderView(ctx, folderID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_list_views",
			mcp.WithDescription("List the views defined on a ClickUp list."),
			mcp.WithString("list_id", mcp.Required(), mcp.Description("List ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			listID, err := req.RequireString("list_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListListViews(ctx, listID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_list_view",
			append([]mcp.ToolOption{
				mcp.WithDescription("Create a view on a ClickUp list."),
				mcp.WithString("list_id", mcp.Required(), mcp.Description("List ID")),
			}, viewBodyOptions...)...,
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			listID, err := req.RequireString("list_id")
			if err != nil {
				return ErrorResult(err)
			}
			body, err := buildViewBody(req)
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.CreateListView(ctx, listID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_get_view",
			mcp.WithDescription("Get a single ClickUp view by ID."),
			mcp.WithString("view_id", mcp.Required(), mcp.Description("View ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			viewID, err := req.RequireString("view_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.GetView(ctx, viewID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_view",
			mcp.WithDescription("Update a ClickUp view."),
			mcp.WithString("view_id", mcp.Required(), mcp.Description("View ID")),
			mcp.WithString("name", mcp.Description("View name")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			viewID, err := req.RequireString("view_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			setString(body, req, "name")
			out, err := c.UpdateView(ctx, viewID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_delete_view",
			mcp.WithDescription("Delete a ClickUp view."),
			mcp.WithString("view_id", mcp.Required(), mcp.Description("View ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			viewID, err := req.RequireString("view_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.DeleteView(ctx, viewID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"deleted": true, "view_id": viewID})
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_get_view_tasks",
			mcp.WithDescription("Get the tasks visible in a ClickUp view (paginated)."),
			mcp.WithString("view_id", mcp.Required(), mcp.Description("View ID")),
			mcp.WithInteger("page", mcp.Description("Page number, 0-indexed")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			viewID, err := req.RequireString("view_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.GetViewTasks(ctx, viewID, req.GetInt("page", 0), hasArg(req, "page"))
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)
}
