package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

// Note: every method in docs.go sets APIVersion: apiV3; since testClient
// points BaseURLv2 and BaseURLv3 at the same fake server, that's exercised
// implicitly by every request below hitting the handler at all.

func TestCreateDocPostsExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"d1"}`))
	})

	out, err := c.CreateDoc(context.Background(), "999", map[string]any{"name": "Runbook"})
	if err != nil {
		t.Fatalf("CreateDoc: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotPath != "/workspaces/999/docs" {
		t.Errorf("path = %q, want /workspaces/999/docs", gotPath)
	}
	if gotBody["name"] != "Runbook" {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok || m["id"] != "d1" {
		t.Errorf("decoded = %+v", out)
	}
}

func TestSearchDocsIncludesQueryParamWhenSet(t *testing.T) {
	var gotPath, gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"docs":[{"id":"d1"}]}`))
	})

	out, err := c.SearchDocs(context.Background(), "999", "runbook")
	if err != nil {
		t.Fatalf("SearchDocs: %v", err)
	}
	if gotPath != "/workspaces/999/docs" {
		t.Errorf("path = %q, want /workspaces/999/docs", gotPath)
	}
	if gotQuery != "query=runbook" {
		t.Errorf("query = %q, want query=runbook", gotQuery)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded type = %T, want map[string]any", out)
	}
	if docs, ok := m["docs"].([]any); !ok || len(docs) != 1 {
		t.Errorf("docs = %+v", m["docs"])
	}
}

func TestSearchDocsOmitsQueryParamWhenEmpty(t *testing.T) {
	var gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"docs":[]}`))
	})

	if _, err := c.SearchDocs(context.Background(), "999", ""); err != nil {
		t.Fatalf("SearchDocs: %v", err)
	}
	if gotQuery != "" {
		t.Errorf("query = %q, want empty (query param omitted)", gotQuery)
	}
}

func TestListDocPagesHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"pages":[{"id":"p1"}]}`))
	})

	out, err := c.ListDocPages(context.Background(), "999", "d1")
	if err != nil {
		t.Fatalf("ListDocPages: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/workspaces/999/docs/d1/pages" {
		t.Errorf("path = %q, want /workspaces/999/docs/d1/pages", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded type = %T, want map[string]any", out)
	}
	if pages, ok := m["pages"].([]any); !ok || len(pages) != 1 {
		t.Errorf("pages = %+v", m["pages"])
	}
}

func TestGetDocPageHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"id":"p1","content":"hello"}`))
	})

	out, err := c.GetDocPage(context.Background(), "999", "d1", "p1")
	if err != nil {
		t.Fatalf("GetDocPage: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/workspaces/999/docs/d1/pages/p1" {
		t.Errorf("path = %q, want /workspaces/999/docs/d1/pages/p1", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok || m["content"] != "hello" {
		t.Errorf("decoded = %+v", out)
	}
}

func TestCreateDocPagePostsExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"p1"}`))
	})

	out, err := c.CreateDocPage(context.Background(), "999", "d1", map[string]any{"name": "Page 1"})
	if err != nil {
		t.Fatalf("CreateDocPage: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotPath != "/workspaces/999/docs/d1/pages" {
		t.Errorf("path = %q, want /workspaces/999/docs/d1/pages", gotPath)
	}
	if gotBody["name"] != "Page 1" {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok || m["id"] != "p1" {
		t.Errorf("decoded = %+v", out)
	}
}

func TestUpdateDocPagePutsExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"p1","content":"updated"}`))
	})

	out, err := c.UpdateDocPage(context.Background(), "999", "d1", "p1", map[string]any{"content": "updated"})
	if err != nil {
		t.Fatalf("UpdateDocPage: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %q, want PUT", gotMethod)
	}
	if gotPath != "/workspaces/999/docs/d1/pages/p1" {
		t.Errorf("path = %q, want /workspaces/999/docs/d1/pages/p1", gotPath)
	}
	if gotBody["content"] != "updated" {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok || m["content"] != "updated" {
		t.Errorf("decoded = %+v", out)
	}
}

func TestCreateDocReturnsAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"err":"Invalid doc name","ECODE":"DOC_001"}`))
	})

	_, err := c.CreateDoc(context.Background(), "999", map[string]any{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest || apiErr.ECode != "DOC_001" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
