package clickup

import (
	"context"
	"net/http"
	"net/url"
)

// AddTaskTag adds an existing space tag to a task.
// POST /task/{task_id}/tag/{tag_name}
func (c *Client) AddTaskTag(ctx context.Context, taskID, tagName string) error {
	return c.do(ctx, requestParams{Method: http.MethodPost, Path: "/task/" + taskID + "/tag/" + url.PathEscape(tagName)}, nil)
}

// RemoveTaskTag removes a tag from a task.
// DELETE /task/{task_id}/tag/{tag_name}
func (c *Client) RemoveTaskTag(ctx context.Context, taskID, tagName string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/task/" + taskID + "/tag/" + url.PathEscape(tagName)}, nil)
}

// ListSpaceTags returns the tags defined in a space.
// GET /space/{space_id}/tag
func (c *Client) ListSpaceTags(ctx context.Context, spaceID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/space/" + spaceID + "/tag"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateSpaceTag creates a tag in a space. body: {"tag": {"name":..., "tag_fg":..., "tag_bg":...}}
// POST /space/{space_id}/tag
func (c *Client) CreateSpaceTag(ctx context.Context, spaceID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/space/" + spaceID + "/tag", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateSpaceTag updates a tag's name/colors.
// PUT /space/{space_id}/tag/{tag_name}
func (c *Client) UpdateSpaceTag(ctx context.Context, spaceID, tagName string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPut, Path: "/space/" + spaceID + "/tag/" + url.PathEscape(tagName), Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteSpaceTag deletes a tag from a space.
// DELETE /space/{space_id}/tag/{tag_name}
func (c *Client) DeleteSpaceTag(ctx context.Context, spaceID, tagName string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/space/" + spaceID + "/tag/" + url.PathEscape(tagName)}, nil)
}
