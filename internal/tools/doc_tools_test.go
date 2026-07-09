package tools

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupCreateDoc(t *testing.T) {
	t.Run("required arg validation", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { hit = true })
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_create_doc", map[string]any{"team_id": "123"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing name)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("argument wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"d1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_create_doc", map[string]any{
			"team_id":     "123",
			"name":        "My Doc",
			"visibility":  "PUBLIC",
			"create_page": true,
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/workspaces/123/docs" {
			t.Errorf("path = %q, want /workspaces/123/docs", gotPath)
		}
		if gotBody["name"] != "My Doc" || gotBody["visibility"] != "PUBLIC" || gotBody["create_page"] != true {
			t.Errorf("body = %+v", gotBody)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_create_doc", map[string]any{"team_id": "123", "name": "My Doc"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupSearchDocs(t *testing.T) {
	t.Run("argument wiring with query", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotQuery string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotQuery = r.URL.Query().Get("query")
			_, _ = w.Write([]byte(`{"docs":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_search_docs", map[string]any{"team_id": "123", "query": "roadmap"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %q, want GET", gotMethod)
		}
		if gotPath != "/workspaces/123/docs" {
			t.Errorf("path = %q, want /workspaces/123/docs", gotPath)
		}
		if gotQuery != "roadmap" {
			t.Errorf("query = %q, want roadmap", gotQuery)
		}
	})

	t.Run("no query arg is allowed (query is optional)", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			_, _ = w.Write([]byte(`{"docs":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_search_docs", map[string]any{"team_id": "123"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if !hit {
			t.Error("expected handler to be hit for an optional-arg tool with no required args missing")
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_search_docs", map[string]any{"team_id": "123"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupListDocPages(t *testing.T) {
	t.Run("required arg validation", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { hit = true })
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_list_doc_pages", map[string]any{"team_id": "123"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing doc_id)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("argument wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"pages":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_list_doc_pages", map[string]any{"team_id": "123", "doc_id": "d1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %q, want GET", gotMethod)
		}
		if gotPath != "/workspaces/123/docs/d1/pages" {
			t.Errorf("path = %q, want /workspaces/123/docs/d1/pages", gotPath)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_list_doc_pages", map[string]any{"team_id": "123", "doc_id": "d1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupGetDocPage(t *testing.T) {
	t.Run("required arg validation", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { hit = true })
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_get_doc_page", map[string]any{"team_id": "123", "doc_id": "d1"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing page_id)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("argument wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"id":"p1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_get_doc_page", map[string]any{"team_id": "123", "doc_id": "d1", "page_id": "p1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %q, want GET", gotMethod)
		}
		if gotPath != "/workspaces/123/docs/d1/pages/p1" {
			t.Errorf("path = %q, want /workspaces/123/docs/d1/pages/p1", gotPath)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_get_doc_page", map[string]any{"team_id": "123", "doc_id": "d1", "page_id": "p1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupCreateDocPage(t *testing.T) {
	t.Run("required arg validation", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { hit = true })
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_create_doc_page", map[string]any{"team_id": "123"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing doc_id)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("argument wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"p1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_create_doc_page", map[string]any{
			"team_id":        "123",
			"doc_id":         "d1",
			"name":           "Page 1",
			"content":        "hello",
			"parent_page_id": "p0",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/workspaces/123/docs/d1/pages" {
			t.Errorf("path = %q, want /workspaces/123/docs/d1/pages", gotPath)
		}
		if gotBody["name"] != "Page 1" || gotBody["content"] != "hello" || gotBody["parent_page_id"] != "p0" {
			t.Errorf("body = %+v", gotBody)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_create_doc_page", map[string]any{"team_id": "123", "doc_id": "d1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupUpdateDocPage(t *testing.T) {
	t.Run("required arg validation", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { hit = true })
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_update_doc_page", map[string]any{"team_id": "123", "doc_id": "d1"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing page_id)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("partial update semantics", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"p1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_update_doc_page", map[string]any{
			"team_id": "123",
			"doc_id":  "d1",
			"page_id": "p1",
			"content": "new content",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPut {
			t.Errorf("method = %q, want PUT", gotMethod)
		}
		if gotPath != "/workspaces/123/docs/d1/pages/p1" {
			t.Errorf("path = %q, want /workspaces/123/docs/d1/pages/p1", gotPath)
		}
		if len(gotBody) != 1 {
			t.Errorf("body = %+v, want exactly one field", gotBody)
		}
		if gotBody["content"] != "new content" {
			t.Errorf("body[content] = %v, want %q", gotBody["content"], "new content")
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDocTools(s, c)
		res := callTool(t, s, "clickup_update_doc_page", map[string]any{"team_id": "123", "doc_id": "d1", "page_id": "p1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}
