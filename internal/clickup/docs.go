package clickup

import (
	"context"
	"net/http"
	"net/url"
)

func docsBase(workspaceID string) string {
	return "/workspaces/" + workspaceID + "/docs"
}

// CreateDoc creates a Doc in a workspace.
// POST /v3/workspaces/{workspace_id}/docs
func (c *Client) CreateDoc(ctx context.Context, workspaceID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, APIVersion: apiV3, Path: docsBase(workspaceID), Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SearchDocs searches the Docs in a workspace.
// GET /v3/workspaces/{workspace_id}/docs
func (c *Client) SearchDocs(ctx context.Context, workspaceID, query string) (any, error) {
	q := url.Values{}
	addParam(q, "query", query)
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, APIVersion: apiV3, Path: docsBase(workspaceID), Query: q}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListDocPages returns the pages belonging to a Doc.
// GET /v3/workspaces/{workspace_id}/docs/{doc_id}/pages
func (c *Client) ListDocPages(ctx context.Context, workspaceID, docID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, APIVersion: apiV3, Path: docsBase(workspaceID) + "/" + docID + "/pages"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDocPage returns a single page in a Doc.
// GET /v3/workspaces/{workspace_id}/docs/{doc_id}/pages/{page_id}
func (c *Client) GetDocPage(ctx context.Context, workspaceID, docID, pageID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, APIVersion: apiV3, Path: docsBase(workspaceID) + "/" + docID + "/pages/" + pageID}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateDocPage creates a page in a Doc.
// POST /v3/workspaces/{workspace_id}/docs/{doc_id}/pages
func (c *Client) CreateDocPage(ctx context.Context, workspaceID, docID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, APIVersion: apiV3, Path: docsBase(workspaceID) + "/" + docID + "/pages", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateDocPage edits a page in a Doc.
// PUT /v3/workspaces/{workspace_id}/docs/{doc_id}/pages/{page_id}
func (c *Client) UpdateDocPage(ctx context.Context, workspaceID, docID, pageID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPut, APIVersion: apiV3, Path: docsBase(workspaceID) + "/" + docID + "/pages/" + pageID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
