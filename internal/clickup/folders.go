package clickup

import (
	"context"
	"net/http"
	"net/url"
)

// ListFolders returns the folders in a space.
// GET /space/{space_id}/folder
func (c *Client) ListFolders(ctx context.Context, spaceID string, archived bool, archivedSet bool) (any, error) {
	q := url.Values{}
	addBoolParam(q, "archived", archived, archivedSet)
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/space/" + spaceID + "/folder", Query: q}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFolder returns a single folder.
// GET /folder/{folder_id}
func (c *Client) GetFolder(ctx context.Context, folderID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/folder/" + folderID}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateFolder creates a folder in a space.
// POST /space/{space_id}/folder
func (c *Client) CreateFolder(ctx context.Context, spaceID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/space/" + spaceID + "/folder", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateFolder renames a folder.
// PUT /folder/{folder_id}
func (c *Client) UpdateFolder(ctx context.Context, folderID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPut, Path: "/folder/" + folderID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteFolder deletes a folder.
// DELETE /folder/{folder_id}
func (c *Client) DeleteFolder(ctx context.Context, folderID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/folder/" + folderID}, nil)
}

// CreateFolderFromTemplate creates a folder in a space from a folder template.
// POST /space/{space_id}/folder_template/{template_id}
func (c *Client) CreateFolderFromTemplate(ctx context.Context, spaceID, templateID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/space/" + spaceID + "/folder_template/" + templateID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
