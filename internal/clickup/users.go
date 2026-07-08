package clickup

import (
	"context"
	"net/http"
)

// InviteWorkspaceUser invites a user to a workspace.
// POST /team/{team_id}/user
func (c *Client) InviteWorkspaceUser(ctx context.Context, teamID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/team/" + teamID + "/user", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateWorkspaceUser updates a workspace member's role/custom role.
// PUT /team/{team_id}/user/{user_id}
func (c *Client) UpdateWorkspaceUser(ctx context.Context, teamID, userID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPut, Path: "/team/" + teamID + "/user/" + userID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// RemoveWorkspaceUser removes a member from a workspace.
// DELETE /team/{team_id}/user/{user_id}
func (c *Client) RemoveWorkspaceUser(ctx context.Context, teamID, userID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/team/" + teamID + "/user/" + userID}, nil)
}
