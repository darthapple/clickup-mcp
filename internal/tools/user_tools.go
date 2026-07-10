package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func RegisterUserTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_invite_workspace_user",
			mcp.WithDescription("Invite a user to a ClickUp workspace. Enterprise plan only; "+
				"non-Enterprise workspaces get an expected 4xx error."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("email", mcp.Required(), mcp.Description("Invitee email address")),
			mcp.WithBoolean("admin", mcp.Description("Invite as a workspace admin")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			email, err := req.RequireString("email")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{"email": email}
			setBool(body, req, "admin")
			out, err := c.InviteWorkspaceUser(ctx, teamIDOrDefault(req, c), body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_workspace_user",
			mcp.WithDescription("Update a ClickUp workspace member's role. Enterprise plan only; "+
				"non-Enterprise workspaces get an expected 4xx error."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("user_id", mcp.Required(), mcp.Description("User ID")),
			mcp.WithString("username", mcp.Description("Display name")),
			mcp.WithBoolean("admin", mcp.Description("Set true to make the user a workspace admin, false for a regular member")),
			mcp.WithNumber("custom_role_id", mcp.Description("Custom role ID to assign")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			userID, err := req.RequireString("user_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			setString(body, req, "username")
			setBool(body, req, "admin")
			setFloat(body, req, "custom_role_id")
			out, err := c.UpdateWorkspaceUser(ctx, teamIDOrDefault(req, c), userID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_remove_workspace_user",
			mcp.WithDescription("Remove a member from a ClickUp workspace. Enterprise plan only; "+
				"non-Enterprise workspaces get an expected 4xx error."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("user_id", mcp.Required(), mcp.Description("User ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			userID, err := req.RequireString("user_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.RemoveWorkspaceUser(ctx, teamIDOrDefault(req, c), userID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"removed": true, "user_id": userID})
		},
	)
}
