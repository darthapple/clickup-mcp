package clickup

import (
	"context"
	"net/http"
)

// CreateAuditLogReport requests an audit log export for a workspace.
// Enterprise plans only; non-Enterprise workspaces get a 4xx APIError, which
// is expected, not a bug.
// POST /v3/workspaces/{workspace_id}/auditlogs
func (c *Client) CreateAuditLogReport(ctx context.Context, workspaceID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, APIVersion: apiV3, Path: "/workspaces/" + workspaceID + "/auditlogs", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
