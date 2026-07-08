package clickup

import (
	"context"
	"net/http"
	"net/url"
)

// ListSpaces returns the spaces in a workspace.
// GET /team/{team_id}/space
func (c *Client) ListSpaces(ctx context.Context, teamID string, archived bool, archivedSet bool) (any, error) {
	q := url.Values{}
	addBoolParam(q, "archived", archived, archivedSet)
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/team/" + teamID + "/space", Query: q}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateSpace creates a space in a workspace. body is the raw JSON request
// body (name, multiple_assignees, features, statuses, ...).
// POST /team/{team_id}/space
func (c *Client) CreateSpace(ctx context.Context, teamID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/team/" + teamID + "/space", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetSpace returns a single space.
// GET /space/{space_id}
func (c *Client) GetSpace(ctx context.Context, spaceID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/space/" + spaceID}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateSpace updates a space.
// PUT /space/{space_id}
func (c *Client) UpdateSpace(ctx context.Context, spaceID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPut, Path: "/space/" + spaceID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteSpace deletes a space.
// DELETE /space/{space_id}
func (c *Client) DeleteSpace(ctx context.Context, spaceID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/space/" + spaceID}, nil)
}
