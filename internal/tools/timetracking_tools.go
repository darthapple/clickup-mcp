package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func RegisterTimeTrackingTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_list_time_entries",
			mcp.WithDescription("List time entries in a ClickUp workspace, with optional filters. "+
				"IMPORTANT: if start_date/end_date are omitted, ClickUp defaults to only "+
				"the last 30 days, AND only entries created by the authenticated token's "+
				"own user — other assignees' and older entries are silently dropped with "+
				"no error. For a complete or historical report, always pass explicit "+
				"start_date/end_date (and assignee, to see another user's time), or "+
				"scope with task_id/list_id/folder_id/space_id instead. For a "+
				"ready-made per-list report (every task plus rolled-up/totaled time) "+
				"or a full per-user timesheet (every entry resolved to space/folder/"+
				"list/task names, with per-task totals), use "+
				"clickup_get_list_time_report or clickup_get_user_time_report "+
				"instead of assembling one by hand from this tool's raw entries. "+
				"Known limitation: this endpoint has no documented pagination, so an "+
				"extremely high-volume query could theoretically be capped by "+
				"ClickUp server-side — not yet observed, but not provably ruled out. "+
				"All returned timestamps (start, end, at) are human-readable UTC "+
				"datetime strings (\"YYYY-MM-DD HH:MM:SS\") — convert to the user's local "+
				"timezone before computing a calendar date, since a UTC-vs-local mismatch "+
				"can shift entries near midnight onto the wrong day. duration is still "+
				"raw milliseconds (an elapsed length, not a point in time)."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("start_date", mcp.Description("Filter: range start, UTC \"YYYY-MM-DD HH:MM:SS\" or bare \"YYYY-MM-DD\" (midnight UTC). Omitting this (along with end_date) restricts results to the last 30 days.")),
			mcp.WithString("end_date", mcp.Description("Filter: range end, UTC \"YYYY-MM-DD HH:MM:SS\" or bare \"YYYY-MM-DD\" (midnight UTC)")),
			mcp.WithString("assignee", mcp.Description("Filter by assignee user ID. Accepts multiple user IDs as a single comma-separated "+
				"string (e.g. \"170440755,87915023,118082738\") to fetch several users' entries in one call — "+
				"this is the only way to get a complete multi-user report, since looping per-user is easy to "+
				"forget and this param is not an array. Omitting this restricts results to the authenticated "+
				"token's own user.")),
			mcp.WithString("space_id", mcp.Description("Filter by space ID")),
			mcp.WithString("folder_id", mcp.Description("Filter by folder ID")),
			mcp.WithString("list_id", mcp.Description("Filter by list ID")),
			mcp.WithString("task_id", mcp.Description("Filter by task ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			filters := clickup.TimeEntryFilters{
				Assignee: req.GetString("assignee", ""),
				SpaceID:  req.GetString("space_id", ""),
				FolderID: req.GetString("folder_id", ""),
				ListID:   req.GetString("list_id", ""),
				TaskID:   req.GetString("task_id", ""),
			}
			if hasArg(req, "start_date") {
				v, err := requireDateTimeArg(req, "start_date")
				if err != nil {
					return ErrorResult(err)
				}
				filters.StartDate = &v
			}
			if hasArg(req, "end_date") {
				v, err := requireDateTimeArg(req, "end_date")
				if err != nil {
					return ErrorResult(err)
				}
				filters.EndDate = &v
			}
			out, err := c.ListTimeEntries(ctx, teamIDOrDefault(req, c), filters)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_time_entry",
			mcp.WithDescription("Create a manual time entry in a ClickUp workspace."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("task_id", mcp.Description("Task ID to log time against")),
			mcp.WithString("start", mcp.Required(), mcp.Description("Start time, UTC \"YYYY-MM-DD HH:MM:SS\" or bare \"YYYY-MM-DD\" (midnight UTC)")),
			mcp.WithNumber("duration", mcp.Required(), mcp.Description("Duration in milliseconds (not seconds)")),
			mcp.WithString("description", mcp.Description("Entry description")),
			mcp.WithBoolean("billable", mcp.Description("Mark as billable")),
			mcp.WithString("assignee", mcp.Description("User ID to log the time for")),
			mcp.WithArray("tags", mcp.WithStringItems(), mcp.Description("Tag names")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			start, err := requireDateTimeArg(req, "start")
			if err != nil {
				return ErrorResult(err)
			}
			duration, err := req.RequireFloat("duration")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{"start": start, "duration": int64(duration)}
			if hasArg(req, "task_id") {
				body["tid"] = req.GetString("task_id", "")
			}
			setString(body, req, "description")
			setBool(body, req, "billable")
			setString(body, req, "assignee")
			setStringSlice(body, req, "tags")
			out, err := c.CreateTimeEntry(ctx, teamIDOrDefault(req, c), body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_get_time_entry",
			mcp.WithDescription("Get a single ClickUp time entry."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("timer_id", mcp.Required(), mcp.Description("Time entry ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			timerID, err := req.RequireString("timer_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.GetTimeEntry(ctx, teamIDOrDefault(req, c), timerID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_time_entry",
			mcp.WithDescription("Update a ClickUp time entry."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("timer_id", mcp.Required(), mcp.Description("Time entry ID")),
			mcp.WithString("start", mcp.Description("Start time, UTC \"YYYY-MM-DD HH:MM:SS\" or bare \"YYYY-MM-DD\" (midnight UTC)")),
			mcp.WithNumber("duration", mcp.Description("Duration in milliseconds (not seconds)")),
			mcp.WithString("description", mcp.Description("Entry description")),
			mcp.WithBoolean("billable", mcp.Description("Mark as billable")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			timerID, err := req.RequireString("timer_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			if err := setDateTime(body, req, "start"); err != nil {
				return ErrorResult(err)
			}
			setFloat(body, req, "duration")
			setString(body, req, "description")
			setBool(body, req, "billable")
			out, err := c.UpdateTimeEntry(ctx, teamIDOrDefault(req, c), timerID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_delete_time_entry",
			mcp.WithDescription("Delete a ClickUp time entry."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("timer_id", mcp.Required(), mcp.Description("Time entry ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			timerID, err := req.RequireString("timer_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.DeleteTimeEntry(ctx, teamIDOrDefault(req, c), timerID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"deleted": true, "timer_id": timerID})
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_get_time_entry_history",
			mcp.WithDescription("Get the edit history of a ClickUp time entry."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("timer_id", mcp.Required(), mcp.Description("Time entry ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			timerID, err := req.RequireString("timer_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.GetTimeEntryHistory(ctx, teamIDOrDefault(req, c), timerID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_get_current_time_entry",
			mcp.WithDescription("Get the currently running ClickUp timer, if any."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("assignee", mcp.Description("User ID to check; defaults to the token owner")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			out, err := c.GetCurrentTimeEntry(ctx, teamIDOrDefault(req, c), req.GetString("assignee", ""))
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_start_time_entry",
			mcp.WithDescription("Start a timer on a ClickUp task."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
			mcp.WithString("description", mcp.Description("Entry description")),
			mcp.WithArray("tags", mcp.WithStringItems(), mcp.Description("Tag names")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			setString(body, req, "description")
			setStringSlice(body, req, "tags")
			out, err := c.StartTimeEntry(ctx, teamIDOrDefault(req, c), taskID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_stop_time_entry",
			mcp.WithDescription("Stop the currently running ClickUp timer."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			out, err := c.StopTimeEntry(ctx, teamIDOrDefault(req, c))
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_time_entry_tags",
			mcp.WithDescription("List all time-tracking tags in a ClickUp workspace."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			out, err := c.ListTimeEntryTags(ctx, teamIDOrDefault(req, c))
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_add_time_entry_tags",
			mcp.WithDescription("Add tags to ClickUp time entries."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithArray("time_entry_ids", mcp.Required(), mcp.WithStringItems(), mcp.Description("Time entry IDs")),
			mcp.WithArray("tags", mcp.Required(), mcp.WithStringItems(), mcp.Description("Tag names")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			ids, err := req.RequireStringSlice("time_entry_ids")
			if err != nil {
				return ErrorResult(err)
			}
			tags, err := req.RequireStringSlice("tags")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{"time_entry_ids": ids, "tags": tagObjects(tags)}
			if err := c.AddTimeEntryTags(ctx, teamIDOrDefault(req, c), body); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"added": true})
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_rename_time_entry_tag",
			mcp.WithDescription("Rename a ClickUp time-tracking tag."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Current tag name")),
			mcp.WithString("new_name", mcp.Required(), mcp.Description("New tag name")),
			mcp.WithString("tag_bg", mcp.Description("Background color, e.g. #ffffff")),
			mcp.WithString("tag_fg", mcp.Description("Foreground color, e.g. #000000")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name, err := req.RequireString("name")
			if err != nil {
				return ErrorResult(err)
			}
			newName, err := req.RequireString("new_name")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{
				"name":     name,
				"new_name": newName,
			}
			setString(body, req, "tag_bg")
			setString(body, req, "tag_fg")
			if err := c.RenameTimeEntryTag(ctx, teamIDOrDefault(req, c), body); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"renamed": true})
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_remove_time_entry_tags",
			mcp.WithDescription("Remove tags from ClickUp time entries."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithArray("time_entry_ids", mcp.Required(), mcp.WithStringItems(), mcp.Description("Time entry IDs")),
			mcp.WithArray("tags", mcp.Required(), mcp.WithStringItems(), mcp.Description("Tag names")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			ids, err := req.RequireStringSlice("time_entry_ids")
			if err != nil {
				return ErrorResult(err)
			}
			tags, err := req.RequireStringSlice("tags")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{"time_entry_ids": ids, "tags": tagObjects(tags)}
			if err := c.RemoveTimeEntryTags(ctx, teamIDOrDefault(req, c), body); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"removed": true})
		},
	)
}

// tagObjects converts plain tag names into ClickUp's {"name": ...} tag
// object shape used by the time-entry tag endpoints.
func tagObjects(names []string) []map[string]any {
	out := make([]map[string]any, len(names))
	for i, n := range names {
		out[i] = map[string]any{"name": n}
	}
	return out
}
