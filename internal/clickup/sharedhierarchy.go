package clickup

import (
	"context"
	"net/http"
)

// GetSharedHierarchy returns the tasks/lists/folders shared directly with the
// token's user in a workspace.
// GET /team/{team_id}/shared
func (c *Client) GetSharedHierarchy(ctx context.Context, teamID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/team/" + teamID + "/shared"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
