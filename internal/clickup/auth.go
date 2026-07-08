package clickup

import (
	"context"
	"net/http"
)

// GetUser returns the user associated with the configured API token.
// GET /user
func (c *Client) GetUser(ctx context.Context) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/user"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListWorkspaces returns the workspaces (teams) the configured token can access.
// GET /team
func (c *Client) ListWorkspaces(ctx context.Context) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/team"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
