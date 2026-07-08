package clickup

import (
	"context"
	"net/http"
)

// ListWebhooks returns the webhooks registered in a workspace.
// GET /team/{team_id}/webhook
func (c *Client) ListWebhooks(ctx context.Context, teamID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/team/" + teamID + "/webhook"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateWebhook registers a webhook on a workspace.
// POST /team/{team_id}/webhook
func (c *Client) CreateWebhook(ctx context.Context, teamID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, Path: "/team/" + teamID + "/webhook", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateWebhook updates a webhook's events/endpoint/status.
// PUT /webhook/{webhook_id}
func (c *Client) UpdateWebhook(ctx context.Context, webhookID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPut, Path: "/webhook/" + webhookID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteWebhook removes a webhook.
// DELETE /webhook/{webhook_id}
func (c *Client) DeleteWebhook(ctx context.Context, webhookID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, Path: "/webhook/" + webhookID}, nil)
}
