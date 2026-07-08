package clickup

import (
	"context"
	"net/http"
)

// ListTaskMembers returns the users who can see a task.
// GET /task/{task_id}/member
func (c *Client) ListTaskMembers(ctx context.Context, taskID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/task/" + taskID + "/member"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListListMembers returns the users who can see a list.
// GET /list/{list_id}/member
func (c *Client) ListListMembers(ctx context.Context, listID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/list/" + listID + "/member"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
