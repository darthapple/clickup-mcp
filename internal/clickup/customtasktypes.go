package clickup

import (
	"context"
	"net/http"
)

// ListCustomTaskTypes returns the custom task types defined in a workspace.
// GET /team/{team_id}/custom_item
func (c *Client) ListCustomTaskTypes(ctx context.Context, teamID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/team/" + teamID + "/custom_item"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
