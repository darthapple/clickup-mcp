package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func RegisterListTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_list_lists",
			mcp.WithDescription("List the lists inside a ClickUp folder."),
			mcp.WithString("folder_id", mcp.Required(), mcp.Description("Folder ID")),
			mcp.WithBoolean("archived", mcp.Description("Include archived lists")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			folderID, err := req.RequireString("folder_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListListsInFolder(ctx, folderID, req.GetBool("archived", false), hasArg(req, "archived"))
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_folderless_lists",
			mcp.WithDescription("List the folderless lists directly under a ClickUp space."),
			mcp.WithString("space_id", mcp.Required(), mcp.Description("Space ID")),
			mcp.WithBoolean("archived", mcp.Description("Include archived lists")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			spaceID, err := req.RequireString("space_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListFolderlessLists(ctx, spaceID, req.GetBool("archived", false), hasArg(req, "archived"))
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_get_list",
			mcp.WithDescription("Get a single ClickUp list by ID."),
			mcp.WithString("list_id", mcp.Required(), mcp.Description("List ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			listID, err := req.RequireString("list_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.GetList(ctx, listID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	listBodyOptions := []mcp.ToolOption{
		mcp.WithString("name", mcp.Required(), mcp.Description("List name")),
		mcp.WithString("content", mcp.Description("List description")),
		mcp.WithNumber("priority", mcp.Description("1=Urgent, 2=High, 3=Normal, 4=Low")),
		mcp.WithString("assignee", mcp.Description("User ID to assign the list to")),
		mcp.WithString("status", mcp.Description("List status color/name")),
	}

	s.AddTool(
		mcp.NewTool("clickup_create_list_in_folder",
			append([]mcp.ToolOption{
				mcp.WithDescription("Create a list inside a ClickUp folder."),
				mcp.WithString("folder_id", mcp.Required(), mcp.Description("Folder ID")),
			}, listBodyOptions...)...,
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			folderID, err := req.RequireString("folder_id")
			if err != nil {
				return ErrorResult(err)
			}
			body, err := buildListBody(req)
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.CreateListInFolder(ctx, folderID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_folderless_list",
			append([]mcp.ToolOption{
				mcp.WithDescription("Create a list directly under a ClickUp space (no folder)."),
				mcp.WithString("space_id", mcp.Required(), mcp.Description("Space ID")),
			}, listBodyOptions...)...,
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			spaceID, err := req.RequireString("space_id")
			if err != nil {
				return ErrorResult(err)
			}
			body, err := buildListBody(req)
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.CreateFolderlessList(ctx, spaceID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_list",
			append([]mcp.ToolOption{
				mcp.WithDescription("Update a ClickUp list."),
				mcp.WithString("list_id", mcp.Required(), mcp.Description("List ID")),
				mcp.WithBoolean("archived", mcp.Description("Archive/unarchive the list")),
			}, listBodyOptions...)...,
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			listID, err := req.RequireString("list_id")
			if err != nil {
				return ErrorResult(err)
			}
			body, err := buildListBody(req)
			if err != nil {
				return ErrorResult(err)
			}
			setBool(body, req, "archived")
			out, err := c.UpdateList(ctx, listID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_delete_list",
			mcp.WithDescription("Delete a ClickUp list. This is irreversible."),
			mcp.WithString("list_id", mcp.Required(), mcp.Description("List ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			listID, err := req.RequireString("list_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.DeleteList(ctx, listID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"deleted": true, "list_id": listID})
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_add_task_to_list",
			mcp.WithDescription("Add an existing task to an additional list."),
			mcp.WithString("list_id", mcp.Required(), mcp.Description("List ID")),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			listID, err := req.RequireString("list_id")
			if err != nil {
				return ErrorResult(err)
			}
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.AddTaskToList(ctx, listID, taskID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"added": true, "list_id": listID, "task_id": taskID})
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_remove_task_from_list",
			mcp.WithDescription("Remove a task from an additional list (does not delete the task)."),
			mcp.WithString("list_id", mcp.Required(), mcp.Description("List ID")),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			listID, err := req.RequireString("list_id")
			if err != nil {
				return ErrorResult(err)
			}
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.RemoveTaskFromList(ctx, listID, taskID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"removed": true, "list_id": listID, "task_id": taskID})
		},
	)
}

func buildListBody(req mcp.CallToolRequest) (map[string]any, error) {
	body := map[string]any{}
	setString(body, req, "name")
	setString(body, req, "content")
	setFloat(body, req, "priority")
	setString(body, req, "assignee")
	setString(body, req, "status")
	return body, nil
}
