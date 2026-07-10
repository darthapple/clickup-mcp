package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

// RegisterReportTools registers cross-cutting time-tracking report tools
// that aren't a single ClickUp REST resource: a per-list report (every task
// in a list plus its tracked time) and a per-user timesheet (every entry a
// user logged, across all spaces/lists/tasks).
func RegisterReportTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_get_list_time_report",
			mcp.WithDescription("Return every task in a ClickUp list plus all time tracked to "+
				"each task within a date range, aggregated across every member of the list "+
				"(not just the calling token's own user) — internally resolves the list's "+
				"members and queries their entries together, so a task worked on by multiple "+
				"people reports everyone's combined time, not just one person's. A user with "+
				"time logged on this list but who is no longer a list member (e.g. removed "+
				"after logging time) is not covered by this aggregation and their entries "+
				"will be missing; use clickup_list_time_entries with an explicit assignee for "+
				"that case. Includes closed/completed tasks, but NOT "+
				"subtasks (ClickUp's list endpoint omits them unless subtasks are requested "+
				"separately) — any time entry logged against something other than a top-level "+
				"task in this list (a subtask, or a deleted task) is returned under top-level "+
				"\"unmatched_entries\" instead of being silently dropped; each unmatched entry "+
				"still carries its own task_id/task_name so the caller can tell what it belongs "+
				"to. Response shape: {list_id, start_date, end_date, tasks: [{task_id, "+
				"task_name, status, duration_ms, duration_formatted, entries: [{id, task_id, "+
				"task_name, start, end, duration_ms, duration_formatted, description}]}], "+
				"total_duration_ms, total_duration_formatted, unmatched_entries: [...same "+
				"entry shape...]}. Tasks with no tracked time in the period still appear, with "+
				"duration_ms 0 and an empty entries array. Each task's duration is reported as "+
				"both raw milliseconds and a zero-padded \"DD:HH:MM:SS\" string "+
				"(days:hours:minutes:seconds); entries within a task carry the same two "+
				"forms. All timestamps (start_date, end_date, start, end) are human-readable "+
				"UTC datetime strings (\"YYYY-MM-DD HH:MM:SS\") — this applies to start_date/"+
				"end_date here even though those same field names render as a bare date "+
				"elsewhere (e.g. a task's own due_date/start_date), since here they're the "+
				"query range boundary, not a calendar-date field. Known limitation: the "+
				"underlying time-entries endpoint has no documented pagination, so an "+
				"extremely high-volume list over a long range could theoretically be capped "+
				"by ClickUp server-side — not yet observed, but not provably ruled out."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("list_id", mcp.Required(), mcp.Description("List ID to report on")),
			mcp.WithString("start_date", mcp.Required(), mcp.Description("Range start, UTC \"YYYY-MM-DD HH:MM:SS\" or bare \"YYYY-MM-DD\" (midnight UTC)")),
			mcp.WithString("end_date", mcp.Required(), mcp.Description("Range end, UTC \"YYYY-MM-DD HH:MM:SS\" or bare \"YYYY-MM-DD\" (midnight UTC)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			listID, err := req.RequireString("list_id")
			if err != nil {
				return ErrorResult(err)
			}
			start, err := requireDateTimeArg(req, "start_date")
			if err != nil {
				return ErrorResult(err)
			}
			end, err := requireDateTimeArg(req, "end_date")
			if err != nil {
				return ErrorResult(err)
			}
			teamID := teamIDOrDefault(req, c)

			tasks, err := fetchAllTasksInList(ctx, c, listID)
			if err != nil {
				return ErrorResult(err)
			}
			grouped, err := fetchTimeEntriesByTask(ctx, c, teamID, listID, start, end)
			if err != nil {
				return ErrorResult(err)
			}

			out := []map[string]any{}
			var totalMs int64
			consumed := map[string]bool{}
			for _, t := range tasks {
				id, _ := t["id"].(string)
				name, _ := t["name"].(string)
				status, _ := t["status"].(map[string]any)
				statusName, _ := status["status"].(string)
				consumed[id] = true
				entries := grouped[id]
				if entries == nil {
					entries = []map[string]any{}
				}
				var taskMs int64
				for _, e := range entries {
					if d, ok := e["duration_ms"].(int64); ok {
						taskMs += d
					}
				}
				totalMs += taskMs
				out = append(out, map[string]any{
					"task_id":            id,
					"task_name":          name,
					"status":             statusName,
					"duration_ms":        taskMs,
					"duration_formatted": formatDuration(taskMs),
					"entries":            entries,
				})
			}
			// Entries whose task_id was never consumed above (e.g. logged
			// against a subtask, which ClickUp's list endpoint doesn't
			// return) would otherwise be silently lost since their grouped
			// bucket is never read — surface them instead.
			unmatched := []map[string]any{}
			for taskID, es := range grouped {
				if !consumed[taskID] {
					unmatched = append(unmatched, es...)
				}
			}
			return JSONResult(map[string]any{
				"list_id":                  listID,
				"start_date":               start,
				"end_date":                 end,
				"tasks":                    out,
				"total_duration_ms":        totalMs,
				"total_duration_formatted": formatDuration(totalMs),
				"unmatched_entries":        unmatched,
			}, "start_date")
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_get_user_time_report",
			mcp.WithDescription("Return a timesheet for one or more ClickUp users: every time "+
				"entry they logged in a date range, across all spaces/lists/tasks (or "+
				"restricted to one space via space_id), with each entry's space, folder, "+
				"list, task, and logging-user name resolved. Includes a per-task rollup and "+
				"a grand total. user_id accepts a single ID or a comma-separated list (e.g. "+
				"\"170440755,87915023\") to cover multiple people in one call — a task "+
				"worked on by several people this way is not silently attributed to only "+
				"one of them, since each entry still carries its own user_id/user_name. "+
				"IMPORTANT: when multiple users are requested, by_task and "+
				"total_duration_ms/total_duration_formatted aggregate across ALL requested "+
				"users combined, not broken down per person — group entries yourself by "+
				"user_id if you need a per-analyst total. Response shape: {user_id, "+
				"space_id, start_date, end_date, entries: [{space_id, space_name, "+
				"folder_id, folder_name, list_id, list_name, task_id, task_name, id, "+
				"user_id, user_name, start, end, duration_ms, duration_formatted, "+
				"description}], by_task: [{task_id, task_name, list_name, folder_name, "+
				"space_name, duration_ms, duration_formatted}], total_duration_ms, "+
				"total_duration_formatted}. If a list's name can't be resolved (e.g. "+
				"deleted since the entry was logged), its space_name/folder_name/list_name "+
				"come back as empty strings rather than failing the whole report — an "+
				"empty name means \"lookup failed\", not \"no folder\". Duration is "+
				"reported as both raw milliseconds and a zero-padded \"DD:HH:MM:SS\" "+
				"string (days:hours:minutes:seconds). All timestamps (start_date, "+
				"end_date, start, end) are human-readable UTC datetime strings "+
				"(\"YYYY-MM-DD HH:MM:SS\") — this applies to start_date/end_date here "+
				"even though those same field names render as a bare date elsewhere (e.g. "+
				"a task's own due_date/start_date), since here they're the query range "+
				"boundary, not a calendar-date field."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("user_id", mcp.Required(), mcp.Description("ClickUp user ID (assignee) to report on. Accepts multiple user IDs as a single comma-separated string (e.g. \"170440755,87915023\") to cover several people in one call.")),
			mcp.WithString("space_id", mcp.Description("Optional: restrict to time entries within this space only. Omit for a full cross-workspace timesheet.")),
			mcp.WithString("start_date", mcp.Required(), mcp.Description("Range start, UTC \"YYYY-MM-DD HH:MM:SS\" or bare \"YYYY-MM-DD\" (midnight UTC)")),
			mcp.WithString("end_date", mcp.Required(), mcp.Description("Range end, UTC \"YYYY-MM-DD HH:MM:SS\" or bare \"YYYY-MM-DD\" (midnight UTC)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			userID, err := req.RequireString("user_id")
			if err != nil {
				return ErrorResult(err)
			}
			spaceID := req.GetString("space_id", "")
			start, err := requireDateTimeArg(req, "start_date")
			if err != nil {
				return ErrorResult(err)
			}
			end, err := requireDateTimeArg(req, "end_date")
			if err != nil {
				return ErrorResult(err)
			}
			teamID := teamIDOrDefault(req, c)

			rawEntries, err := fetchUserTimeEntries(ctx, c, teamID, userID, spaceID, start, end)
			if err != nil {
				return ErrorResult(err)
			}

			resolveList := newListInfoResolver(ctx, c)
			entries := []map[string]any{}
			type taskAgg struct {
				taskName, listName, folderName, spaceName string
				durationMs                                int64
			}
			byTask := map[string]*taskAgg{}
			var totalMs int64

			for _, em := range rawEntries {
				taskObj, _ := em["task"].(map[string]any)
				taskID, _ := taskObj["id"].(string)
				taskName, _ := taskObj["name"].(string)
				loc, _ := em["task_location"].(map[string]any)
				listID, _ := loc["list_id"].(string)
				info := resolveList(listID)

				dur := parseMsField(em, "duration")
				id, _ := em["id"].(string)
				desc, _ := em["description"].(string)
				entryUserID, entryUserName := entryUser(em)

				entries = append(entries, map[string]any{
					"space_id":           info.spaceID,
					"space_name":         info.spaceName,
					"folder_id":          loc["folder_id"],
					"folder_name":        info.folderName,
					"list_id":            listID,
					"list_name":          info.listName,
					"task_id":            taskID,
					"task_name":          taskName,
					"id":                 id,
					"user_id":            entryUserID,
					"user_name":          entryUserName,
					"start":              parseMsField(em, "start"),
					"end":                parseMsField(em, "end"),
					"duration_ms":        dur,
					"duration_formatted": formatDuration(dur),
					"description":        desc,
				})

				totalMs += dur
				agg, ok := byTask[taskID]
				if !ok {
					agg = &taskAgg{taskName: taskName, listName: info.listName, folderName: info.folderName, spaceName: info.spaceName}
					byTask[taskID] = agg
				}
				agg.durationMs += dur
			}

			byTaskOut := []map[string]any{}
			for taskID, agg := range byTask {
				byTaskOut = append(byTaskOut, map[string]any{
					"task_id":            taskID,
					"task_name":          agg.taskName,
					"list_name":          agg.listName,
					"folder_name":        agg.folderName,
					"space_name":         agg.spaceName,
					"duration_ms":        agg.durationMs,
					"duration_formatted": formatDuration(agg.durationMs),
				})
			}

			return JSONResult(map[string]any{
				"user_id":                  userID,
				"space_id":                 spaceID,
				"start_date":               start,
				"end_date":                 end,
				"entries":                  entries,
				"by_task":                  byTaskOut,
				"total_duration_ms":        totalMs,
				"total_duration_formatted": formatDuration(totalMs),
			}, "start_date")
		},
	)
}

// fetchAllTasksInList paginates GetTasksInList until ClickUp reports
// last_page, since the client only fetches one page per call.
func fetchAllTasksInList(ctx context.Context, c *clickup.Client, listID string) ([]map[string]any, error) {
	var all []map[string]any
	includeClosed := true
	page := 0
	for {
		raw, err := c.GetTasksInList(ctx, listID, clickup.TaskQueryFilters{Page: &page, IncludeClosed: &includeClosed})
		if err != nil {
			return nil, fmt.Errorf("fetching tasks (page %d): %w", page, err)
		}
		m, _ := raw.(map[string]any)
		tasksRaw, _ := m["tasks"].([]any)
		if len(tasksRaw) == 0 {
			break
		}
		for _, t := range tasksRaw {
			if tm, ok := t.(map[string]any); ok {
				all = append(all, tm)
			}
		}
		if lastPage, _ := m["last_page"].(bool); lastPage {
			break
		}
		page++
		if page > 500 { // safety cap in case last_page is ever missing
			return nil, fmt.Errorf("aborted after 500 pages for list %s", listID)
		}
	}
	return all, nil
}

// fetchTimeEntriesByTask fetches all time entries for a list in one call and
// groups them by task ID. Each entry carries its own task_id/task_name so
// entries whose task didn't turn up in the task listing (e.g. a deleted
// task, or a subtask, since ClickUp's list endpoint omits subtasks) can
// still be surfaced and identified by the caller instead of silently
// dropped.
//
// The underlying time-entries endpoint defaults to only the calling token's
// own entries when no assignee is given, so omitting it here would silently
// under-report every other user's time (a full report reading back as "no
// time tracked" rather than erroring). This resolves the list's members
// first and passes them all as a comma-separated assignee filter, which
// ClickUp's API accepts, so the report aggregates everyone's time by default.
func fetchTimeEntriesByTask(ctx context.Context, c *clickup.Client, teamID, listID string, start, end int64) (map[string][]map[string]any, error) {
	assignees, err := fetchListMemberIDs(ctx, c, listID)
	if err != nil {
		return nil, fmt.Errorf("fetching members of list %s: %w", listID, err)
	}
	raw, err := c.ListTimeEntries(ctx, teamID, clickup.TimeEntryFilters{ListID: listID, Assignee: assignees, StartDate: &start, EndDate: &end})
	if err != nil {
		return nil, fmt.Errorf("fetching time entries for list %s: %w", listID, err)
	}
	m, _ := raw.(map[string]any)
	entriesRaw, _ := m["data"].([]any)
	grouped := map[string][]map[string]any{}
	for _, e := range entriesRaw {
		em, ok := e.(map[string]any)
		if !ok {
			continue
		}
		taskObj, _ := em["task"].(map[string]any)
		taskID, _ := taskObj["id"].(string)
		taskName, _ := taskObj["name"].(string)
		dur := parseMsField(em, "duration")
		id, _ := em["id"].(string)
		desc, _ := em["description"].(string)
		entry := map[string]any{
			"id":                 id,
			"task_id":            taskID,
			"task_name":          taskName,
			"start":              parseMsField(em, "start"),
			"end":                parseMsField(em, "end"),
			"duration_ms":        dur,
			"duration_formatted": formatDuration(dur),
			"description":        desc,
		}
		grouped[taskID] = append(grouped[taskID], entry)
	}
	return grouped, nil
}

// fetchListMemberIDs returns every member of a list as a single
// comma-separated string of user IDs, suitable for TimeEntryFilters.Assignee
// — the time-entries endpoint accepts a comma-separated list there. Member
// IDs come back from ClickUp as JSON numbers, not strings (unlike task/list
// IDs elsewhere in this API), so they're reformatted rather than type-asserted.
func fetchListMemberIDs(ctx context.Context, c *clickup.Client, listID string) (string, error) {
	raw, err := c.ListListMembers(ctx, listID)
	if err != nil {
		return "", err
	}
	m, _ := raw.(map[string]any)
	membersRaw, _ := m["members"].([]any)
	ids := make([]string, 0, len(membersRaw))
	for _, mem := range membersRaw {
		mm, ok := mem.(map[string]any)
		if !ok {
			continue
		}
		if id, ok := mm["id"].(float64); ok {
			ids = append(ids, strconv.FormatInt(int64(id), 10))
		}
	}
	return strings.Join(ids, ","), nil
}

// fetchUserTimeEntries fetches every time entry logged by userID (a single
// ID, or ClickUp's accepted comma-separated list of IDs for a multi-user
// report) in a date range, optionally restricted to one space via spaceID
// (empty string means no space restriction — the full cross-workspace
// timesheet).
func fetchUserTimeEntries(ctx context.Context, c *clickup.Client, teamID, userID, spaceID string, start, end int64) ([]map[string]any, error) {
	raw, err := c.ListTimeEntries(ctx, teamID, clickup.TimeEntryFilters{Assignee: userID, SpaceID: spaceID, StartDate: &start, EndDate: &end})
	if err != nil {
		return nil, fmt.Errorf("fetching time entries for user %s: %w", userID, err)
	}
	m, _ := raw.(map[string]any)
	entriesRaw, _ := m["data"].([]any)
	var out []map[string]any
	for _, e := range entriesRaw {
		if em, ok := e.(map[string]any); ok {
			out = append(out, em)
		}
	}
	return out, nil
}

// entryUser extracts the id/username of whoever logged a raw time entry,
// from its embedded "user" object. Like list-member IDs (fetchListMemberIDs
// above), ClickUp returns this id as a JSON number, not a string. Missing or
// malformed data degrades to an empty string rather than failing the whole
// report, consistent with this file's other display-enrichment lookups.
func entryUser(em map[string]any) (id, name string) {
	userObj, _ := em["user"].(map[string]any)
	if userObj == nil {
		return "", ""
	}
	if idNum, ok := userObj["id"].(float64); ok {
		id = strconv.FormatInt(int64(idNum), 10)
	}
	name, _ = userObj["username"].(string)
	return id, name
}

// listInfo is the display metadata resolved for a list ID: its own name plus
// its parent folder's and space's names.
type listInfo struct {
	listName, folderName, spaceID, spaceName string
}

// newListInfoResolver returns a memoized resolver: GetList embeds both
// folder.name and space.name inline, so one call per distinct list ID is
// enough (never one call per entry). A lookup failure degrades to empty
// names rather than aborting the report, since this is display enrichment,
// not primary requested data.
func newListInfoResolver(ctx context.Context, c *clickup.Client) func(listID string) listInfo {
	cache := map[string]listInfo{}
	return func(listID string) listInfo {
		if v, ok := cache[listID]; ok {
			return v
		}
		info := listInfo{}
		if raw, err := c.GetList(ctx, listID); err == nil {
			if m, ok := raw.(map[string]any); ok {
				info.listName, _ = m["name"].(string)
				if folder, ok := m["folder"].(map[string]any); ok {
					info.folderName, _ = folder["name"].(string)
				}
				if space, ok := m["space"].(map[string]any); ok {
					info.spaceID, _ = space["id"].(string)
					info.spaceName, _ = space["name"].(string)
				}
			}
		}
		cache[listID] = info
		return info
	}
}

// parseMsField reads a ClickUp millisecond-epoch field, which is encoded as
// a JSON string, not a number — parsed with ParseInt (never via float) to
// avoid precision loss on large timestamp values.
func parseMsField(m map[string]any, key string) int64 {
	s, _ := m[key].(string)
	if s == "" {
		return 0
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return n
}

// formatDuration renders a millisecond duration as zero-padded
// "DD:HH:MM:SS" (days, hours, minutes, seconds), using integer arithmetic
// throughout.
func formatDuration(ms int64) string {
	if ms < 0 {
		ms = 0
	}
	total := ms / 1000
	days, total := total/86400, total%86400
	hours, total := total/3600, total%3600
	minutes, seconds := total/60, total%60
	return fmt.Sprintf("%02d:%02d:%02d:%02d", days, hours, minutes, seconds)
}
