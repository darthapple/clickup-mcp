package clickup

import (
	"context"
	"net/http"
	"net/url"
)

// TimeEntryFilters are the optional filters for ListTimeEntries.
type TimeEntryFilters struct {
	StartDate  *int64
	EndDate    *int64
	Assignee   string
	SpaceID    string
	FolderID   string
	ListID     string
	TaskID     string
	IncludeAll *bool
}

func (f TimeEntryFilters) apply(q url.Values) {
	if f.StartDate != nil {
		addIntParam(q, "start_date", int(*f.StartDate), true)
	}
	if f.EndDate != nil {
		addIntParam(q, "end_date", int(*f.EndDate), true)
	}
	addParam(q, "assignee", f.Assignee)
	addParam(q, "space_id", f.SpaceID)
	addParam(q, "folder_id", f.FolderID)
	addParam(q, "list_id", f.ListID)
	addParam(q, "task_id", f.TaskID)
	if f.IncludeAll != nil {
		addBoolParam(q, "include_task_tags", *f.IncludeAll, true)
	}
}

// ListTimeEntries returns time entries in a workspace, with optional filters.
// GET /team/{team_id}/time_entries
func (c *Client) ListTimeEntries(ctx context.Context, teamID string, filters TimeEntryFilters) (any, error) {
	q := url.Values{}
	filters.apply(q)
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/team/" + teamID + "/time_entries", Query: q}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateTimeEntry creates a manual time entry.
// POST /team/{team_id}/time_entries
func (c *Client) CreateTimeEntry(ctx context.Context, teamID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/team/" + teamID + "/time_entries", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTimeEntry returns a single time entry.
// GET /team/{team_id}/time_entries/{timer_id}
func (c *Client) GetTimeEntry(ctx context.Context, teamID, timerID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/team/" + teamID + "/time_entries/" + timerID}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateTimeEntry updates a time entry.
// PUT /team/{team_id}/time_entries/{timer_id}
func (c *Client) UpdateTimeEntry(ctx context.Context, teamID, timerID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPut, Path: "/team/" + teamID + "/time_entries/" + timerID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteTimeEntry deletes a time entry.
// DELETE /team/{team_id}/time_entries/{timer_id}
func (c *Client) DeleteTimeEntry(ctx context.Context, teamID, timerID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/team/" + teamID + "/time_entries/" + timerID}, nil)
}

// GetTimeEntryHistory returns the edit history of a time entry.
// GET /team/{team_id}/time_entries/{timer_id}/history
func (c *Client) GetTimeEntryHistory(ctx context.Context, teamID, timerID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/team/" + teamID + "/time_entries/" + timerID + "/history"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCurrentTimeEntry returns the currently running timer, if any.
// GET /team/{team_id}/time_entries/current
func (c *Client) GetCurrentTimeEntry(ctx context.Context, teamID, assignee string) (any, error) {
	q := url.Values{}
	addParam(q, "assignee", assignee)
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/team/" + teamID + "/time_entries/current", Query: q}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// StartTimeEntry starts a timer on a task.
// POST /team/{team_id}/time_entries/start/{task_id}
func (c *Client) StartTimeEntry(ctx context.Context, teamID, taskID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/team/" + teamID + "/time_entries/start/" + taskID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// StopTimeEntry stops the currently running timer.
// POST /team/{team_id}/time_entries/stop
func (c *Client) StopTimeEntry(ctx context.Context, teamID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/team/" + teamID + "/time_entries/stop"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListTimeEntryTags returns all time-tracking tags in a workspace.
// GET /team/{team_id}/time_entries/tags
func (c *Client) ListTimeEntryTags(ctx context.Context, teamID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/team/" + teamID + "/time_entries/tags"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AddTimeEntryTags adds tags to time entries. body: {"time_entry_ids": [...], "tags": [...]}
// POST /team/{team_id}/time_entries/tags
func (c *Client) AddTimeEntryTags(ctx context.Context, teamID string, body map[string]any) error {
	return c.do(ctx, requestParams{Method: http.MethodPost, Path: "/team/" + teamID + "/time_entries/tags", Body: body}, nil)
}

// RenameTimeEntryTag renames a time-tracking tag.
// PUT /team/{team_id}/time_entries/tags
func (c *Client) RenameTimeEntryTag(ctx context.Context, teamID string, body map[string]any) error {
	return c.do(ctx, requestParams{Method: http.MethodPut, Path: "/team/" + teamID + "/time_entries/tags", Body: body}, nil)
}

// RemoveTimeEntryTags removes tags from time entries. body: {"time_entry_ids": [...], "tags": [...]}
// DELETE /team/{team_id}/time_entries/tags
func (c *Client) RemoveTimeEntryTags(ctx context.Context, teamID string, body map[string]any) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/team/" + teamID + "/time_entries/tags", Body: body}, nil)
}
