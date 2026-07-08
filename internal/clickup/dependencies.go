package clickup

import (
	"context"
	"net/http"
	"net/url"
)

// AddTaskDependency creates a dependency relationship between two tasks.
// body: {"depends_on": "..."} or {"dependency_of": "..."}
// POST /task/{task_id}/dependency
func (c *Client) AddTaskDependency(ctx context.Context, taskID string, body map[string]any) error {
	return c.do(ctx, requestParams{Method: http.MethodPost, Path: "/task/" + taskID + "/dependency", Body: body}, nil)
}

// RemoveTaskDependency removes a dependency relationship between two tasks.
// DELETE /task/{task_id}/dependency
func (c *Client) RemoveTaskDependency(ctx context.Context, taskID, dependsOn, dependencyOf string) error {
	q := url.Values{}
	addParam(q, "depends_on", dependsOn)
	addParam(q, "dependency_of", dependencyOf)
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/task/" + taskID + "/dependency", Query: q}, nil)
}
