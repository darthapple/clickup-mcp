package clickup

import (
	"context"
	"net/http"
)

// CreateChecklist creates a checklist on a task. body: {"name": "..."}
// POST /task/{task_id}/checklist
func (c *Client) CreateChecklist(ctx context.Context, taskID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/task/" + taskID + "/checklist", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateChecklist renames a checklist or changes its order.
// PUT /checklist/{checklist_id}
func (c *Client) UpdateChecklist(ctx context.Context, checklistID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPut, Path: "/checklist/" + checklistID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteChecklist deletes a checklist.
// DELETE /checklist/{checklist_id}
func (c *Client) DeleteChecklist(ctx context.Context, checklistID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/checklist/" + checklistID}, nil)
}

// CreateChecklistItem adds an item to a checklist.
// POST /checklist/{checklist_id}/checklist_item
func (c *Client) CreateChecklistItem(ctx context.Context, checklistID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/checklist/" + checklistID + "/checklist_item", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateChecklistItem updates a checklist item (name/assignee/resolved/parent).
// PUT /checklist/{checklist_id}/checklist_item/{checklist_item_id}
func (c *Client) UpdateChecklistItem(ctx context.Context, checklistID, itemID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPut, Path: "/checklist/" + checklistID + "/checklist_item/" + itemID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteChecklistItem deletes a checklist item.
// DELETE /checklist/{checklist_id}/checklist_item/{checklist_item_id}
func (c *Client) DeleteChecklistItem(ctx context.Context, checklistID, itemID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/checklist/" + checklistID + "/checklist_item/" + itemID}, nil)
}
