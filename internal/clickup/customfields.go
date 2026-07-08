package clickup

import (
	"context"
	"net/http"
)

// ListListFields returns the custom fields accessible from a list.
// GET /list/{list_id}/field
func (c *Client) ListListFields(ctx context.Context, listID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/list/" + listID + "/field"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListFolderFields returns the custom fields accessible from a folder.
// GET /folder/{folder_id}/field
func (c *Client) ListFolderFields(ctx context.Context, folderID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/folder/" + folderID + "/field"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListSpaceFields returns the custom fields accessible from a space.
// GET /space/{space_id}/field
func (c *Client) ListSpaceFields(ctx context.Context, spaceID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/space/" + spaceID + "/field"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListWorkspaceFields returns every custom field defined in a workspace.
// GET /team/{team_id}/field
func (c *Client) ListWorkspaceFields(ctx context.Context, teamID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/team/" + teamID + "/field"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SetTaskCustomField sets a custom field's value on a task. body: {"value": ...}
// POST /task/{task_id}/field/{field_id}
func (c *Client) SetTaskCustomField(ctx context.Context, taskID, fieldID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/task/" + taskID + "/field/" + fieldID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// RemoveTaskCustomField clears a custom field's value on a task.
// DELETE /task/{task_id}/field/{field_id}
func (c *Client) RemoveTaskCustomField(ctx context.Context, taskID, fieldID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/task/" + taskID + "/field/" + fieldID}, nil)
}
