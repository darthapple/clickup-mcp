package clickup

import (
	"context"
	"net/http"
	"net/url"
)

// ListListsInFolder returns the lists in a folder.
// GET /folder/{folder_id}/list
func (c *Client) ListListsInFolder(ctx context.Context, folderID string, archived bool, archivedSet bool) (any, error) {
	q := url.Values{}
	addBoolParam(q, "archived", archived, archivedSet)
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/folder/" + folderID + "/list", Query: q}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListFolderlessLists returns the folderless lists directly under a space.
// GET /space/{space_id}/list
func (c *Client) ListFolderlessLists(ctx context.Context, spaceID string, archived bool, archivedSet bool) (any, error) {
	q := url.Values{}
	addBoolParam(q, "archived", archived, archivedSet)
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/space/" + spaceID + "/list", Query: q}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetList returns a single list.
// GET /list/{list_id}
func (c *Client) GetList(ctx context.Context, listID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/list/" + listID}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateListInFolder creates a list inside a folder.
// POST /folder/{folder_id}/list
func (c *Client) CreateListInFolder(ctx context.Context, folderID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/folder/" + folderID + "/list", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateFolderlessList creates a list directly under a space.
// POST /space/{space_id}/list
func (c *Client) CreateFolderlessList(ctx context.Context, spaceID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/space/" + spaceID + "/list", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateList updates a list.
// PUT /list/{list_id}
func (c *Client) UpdateList(ctx context.Context, listID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPut, Path: "/list/" + listID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteList deletes a list.
// DELETE /list/{list_id}
func (c *Client) DeleteList(ctx context.Context, listID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/list/" + listID}, nil)
}

// AddTaskToList adds an existing task to an additional list.
// POST /list/{list_id}/task/{task_id}
func (c *Client) AddTaskToList(ctx context.Context, listID, taskID string) error {
	return c.do(ctx, requestParams{Method: http.MethodPost, Path: "/list/" + listID + "/task/" + taskID}, nil)
}

// RemoveTaskFromList removes a task from an additional list (does not delete it).
// DELETE /list/{list_id}/task/{task_id}
func (c *Client) RemoveTaskFromList(ctx context.Context, listID, taskID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/list/" + listID + "/task/" + taskID}, nil)
}

// CreateListFromTemplateInFolder creates a list in a folder from a list template.
// POST /folder/{folder_id}/list_template/{template_id}
func (c *Client) CreateListFromTemplateInFolder(ctx context.Context, folderID, templateID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/folder/" + folderID + "/list_template/" + templateID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateListFromTemplateInSpace creates a folderless list in a space from a list template.
// POST /space/{space_id}/list_template/{template_id}
func (c *Client) CreateListFromTemplateInSpace(ctx context.Context, spaceID, templateID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/space/" + spaceID + "/list_template/" + templateID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
