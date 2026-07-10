package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

// webhookEventsDescription lists every valid ClickUp webhook event name, so
// an agent isn't left guessing at names beyond the 2 examples ClickUp's own
// docs typically lead with.
const webhookEventsDescription = `Event names: taskCreated, taskUpdated, taskDeleted, taskPriorityUpdated, taskStatusUpdated, taskAssigneeUpdated, taskDueDateUpdated, taskTagUpdated, taskMoved, taskCommentPosted, taskCommentUpdated, taskTimeEstimateUpdated, taskTimeTrackedUpdated, listCreated, listUpdated, listDeleted, folderCreated, folderUpdated, folderDeleted, spaceCreated, spaceUpdated, spaceDeleted, goalCreated, goalUpdated, goalDeleted, keyResultCreated, keyResultUpdated, keyResultDeleted — or ["*"] for all.`

func RegisterWebhookTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_list_webhooks",
			mcp.WithDescription("List the webhooks registered on a ClickUp workspace."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			out, err := c.ListWebhooks(ctx, teamIDOrDefault(req, c))
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_webhook",
			mcp.WithDescription("Register a webhook on a ClickUp workspace. Only one of "+
				"space_id/folder_id/list_id/task_id takes effect — if multiple are "+
				"given, the most specific (task > list > folder > space) silently "+
				"wins with no error. If none are given, the webhook is workspace-wide."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("endpoint", mcp.Required(), mcp.Description("HTTPS URL to receive webhook events")),
			mcp.WithArray("events", mcp.Required(), mcp.WithStringItems(), mcp.Description(webhookEventsDescription)),
			mcp.WithString("space_id", mcp.Description("Restrict to a space; ignored if folder_id/list_id/task_id is also given")),
			mcp.WithString("folder_id", mcp.Description("Restrict to a folder; ignored if list_id/task_id is also given")),
			mcp.WithString("list_id", mcp.Description("Restrict to a list; ignored if task_id is also given")),
			mcp.WithString("task_id", mcp.Description("Restrict to a task; takes precedence over space_id/folder_id/list_id if multiple are given")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			endpoint, err := req.RequireString("endpoint")
			if err != nil {
				return ErrorResult(err)
			}
			events, err := req.RequireStringSlice("events")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{"endpoint": endpoint, "events": events}
			setString(body, req, "space_id")
			setString(body, req, "folder_id")
			setString(body, req, "list_id")
			setString(body, req, "task_id")
			out, err := c.CreateWebhook(ctx, teamIDOrDefault(req, c), body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_webhook",
			mcp.WithDescription("Update a ClickUp webhook's endpoint, events, or status."),
			mcp.WithString("webhook_id", mcp.Required(), mcp.Description("Webhook ID")),
			mcp.WithString("endpoint", mcp.Description("HTTPS URL to receive webhook events")),
			mcp.WithArray("events", mcp.WithStringItems(), mcp.Description("Event names to subscribe to (see clickup_create_webhook's events parameter for the full list of valid names)")),
			mcp.WithString("status", mcp.Description("active or suspended")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			webhookID, err := req.RequireString("webhook_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			setString(body, req, "endpoint")
			setStringSlice(body, req, "events")
			setString(body, req, "status")
			out, err := c.UpdateWebhook(ctx, webhookID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_delete_webhook",
			mcp.WithDescription("Delete a ClickUp webhook."),
			mcp.WithString("webhook_id", mcp.Required(), mcp.Description("Webhook ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			webhookID, err := req.RequireString("webhook_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.DeleteWebhook(ctx, webhookID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"deleted": true, "webhook_id": webhookID})
		},
	)
}
