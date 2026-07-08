package clickup

import (
	"context"
	"net/http"
)

// ListCustomRoles returns the custom roles defined in a workspace.
// GET /team/{team_id}/customroles
func (c *Client) ListCustomRoles(ctx context.Context, teamID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/team/" + teamID + "/customroles"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
