package clickup

import (
	"context"
	"net/http"
)

// ListTaskComments returns comments on a task.
// GET /task/{task_id}/comment
func (c *Client) ListTaskComments(ctx context.Context, taskID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/task/" + taskID + "/comment"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateTaskComment adds a comment to a task.
// POST /task/{task_id}/comment
func (c *Client) CreateTaskComment(ctx context.Context, taskID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/task/" + taskID + "/comment", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListListComments returns comments on a list.
// GET /list/{list_id}/comment
func (c *Client) ListListComments(ctx context.Context, listID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/list/" + listID + "/comment"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateListComment adds a comment to a list.
// POST /list/{list_id}/comment
func (c *Client) CreateListComment(ctx context.Context, listID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/list/" + listID + "/comment", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListViewComments returns comments on a view (used for chat views).
// GET /view/{view_id}/comment
func (c *Client) ListViewComments(ctx context.Context, viewID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/view/" + viewID + "/comment"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateViewComment adds a comment to a view.
// POST /view/{view_id}/comment
func (c *Client) CreateViewComment(ctx context.Context, viewID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/view/" + viewID + "/comment", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateComment updates a comment's text/assignee/resolved state.
// PUT /comment/{comment_id}
func (c *Client) UpdateComment(ctx context.Context, commentID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPut, Path: "/comment/" + commentID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteComment deletes a comment.
// DELETE /comment/{comment_id}
func (c *Client) DeleteComment(ctx context.Context, commentID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/comment/" + commentID}, nil)
}

// ListCommentReplies returns the threaded replies on a comment.
// GET /comment/{comment_id}/reply
func (c *Client) ListCommentReplies(ctx context.Context, commentID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/comment/" + commentID + "/reply"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateCommentReply adds a threaded reply to a comment.
// POST /comment/{comment_id}/reply
func (c *Client) CreateCommentReply(ctx context.Context, commentID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/comment/" + commentID + "/reply", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
