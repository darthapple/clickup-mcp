package clickup

import (
	"context"
	"net/http"
)

// ListGoals returns the goals in a workspace.
// GET /team/{team_id}/goal
func (c *Client) ListGoals(ctx context.Context, teamID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/team/" + teamID + "/goal"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateGoal creates a goal in a workspace.
// POST /team/{team_id}/goal
func (c *Client) CreateGoal(ctx context.Context, teamID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/team/" + teamID + "/goal", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetGoal returns a single goal.
// GET /goal/{goal_id}
func (c *Client) GetGoal(ctx context.Context, goalID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/goal/" + goalID}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateGoal updates a goal.
// PUT /goal/{goal_id}
func (c *Client) UpdateGoal(ctx context.Context, goalID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPut, Path: "/goal/" + goalID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteGoal deletes a goal.
// DELETE /goal/{goal_id}
func (c *Client) DeleteGoal(ctx context.Context, goalID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/goal/" + goalID}, nil)
}

// CreateKeyResult creates a key result (target) on a goal.
// POST /goal/{goal_id}/key_result
func (c *Client) CreateKeyResult(ctx context.Context, goalID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/goal/" + goalID + "/key_result", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateKeyResult updates a key result's progress/value.
// PUT /key_result/{key_result_id}
func (c *Client) UpdateKeyResult(ctx context.Context, keyResultID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPut, Path: "/key_result/" + keyResultID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteKeyResult deletes a key result.
// DELETE /key_result/{key_result_id}
func (c *Client) DeleteKeyResult(ctx context.Context, keyResultID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/key_result/" + keyResultID}, nil)
}
