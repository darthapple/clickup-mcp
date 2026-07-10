package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func RegisterDocTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_create_doc",
			mcp.WithDescription("Create a Doc in a ClickUp workspace."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Doc name")),
			mcp.WithString("visibility", mcp.Description("PUBLIC, PRIVATE, PERSONAL, or HIDDEN")),
			mcp.WithBoolean("create_page", mcp.Description("Create an initial page; defaults to true")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name, err := req.RequireString("name")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{"name": name}
			setString(body, req, "visibility")
			setBool(body, req, "create_page")
			out, err := c.CreateDoc(ctx, teamIDOrDefault(req, c), body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_search_docs",
			mcp.WithDescription("Search the Docs in a ClickUp workspace. Returns at most 50 "+
				"matching docs; if the workspace has more, only the first 50 are "+
				"returned — this tool does not page further."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("query", mcp.Description("Search text")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			out, err := c.SearchDocs(ctx, teamIDOrDefault(req, c), req.GetString("query", ""))
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_doc_pages",
			mcp.WithDescription("List the pages belonging to a ClickUp Doc."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("doc_id", mcp.Required(), mcp.Description("Doc ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			docID, err := req.RequireString("doc_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListDocPages(ctx, teamIDOrDefault(req, c), docID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_get_doc_page",
			mcp.WithDescription("Get a single page in a ClickUp Doc."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("doc_id", mcp.Required(), mcp.Description("Doc ID")),
			mcp.WithString("page_id", mcp.Required(), mcp.Description("Page ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			docID, err := req.RequireString("doc_id")
			if err != nil {
				return ErrorResult(err)
			}
			pageID, err := req.RequireString("page_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.GetDocPage(ctx, teamIDOrDefault(req, c), docID, pageID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_doc_page",
			mcp.WithDescription("Create a page in a ClickUp Doc."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("doc_id", mcp.Required(), mcp.Description("Doc ID")),
			mcp.WithString("name", mcp.Description("Page name")),
			mcp.WithString("content", mcp.Description("Page content (markdown)")),
			mcp.WithString("parent_page_id", mcp.Description("Parent page ID, to nest this page")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			docID, err := req.RequireString("doc_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			setString(body, req, "name")
			setString(body, req, "content")
			setString(body, req, "parent_page_id")
			out, err := c.CreateDocPage(ctx, teamIDOrDefault(req, c), docID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_doc_page",
			mcp.WithDescription("Edit a page in a ClickUp Doc. IMPORTANT: if content_edit_mode "+
				"is omitted, ClickUp defaults to \"replace\" — supplying content without "+
				"explicitly passing content_edit_mode will silently overwrite/destroy "+
				"the entire existing page, not append to it."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("doc_id", mcp.Required(), mcp.Description("Doc ID")),
			mcp.WithString("page_id", mcp.Required(), mcp.Description("Page ID")),
			mcp.WithString("name", mcp.Description("Page name")),
			mcp.WithString("content", mcp.Description("Page content (markdown)")),
			mcp.WithString("content_edit_mode", mcp.Description("replace (default if omitted — overwrites the entire page), append, or prepend. To add text without destroying existing content, pass \"append\" or \"prepend\" explicitly.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			docID, err := req.RequireString("doc_id")
			if err != nil {
				return ErrorResult(err)
			}
			pageID, err := req.RequireString("page_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			setString(body, req, "name")
			setString(body, req, "content")
			setString(body, req, "content_edit_mode")
			out, err := c.UpdateDocPage(ctx, teamIDOrDefault(req, c), docID, pageID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)
}
