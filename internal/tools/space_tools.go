package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func teamIDOrDefault(req mcp.CallToolRequest, c *clickup.Client) string {
	return req.GetString("team_id", c.DefaultTeamID())
}

func RegisterSpaceTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_list_spaces",
			mcp.WithDescription("List the spaces in a ClickUp workspace."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithBoolean("archived", mcp.Description("Include archived spaces")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			out, err := c.ListSpaces(ctx, teamIDOrDefault(req, c), req.GetBool("archived", false), hasArg(req, "archived"))
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_get_space",
			mcp.WithDescription("Get a single ClickUp space by ID."),
			mcp.WithString("space_id", mcp.Required(), mcp.Description("Space ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			spaceID, err := req.RequireString("space_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.GetSpace(ctx, spaceID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_space",
			mcp.WithDescription("Create a space in a ClickUp workspace."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Space name")),
			mcp.WithBoolean("multiple_assignees", mcp.Description("Allow multiple assignees on tasks in this space")),
			mcp.WithString("statuses_json", mcp.Description(`JSON array of custom statuses, e.g. [{"status":"to do","color":"#000","type":"open","orderindex":0}]`)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name, err := req.RequireString("name")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{"name": name}
			setBool(body, req, "multiple_assignees")
			if err := setRawJSON(body, req, "statuses_json", "statuses"); err != nil {
				return ErrorResult(err)
			}
			out, err := c.CreateSpace(ctx, teamIDOrDefault(req, c), body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_space",
			mcp.WithDescription("Update a ClickUp space."),
			mcp.WithString("space_id", mcp.Required(), mcp.Description("Space ID")),
			mcp.WithString("name", mcp.Description("New space name")),
			mcp.WithBoolean("multiple_assignees", mcp.Description("Allow multiple assignees on tasks in this space")),
			mcp.WithBoolean("archived", mcp.Description("Archive/unarchive the space")),
			mcp.WithString("statuses_json", mcp.Description(`JSON array of custom statuses, e.g. [{"status":"to do","color":"#000","type":"open","orderindex":0}]`)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			spaceID, err := req.RequireString("space_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			setString(body, req, "name")
			setBool(body, req, "multiple_assignees")
			setBool(body, req, "archived")
			if err := setRawJSON(body, req, "statuses_json", "statuses"); err != nil {
				return ErrorResult(err)
			}
			out, err := c.UpdateSpace(ctx, spaceID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_delete_space",
			mcp.WithDescription("Delete a ClickUp space. This is irreversible."),
			mcp.WithString("space_id", mcp.Required(), mcp.Description("Space ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			spaceID, err := req.RequireString("space_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.DeleteSpace(ctx, spaceID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"deleted": true, "space_id": spaceID})
		},
	)
}
