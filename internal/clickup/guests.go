package clickup

import (
	"context"
	"net/http"
)

// InviteGuest invites a guest to a workspace.
// POST /team/{team_id}/guest
func (c *Client) InviteGuest(ctx context.Context, teamID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/team/" + teamID + "/guest", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetGuest returns a single guest.
// GET /team/{team_id}/guest/{guest_id}
func (c *Client) GetGuest(ctx context.Context, teamID, guestID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/team/" + teamID + "/guest/" + guestID}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateGuest updates a guest's permissions.
// PUT /team/{team_id}/guest/{guest_id}
func (c *Client) UpdateGuest(ctx context.Context, teamID, guestID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPut, Path: "/team/" + teamID + "/guest/" + guestID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// RemoveGuestFromWorkspace removes a guest from a workspace entirely.
// DELETE /team/{team_id}/guest/{guest_id}
func (c *Client) RemoveGuestFromWorkspace(ctx context.Context, teamID, guestID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/team/" + teamID + "/guest/" + guestID}, nil)
}

// AddGuestToSpace shares a space with a guest.
// POST /space/{space_id}/guest/{guest_id}
func (c *Client) AddGuestToSpace(ctx context.Context, spaceID, guestID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/space/" + spaceID + "/guest/" + guestID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// RemoveGuestFromSpace revokes a guest's access to a space.
// DELETE /space/{space_id}/guest/{guest_id}
func (c *Client) RemoveGuestFromSpace(ctx context.Context, spaceID, guestID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/space/" + spaceID + "/guest/" + guestID}, nil)
}

// AddGuestToFolder shares a folder with a guest.
// POST /folder/{folder_id}/guest/{guest_id}
func (c *Client) AddGuestToFolder(ctx context.Context, folderID, guestID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/folder/" + folderID + "/guest/" + guestID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// RemoveGuestFromFolder revokes a guest's access to a folder.
// DELETE /folder/{folder_id}/guest/{guest_id}
func (c *Client) RemoveGuestFromFolder(ctx context.Context, folderID, guestID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/folder/" + folderID + "/guest/" + guestID}, nil)
}

// AddGuestToList shares a list with a guest.
// POST /list/{list_id}/guest/{guest_id}
func (c *Client) AddGuestToList(ctx context.Context, listID, guestID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/list/" + listID + "/guest/" + guestID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// RemoveGuestFromList revokes a guest's access to a list.
// DELETE /list/{list_id}/guest/{guest_id}
func (c *Client) RemoveGuestFromList(ctx context.Context, listID, guestID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/list/" + listID + "/guest/" + guestID}, nil)
}

// AddGuestToTask shares a task with a guest.
// POST /task/{task_id}/guest/{guest_id}
func (c *Client) AddGuestToTask(ctx context.Context, taskID, guestID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/task/" + taskID + "/guest/" + guestID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// RemoveGuestFromTask revokes a guest's access to a task.
// DELETE /task/{task_id}/guest/{guest_id}
func (c *Client) RemoveGuestFromTask(ctx context.Context, taskID, guestID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/task/" + taskID + "/guest/" + guestID}, nil)
}
