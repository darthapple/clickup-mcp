package clickup

import (
	"context"
	"net/http"
	"net/url"
)

// ListTeamViews returns the views defined at the workspace (Everything) level.
// GET /team/{team_id}/view
func (c *Client) ListTeamViews(ctx context.Context, teamID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/team/" + teamID + "/view"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListSpaceViews returns the views defined on a space.
// GET /space/{space_id}/view
func (c *Client) ListSpaceViews(ctx context.Context, spaceID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/space/" + spaceID + "/view"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateSpaceView creates a view on a space.
// POST /space/{space_id}/view
func (c *Client) CreateSpaceView(ctx context.Context, spaceID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/space/" + spaceID + "/view", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListFolderViews returns the views defined on a folder.
// GET /folder/{folder_id}/view
func (c *Client) ListFolderViews(ctx context.Context, folderID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/folder/" + folderID + "/view"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateFolderView creates a view on a folder.
// POST /folder/{folder_id}/view
func (c *Client) CreateFolderView(ctx context.Context, folderID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/folder/" + folderID + "/view", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListListViews returns the views defined on a list.
// GET /list/{list_id}/view
func (c *Client) ListListViews(ctx context.Context, listID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/list/" + listID + "/view"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateListView creates a view on a list.
// POST /list/{list_id}/view
func (c *Client) CreateListView(ctx context.Context, listID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/list/" + listID + "/view", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetView returns a single view.
// GET /view/{view_id}
func (c *Client) GetView(ctx context.Context, viewID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/view/" + viewID}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateView updates a view.
// PUT /view/{view_id}
func (c *Client) UpdateView(ctx context.Context, viewID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPut, Path: "/view/" + viewID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteView deletes a view.
// DELETE /view/{view_id}
func (c *Client) DeleteView(ctx context.Context, viewID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/view/" + viewID}, nil)
}

// GetViewTasks returns the tasks visible in a view (paginated).
// GET /view/{view_id}/task
func (c *Client) GetViewTasks(ctx context.Context, viewID string, page int, pageSet bool) (any, error) {
	q := url.Values{}
	addIntParam(q, "page", page, pageSet)
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/view/" + viewID + "/task", Query: q}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
