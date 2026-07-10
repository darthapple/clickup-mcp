package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func buildGuestPermissionBody(req mcp.CallToolRequest) map[string]any {
	body := map[string]any{}
	setFloat(body, req, "permission_level")
	return body
}

var guestPermissionOption = mcp.WithNumber("permission_level", mcp.Description("1=read, 2=comment, 3=edit, 4=create. If omitted, ClickUp applies its own default permission level for this guest — pass explicitly to guarantee the access level."))

// guestEnterpriseNote documents that every guest-management endpoint is
// gated to Enterprise-plan workspaces at the API level, confirmed against
// ClickUp's plan-availability reference (independent of whatever the web UI
// allows) — same phrasing as auditlog_tools.go's equivalent gate.
const guestEnterpriseNote = " Enterprise plan only; non-Enterprise workspaces get an expected 4xx error."

func RegisterGuestTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_invite_guest",
			mcp.WithDescription("Invite a guest to a ClickUp workspace."+guestEnterpriseNote),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("email", mcp.Required(), mcp.Description("Guest email address")),
			mcp.WithBoolean("can_edit_tags", mcp.Description("Allow the guest to edit tags")),
			mcp.WithBoolean("can_see_time_spent", mcp.Description("Allow the guest to see time spent")),
			mcp.WithBoolean("can_see_time_estimated", mcp.Description("Allow the guest to see time estimates")),
			mcp.WithBoolean("can_create_views", mcp.Description("Allow the guest to create views")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			email, err := req.RequireString("email")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{"email": email}
			setBool(body, req, "can_edit_tags")
			setBool(body, req, "can_see_time_spent")
			setBool(body, req, "can_see_time_estimated")
			setBool(body, req, "can_create_views")
			out, err := c.InviteGuest(ctx, teamIDOrDefault(req, c), body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_get_guest",
			mcp.WithDescription("Get a single ClickUp guest by ID."+guestEnterpriseNote+" guest_id is only "+
				"obtainable from clickup_invite_guest's response — there is no tool to "+
				"list or search existing guests."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("guest_id", mcp.Required(), mcp.Description("Guest ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			guestID, err := req.RequireString("guest_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.GetGuest(ctx, teamIDOrDefault(req, c), guestID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_guest",
			mcp.WithDescription("Update a ClickUp guest's workspace-wide permissions."+guestEnterpriseNote),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("guest_id", mcp.Required(), mcp.Description("Guest ID")),
			mcp.WithBoolean("can_edit_tags", mcp.Description("Allow the guest to edit tags")),
			mcp.WithBoolean("can_see_time_spent", mcp.Description("Allow the guest to see time spent")),
			mcp.WithBoolean("can_see_time_estimated", mcp.Description("Allow the guest to see time estimates")),
			mcp.WithBoolean("can_create_views", mcp.Description("Allow the guest to create views")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			guestID, err := req.RequireString("guest_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			setBool(body, req, "can_edit_tags")
			setBool(body, req, "can_see_time_spent")
			setBool(body, req, "can_see_time_estimated")
			setBool(body, req, "can_create_views")
			out, err := c.UpdateGuest(ctx, teamIDOrDefault(req, c), guestID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_remove_guest_from_workspace",
			mcp.WithDescription("Remove a guest from a ClickUp workspace entirely."+guestEnterpriseNote),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("guest_id", mcp.Required(), mcp.Description("Guest ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			guestID, err := req.RequireString("guest_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.RemoveGuestFromWorkspace(ctx, teamIDOrDefault(req, c), guestID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"removed": true, "guest_id": guestID})
		},
	)

	registerGuestScope(s, c, "space", "Space",
		func(ctx context.Context, id, guestID string, body map[string]any) (any, error) {
			return c.AddGuestToSpace(ctx, id, guestID, body)
		},
		func(ctx context.Context, id, guestID string) error { return c.RemoveGuestFromSpace(ctx, id, guestID) },
	)
	registerGuestScope(s, c, "folder", "Folder",
		func(ctx context.Context, id, guestID string, body map[string]any) (any, error) {
			return c.AddGuestToFolder(ctx, id, guestID, body)
		},
		func(ctx context.Context, id, guestID string) error { return c.RemoveGuestFromFolder(ctx, id, guestID) },
	)
	registerGuestScope(s, c, "list", "List",
		func(ctx context.Context, id, guestID string, body map[string]any) (any, error) {
			return c.AddGuestToList(ctx, id, guestID, body)
		},
		func(ctx context.Context, id, guestID string) error { return c.RemoveGuestFromList(ctx, id, guestID) },
	)
	registerGuestScope(s, c, "task", "Task",
		func(ctx context.Context, id, guestID string, body map[string]any) (any, error) {
			return c.AddGuestToTask(ctx, id, guestID, body)
		},
		func(ctx context.Context, id, guestID string) error { return c.RemoveGuestFromTask(ctx, id, guestID) },
	)
}

// registerGuestScope registers the add/remove guest tools for one resource
// scope (space, folder, list, task) using the given client calls.
func registerGuestScope(
	s *server.MCPServer,
	c *clickup.Client,
	scopeParam, scopeLabel string,
	add func(ctx context.Context, id, guestID string, body map[string]any) (any, error),
	remove func(ctx context.Context, id, guestID string) error,
) {
	idParam := scopeParam + "_id"

	s.AddTool(
		mcp.NewTool("clickup_add_guest_to_"+scopeParam,
			mcp.WithDescription("Share a ClickUp "+scopeParam+" with a guest."+guestEnterpriseNote),
			mcp.WithString(idParam, mcp.Required(), mcp.Description(scopeLabel+" ID")),
			mcp.WithString("guest_id", mcp.Required(), mcp.Description("Guest ID")),
			guestPermissionOption,
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, err := req.RequireString(idParam)
			if err != nil {
				return ErrorResult(err)
			}
			guestID, err := req.RequireString("guest_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := add(ctx, id, guestID, buildGuestPermissionBody(req))
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_remove_guest_from_"+scopeParam,
			mcp.WithDescription("Revoke a guest's access to a ClickUp "+scopeParam+"."+guestEnterpriseNote),
			mcp.WithString(idParam, mcp.Required(), mcp.Description(scopeLabel+" ID")),
			mcp.WithString("guest_id", mcp.Required(), mcp.Description("Guest ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, err := req.RequireString(idParam)
			if err != nil {
				return ErrorResult(err)
			}
			guestID, err := req.RequireString("guest_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := remove(ctx, id, guestID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"removed": true, idParam: id, "guest_id": guestID})
		},
	)
}
