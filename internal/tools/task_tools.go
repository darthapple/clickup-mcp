package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func taskQueryFilters(req mcp.CallToolRequest) clickup.TaskQueryFilters {
	f := clickup.TaskQueryFilters{
		OrderBy:            req.GetString("order_by", ""),
		Statuses:           req.GetStringSlice("statuses", nil),
		Assignees:          req.GetStringSlice("assignees", nil),
		Tags:               req.GetStringSlice("tags", nil),
		CustomFieldsFilter: req.GetString("custom_fields_json", ""),
	}
	if hasArg(req, "archived") {
		v := req.GetBool("archived", false)
		f.Archived = &v
	}
	if hasArg(req, "include_closed") {
		v := req.GetBool("include_closed", false)
		f.IncludeClosed = &v
	}
	if hasArg(req, "subtasks") {
		v := req.GetBool("subtasks", false)
		f.Subtasks = &v
	}
	if hasArg(req, "reverse") {
		v := req.GetBool("reverse", false)
		f.Reverse = &v
	}
	if hasArg(req, "page") {
		v := req.GetInt("page", 0)
		f.Page = &v
	}
	return f
}

var taskFilterOptions = []mcp.ToolOption{
	mcp.WithBoolean("archived", mcp.Description("Include archived tasks")),
	mcp.WithBoolean("include_closed", mcp.Description("Include closed tasks")),
	mcp.WithBoolean("subtasks", mcp.Description("Include subtasks")),
	mcp.WithInteger("page", mcp.Description("Page number, 0-indexed. This tool returns one page per call and does not auto-paginate — keep incrementing page and re-calling until a response has fewer tasks than the prior page (or is empty) to know you've seen everything.")),
	mcp.WithString("order_by", mcp.Description("Field to sort by: id, created, updated, due_date")),
	mcp.WithBoolean("reverse", mcp.Description("Reverse sort order")),
	mcp.WithArray("statuses", mcp.WithStringItems(), mcp.Description("Filter by status names")),
	mcp.WithArray("assignees", mcp.WithStringItems(), mcp.Description("Filter by assignee user IDs; a task matches if it has ANY of the given assignees (OR, not AND).")),
	mcp.WithArray("tags", mcp.WithStringItems(), mcp.Description("Filter by tag names")),
	mcp.WithString("custom_fields_json", mcp.Description(`Raw JSON array of custom field filters, ClickUp's custom_fields query param, e.g. [{"field_id":"abc-123","operator":"=","value":"42"}]. Get field_id from clickup_list_list_fields/list_folder_fields/list_space_fields/list_workspace_fields; valid operators (=, !=, >, <, IN, NOT IN, ...) vary by field type.`)),
}

func RegisterTaskTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_get_task",
			mcp.WithDescription("Get a single ClickUp task by ID."),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
			mcp.WithBoolean("custom_task_ids", mcp.Description("Treat task_id as a custom task ID")),
			mcp.WithString("team_id", mcp.Description("Workspace ID; required if custom_task_ids is true, defaults to CLICKUP_TEAM_ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.GetTask(ctx, taskID, req.GetBool("custom_task_ids", false), teamIDOrDefault(req, c))
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_task",
			mcp.WithDescription("Create a task in a ClickUp list."),
			mcp.WithString("list_id", mcp.Required(), mcp.Description("List ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Task name")),
			mcp.WithString("description", mcp.Description("Task description (plain text/markdown)")),
			mcp.WithString("status", mcp.Description("Status name")),
			mcp.WithNumber("priority", mcp.Description("1=Urgent, 2=High, 3=Normal, 4=Low")),
			mcp.WithArray("assignees", mcp.WithStringItems(), mcp.Description("Assignee user IDs")),
			mcp.WithArray("tags", mcp.WithStringItems(), mcp.Description("Tag names")),
			mcp.WithString("due_date", mcp.Description("Due date, UTC \"YYYY-MM-DD HH:MM:SS\" or bare \"YYYY-MM-DD\" (midnight UTC); renders as bare date in responses")),
			mcp.WithString("start_date", mcp.Description("Start date, UTC \"YYYY-MM-DD HH:MM:SS\" or bare \"YYYY-MM-DD\" (midnight UTC); renders as bare date in responses")),
			mcp.WithString("parent", mcp.Description("Parent task ID, to create this as a subtask")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			listID, err := req.RequireString("list_id")
			if err != nil {
				return ErrorResult(err)
			}
			name, err := req.RequireString("name")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{"name": name}
			setString(body, req, "description")
			setString(body, req, "status")
			setFloat(body, req, "priority")
			setStringSlice(body, req, "assignees")
			setStringSlice(body, req, "tags")
			if err := setDateTime(body, req, "due_date"); err != nil {
				return ErrorResult(err)
			}
			if err := setDateTime(body, req, "start_date"); err != nil {
				return ErrorResult(err)
			}
			setString(body, req, "parent")
			out, err := c.CreateTask(ctx, listID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_task",
			mcp.WithDescription("Update a ClickUp task."),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID"+taskIDCaveat)),
			mcp.WithString("name", mcp.Description("Task name")),
			mcp.WithString("description", mcp.Description("Task description (plain text/markdown)")),
			mcp.WithString("status", mcp.Description("Status name")),
			mcp.WithNumber("priority", mcp.Description("1=Urgent, 2=High, 3=Normal, 4=Low")),
			mcp.WithString("due_date", mcp.Description("Due date, UTC \"YYYY-MM-DD HH:MM:SS\" or bare \"YYYY-MM-DD\" (midnight UTC); renders as bare date in responses")),
			mcp.WithString("start_date", mcp.Description("Start date, UTC \"YYYY-MM-DD HH:MM:SS\" or bare \"YYYY-MM-DD\" (midnight UTC); renders as bare date in responses")),
			mcp.WithBoolean("archived", mcp.Description("Archive/unarchive the task")),
			mcp.WithArray("assignees_add", mcp.WithStringItems(), mcp.Description("Assignee user IDs to add")),
			mcp.WithArray("assignees_rem", mcp.WithStringItems(), mcp.Description("Assignee user IDs to remove")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			setString(body, req, "name")
			setString(body, req, "description")
			setString(body, req, "status")
			setFloat(body, req, "priority")
			if err := setDateTime(body, req, "due_date"); err != nil {
				return ErrorResult(err)
			}
			if err := setDateTime(body, req, "start_date"); err != nil {
				return ErrorResult(err)
			}
			setBool(body, req, "archived")
			if hasArg(req, "assignees_add") || hasArg(req, "assignees_rem") {
				assignees := map[string]any{}
				if hasArg(req, "assignees_add") {
					assignees["add"] = req.GetStringSlice("assignees_add", nil)
				}
				if hasArg(req, "assignees_rem") {
					assignees["rem"] = req.GetStringSlice("assignees_rem", nil)
				}
				body["assignees"] = assignees
			}
			out, err := c.UpdateTask(ctx, taskID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_delete_task",
			mcp.WithDescription("Delete a ClickUp task. This is irreversible."),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID"+taskIDCaveat)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.DeleteTask(ctx, taskID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"deleted": true, "task_id": taskID})
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_tasks",
			append([]mcp.ToolOption{
				mcp.WithDescription("List the tasks in a ClickUp list, with optional filters."),
				mcp.WithString("list_id", mcp.Required(), mcp.Description("List ID")),
			}, taskFilterOptions...)...,
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			listID, err := req.RequireString("list_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.GetTasksInList(ctx, listID, taskQueryFilters(req))
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_search_tasks",
			append([]mcp.ToolOption{
				mcp.WithDescription("Search tasks across an entire ClickUp workspace, with optional filters."),
				mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
				mcp.WithArray("space_ids", mcp.WithStringItems(), mcp.Description("Restrict search to these space IDs")),
				mcp.WithArray("list_ids", mcp.WithStringItems(), mcp.Description("Restrict search to these list IDs")),
				mcp.WithArray("folder_ids", mcp.WithStringItems(), mcp.Description("Restrict search to these folder IDs")),
			}, taskFilterOptions...)...,
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			out, err := c.SearchTasks(ctx, teamIDOrDefault(req, c), clickup.SearchTasksOptions{
				TaskQueryFilters: taskQueryFilters(req),
				SpaceIDs:         req.GetStringSlice("space_ids", nil),
				ListIDs:          req.GetStringSlice("list_ids", nil),
				FolderIDs:        req.GetStringSlice("folder_ids", nil),
			})
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_get_task_time_in_status",
			mcp.WithDescription("Get how long a task has spent in each status."),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID"+taskIDCaveat)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.GetTaskTimeInStatus(ctx, taskID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_get_bulk_time_in_status",
			mcp.WithDescription("Get time-in-status for multiple tasks at once."),
			mcp.WithArray("task_ids", mcp.Required(), mcp.WithStringItems(), mcp.Description("Task IDs"+taskIDCaveat)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskIDs, err := req.RequireStringSlice("task_ids")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.GetBulkTimeInStatus(ctx, taskIDs)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_add_task_link",
			mcp.WithDescription("Link two ClickUp tasks together."),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID"+taskIDCaveat)),
			mcp.WithString("links_to", mcp.Required(), mcp.Description("ID of the task to link to"+taskIDCaveat)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			linksTo, err := req.RequireString("links_to")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.AddTaskLink(ctx, taskID, linksTo)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_remove_task_link",
			mcp.WithDescription("Remove a link between two ClickUp tasks."),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID"+taskIDCaveat)),
			mcp.WithString("links_to", mcp.Required(), mcp.Description("ID of the linked task"+taskIDCaveat)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			linksTo, err := req.RequireString("links_to")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.RemoveTaskLink(ctx, taskID, linksTo)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)
}
