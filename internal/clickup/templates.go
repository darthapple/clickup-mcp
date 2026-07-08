package clickup

import (
	"context"
	"net/http"
)

// CreateTaskFromTemplate creates a task in a list from a task template.
// POST /list/{list_id}/taskTemplate/{template_id}
func (c *Client) CreateTaskFromTemplate(ctx context.Context, listID, templateID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/list/" + listID + "/taskTemplate/" + templateID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
