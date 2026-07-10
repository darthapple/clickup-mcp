package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

// commentPaginationCaveat documents that the comment-listing endpoints in
// this file call ClickUp with no start/start_id pagination params, so only
// the most recent page (ClickUp's server-side default limit) is ever
// reachable — on a task/list/view with more comments than that, older ones
// are permanently unreachable through this tool with no signal that any
// were left out.
const commentPaginationCaveat = " Returns only the most recent page of comments (ClickUp's server-side default page limit); older comments beyond that cannot be retrieved through this tool since start/start_id pagination isn't exposed."

func buildCommentBody(req mcp.CallToolRequest) (map[string]any, error) {
	commentText, err := req.RequireString("comment_text")
	if err != nil {
		return nil, err
	}
	body := map[string]any{"comment_text": commentText}
	setBool(body, req, "notify_all")
	setString(body, req, "assignee")
	return body, nil
}

func RegisterCommentTools(s *server.MCPServer, c *clickup.Client) {
	commentBodyOptions := []mcp.ToolOption{
		mcp.WithString("comment_text", mcp.Required(), mcp.Description("Comment text")),
		mcp.WithBoolean("notify_all", mcp.Description("Notify all task watchers")),
		mcp.WithString("assignee", mcp.Description("User ID to assign the comment to")),
	}

	s.AddTool(
		mcp.NewTool("clickup_list_task_comments",
			mcp.WithDescription("List comments on a ClickUp task."+commentPaginationCaveat),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID"+taskIDCaveat)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListTaskComments(ctx, taskID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_task_comment",
			append([]mcp.ToolOption{
				mcp.WithDescription("Add a comment to a ClickUp task."),
				mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID"+taskIDCaveat)),
			}, commentBodyOptions...)...,
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return ErrorResult(err)
			}
			body, err := buildCommentBody(req)
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.CreateTaskComment(ctx, taskID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_list_comments",
			mcp.WithDescription("List comments on a ClickUp list."+commentPaginationCaveat),
			mcp.WithString("list_id", mcp.Required(), mcp.Description("List ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			listID, err := req.RequireString("list_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListListComments(ctx, listID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_list_comment",
			append([]mcp.ToolOption{
				mcp.WithDescription("Add a comment to a ClickUp list."),
				mcp.WithString("list_id", mcp.Required(), mcp.Description("List ID")),
			}, commentBodyOptions...)...,
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			listID, err := req.RequireString("list_id")
			if err != nil {
				return ErrorResult(err)
			}
			body, err := buildCommentBody(req)
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.CreateListComment(ctx, listID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_view_comments",
			mcp.WithDescription("List comments on a ClickUp view."+commentPaginationCaveat),
			mcp.WithString("view_id", mcp.Required(), mcp.Description("View ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			viewID, err := req.RequireString("view_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListViewComments(ctx, viewID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_view_comment",
			append([]mcp.ToolOption{
				mcp.WithDescription("Add a comment to a ClickUp view."),
				mcp.WithString("view_id", mcp.Required(), mcp.Description("View ID")),
			}, commentBodyOptions...)...,
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			viewID, err := req.RequireString("view_id")
			if err != nil {
				return ErrorResult(err)
			}
			body, err := buildCommentBody(req)
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.CreateViewComment(ctx, viewID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_comment",
			mcp.WithDescription("Update a ClickUp comment."),
			mcp.WithString("comment_id", mcp.Required(), mcp.Description("Comment ID")),
			mcp.WithString("comment_text", mcp.Description("New comment text")),
			mcp.WithString("assignee", mcp.Description("User ID to assign the comment to")),
			mcp.WithBoolean("resolved", mcp.Description("Mark the comment resolved/unresolved")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			commentID, err := req.RequireString("comment_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			setString(body, req, "comment_text")
			setString(body, req, "assignee")
			setBool(body, req, "resolved")
			out, err := c.UpdateComment(ctx, commentID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_delete_comment",
			mcp.WithDescription("Delete a ClickUp comment."),
			mcp.WithString("comment_id", mcp.Required(), mcp.Description("Comment ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			commentID, err := req.RequireString("comment_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.DeleteComment(ctx, commentID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"deleted": true, "comment_id": commentID})
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_comment_replies",
			mcp.WithDescription("List threaded replies on a ClickUp comment."+commentPaginationCaveat),
			mcp.WithString("comment_id", mcp.Required(), mcp.Description("Comment ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			commentID, err := req.RequireString("comment_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListCommentReplies(ctx, commentID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_comment_reply",
			append([]mcp.ToolOption{
				mcp.WithDescription("Add a threaded reply to a ClickUp comment."),
				mcp.WithString("comment_id", mcp.Required(), mcp.Description("Comment ID")),
			}, commentBodyOptions...)...,
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			commentID, err := req.RequireString("comment_id")
			if err != nil {
				return ErrorResult(err)
			}
			body, err := buildCommentBody(req)
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.CreateCommentReply(ctx, commentID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)
}
