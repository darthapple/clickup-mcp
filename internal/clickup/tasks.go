package clickup

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// TaskQueryFilters are the optional list/search filters shared by
// GetTasksInList and SearchTasks.
type TaskQueryFilters struct {
	Archived           *bool
	IncludeClosed      *bool
	Subtasks           *bool
	Page               *int
	OrderBy            string
	Reverse            *bool
	Statuses           []string
	Assignees          []string
	Tags               []string
	DueDateGT          *int64
	DueDateLT          *int64
	DateCreatedGT      *int64
	DateCreatedLT      *int64
	DateUpdatedGT      *int64
	DateUpdatedLT      *int64
	CustomFieldsFilter string // raw JSON array, ClickUp's `custom_fields` query param
}

func (f TaskQueryFilters) apply(q url.Values) {
	if f.Archived != nil {
		addBoolParam(q, "archived", *f.Archived, true)
	}
	if f.IncludeClosed != nil {
		addBoolParam(q, "include_closed", *f.IncludeClosed, true)
	}
	if f.Subtasks != nil {
		addBoolParam(q, "subtasks", *f.Subtasks, true)
	}
	if f.Page != nil {
		addIntParam(q, "page", *f.Page, true)
	}
	addParam(q, "order_by", f.OrderBy)
	if f.Reverse != nil {
		addBoolParam(q, "reverse", *f.Reverse, true)
	}
	addArrayParam(q, "statuses[]", f.Statuses)
	addArrayParam(q, "assignees[]", f.Assignees)
	addArrayParam(q, "tags[]", f.Tags)
	if f.DueDateGT != nil {
		q.Set("due_date_gt", strconv.FormatInt(*f.DueDateGT, 10))
	}
	if f.DueDateLT != nil {
		q.Set("due_date_lt", strconv.FormatInt(*f.DueDateLT, 10))
	}
	if f.DateCreatedGT != nil {
		q.Set("date_created_gt", strconv.FormatInt(*f.DateCreatedGT, 10))
	}
	if f.DateCreatedLT != nil {
		q.Set("date_created_lt", strconv.FormatInt(*f.DateCreatedLT, 10))
	}
	if f.DateUpdatedGT != nil {
		q.Set("date_updated_gt", strconv.FormatInt(*f.DateUpdatedGT, 10))
	}
	if f.DateUpdatedLT != nil {
		q.Set("date_updated_lt", strconv.FormatInt(*f.DateUpdatedLT, 10))
	}
	addParam(q, "custom_fields", f.CustomFieldsFilter)
}

// GetTask returns a single task.
// GET /task/{task_id}
func (c *Client) GetTask(ctx context.Context, taskID string, customTaskIDs bool, teamID string) (any, error) {
	q := url.Values{}
	addBoolParam(q, "custom_task_ids", customTaskIDs, customTaskIDs)
	addParam(q, "team_id", teamID)
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/task/" + taskID, Query: q}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateTask creates a task in a list.
// POST /list/{list_id}/task
func (c *Client) CreateTask(ctx context.Context, listID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/list/" + listID + "/task", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateTask updates a task.
// PUT /task/{task_id}
func (c *Client) UpdateTask(ctx context.Context, taskID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPut, Path: "/task/" + taskID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteTask deletes a task.
// DELETE /task/{task_id}
func (c *Client) DeleteTask(ctx context.Context, taskID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/task/" + taskID}, nil)
}

// GetTasksInList returns the tasks in a list, with optional filters.
// GET /list/{list_id}/task
func (c *Client) GetTasksInList(ctx context.Context, listID string, filters TaskQueryFilters) (any, error) {
	q := url.Values{}
	filters.apply(q)
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/list/" + listID + "/task", Query: q}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SearchTasksOptions adds workspace-search-only filters on top of TaskQueryFilters.
type SearchTasksOptions struct {
	TaskQueryFilters
	SpaceIDs  []string
	ListIDs   []string
	FolderIDs []string
}

// SearchTasks searches tasks across an entire workspace.
// GET /team/{team_id}/task
func (c *Client) SearchTasks(ctx context.Context, teamID string, opts SearchTasksOptions) (any, error) {
	q := url.Values{}
	opts.TaskQueryFilters.apply(q)
	addArrayParam(q, "space_ids[]", opts.SpaceIDs)
	addArrayParam(q, "list_ids[]", opts.ListIDs)
	addArrayParam(q, "project_ids[]", opts.FolderIDs)
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/team/" + teamID + "/task", Query: q}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTaskTimeInStatus returns how long a task has spent in each status.
// GET /task/{task_id}/time_in_status
func (c *Client) GetTaskTimeInStatus(ctx context.Context, taskID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/task/" + taskID + "/time_in_status"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetBulkTimeInStatus returns time-in-status for multiple tasks at once.
// GET /task/bulk_time_in_status
func (c *Client) GetBulkTimeInStatus(ctx context.Context, taskIDs []string) (any, error) {
	q := url.Values{}
	addArrayParam(q, "task_ids[]", taskIDs)
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/task/bulk_time_in_status", Query: q}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AddTaskLink links two tasks together.
// POST /task/{task_id}/link/{links_to}
func (c *Client) AddTaskLink(ctx context.Context, taskID, linksTo string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/task/" + taskID + "/link/" + linksTo}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// RemoveTaskLink removes a link between two tasks.
// DELETE /task/{task_id}/link/{links_to}
func (c *Client) RemoveTaskLink(ctx context.Context, taskID, linksTo string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/task/" + taskID + "/link/" + linksTo}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
