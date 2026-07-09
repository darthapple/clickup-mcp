package tools

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

// unreachableHandler fails the test if the fake ClickUp server is hit at
// all, for asserting that required-arg validation short-circuits before any
// HTTP request is made.
func unreachableHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request to fake server: %s %s", r.Method, r.URL.String())
	}
}

func TestClickupListSpaces(t *testing.T) {
	t.Run("defaults team_id to configured default", func(t *testing.T) {
		var gotPath, gotQuery string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"spaces":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterSpaceTools(s, c)

		res := callTool(t, s, "clickup_list_spaces", map[string]any{})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotPath != "/team/999/space" {
			t.Errorf("path = %q, want /team/999/space", gotPath)
		}
		if gotQuery != "" {
			t.Errorf("query = %q, want empty (archived not supplied)", gotQuery)
		}
	})

	t.Run("explicit team_id overrides default", func(t *testing.T) {
		var gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"spaces":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterSpaceTools(s, c)

		res := callTool(t, s, "clickup_list_spaces", map[string]any{"team_id": "12345"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotPath != "/team/12345/space" {
			t.Errorf("path = %q, want /team/12345/space", gotPath)
		}
	})

	t.Run("wires archived flag", func(t *testing.T) {
		var gotQuery string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"spaces":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterSpaceTools(s, c)

		res := callTool(t, s, "clickup_list_spaces", map[string]any{"archived": true})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotQuery != "archived=true" {
			t.Errorf("query = %q, want archived=true", gotQuery)
		}
	})
}

func TestClickupGetSpace(t *testing.T) {
	t.Run("requires space_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterSpaceTools(s, c)

		res := callTool(t, s, "clickup_get_space", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires space_id into path", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"id":"space1","name":"Engineering"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterSpaceTools(s, c)

		res := callTool(t, s, "clickup_get_space", map[string]any{"space_id": "space1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %s, want GET", gotMethod)
		}
		if gotPath != "/space/space1" {
			t.Errorf("path = %q, want /space/space1", gotPath)
		}
		if !strings.Contains(textOf(t, res), "Engineering") {
			t.Errorf("body = %q, want it to contain Engineering", textOf(t, res))
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterSpaceTools(s, c)

		res := callTool(t, s, "clickup_get_space", map[string]any{"space_id": "missing"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupCreateSpace(t *testing.T) {
	t.Run("requires name", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterSpaceTools(s, c)

		res := callTool(t, s, "clickup_create_space", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires team_id, name and optional fields", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"space2","name":"New Space"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterSpaceTools(s, c)

		res := callTool(t, s, "clickup_create_space", map[string]any{
			"team_id":            "555",
			"name":               "New Space",
			"multiple_assignees": true,
			"statuses_json":      `[{"status":"to do","color":"#000","type":"open","orderindex":0}]`,
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %s, want POST", gotMethod)
		}
		if gotPath != "/team/555/space" {
			t.Errorf("path = %q, want /team/555/space", gotPath)
		}
		if gotBody["name"] != "New Space" {
			t.Errorf("body[name] = %v, want New Space", gotBody["name"])
		}
		if gotBody["multiple_assignees"] != true {
			t.Errorf("body[multiple_assignees] = %v, want true", gotBody["multiple_assignees"])
		}
		statuses, ok := gotBody["statuses"].([]any)
		if !ok || len(statuses) != 1 {
			t.Errorf("body[statuses] = %v, want a single-element array", gotBody["statuses"])
		}
	})

	t.Run("invalid statuses_json fails before hitting server", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterSpaceTools(s, c)

		res := callTool(t, s, "clickup_create_space", map[string]any{
			"name":          "Bad Space",
			"statuses_json": `not json`,
		})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})
}

func TestClickupUpdateSpace(t *testing.T) {
	t.Run("requires space_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterSpaceTools(s, c)

		res := callTool(t, s, "clickup_update_space", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("partial update sends only supplied field", func(t *testing.T) {
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"space1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterSpaceTools(s, c)

		res := callTool(t, s, "clickup_update_space", map[string]any{
			"space_id": "space1",
			"name":     "Renamed",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if len(gotBody) != 1 {
			t.Fatalf("body = %+v, want exactly one field", gotBody)
		}
		if gotBody["name"] != "Renamed" {
			t.Errorf("body[name] = %v, want Renamed", gotBody["name"])
		}
	})

	t.Run("wires space_id into path and method", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"id":"space1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterSpaceTools(s, c)

		res := callTool(t, s, "clickup_update_space", map[string]any{
			"space_id": "space1",
			"archived": true,
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPut {
			t.Errorf("method = %s, want PUT", gotMethod)
		}
		if gotPath != "/space/space1" {
			t.Errorf("path = %q, want /space/space1", gotPath)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterSpaceTools(s, c)

		res := callTool(t, s, "clickup_update_space", map[string]any{"space_id": "space1", "name": "x"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupDeleteSpace(t *testing.T) {
	t.Run("requires space_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterSpaceTools(s, c)

		res := callTool(t, s, "clickup_delete_space", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires space_id and reports deleted", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterSpaceTools(s, c)

		res := callTool(t, s, "clickup_delete_space", map[string]any{"space_id": "space1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", gotMethod)
		}
		if gotPath != "/space/space1" {
			t.Errorf("path = %q, want /space/space1", gotPath)
		}
		if !strings.Contains(textOf(t, res), `"deleted":true`) {
			t.Errorf("body = %q, want it to report deleted:true", textOf(t, res))
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"err":"forbidden","ECODE":"OAUTH_027"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterSpaceTools(s, c)

		res := callTool(t, s, "clickup_delete_space", map[string]any{"space_id": "space1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 403 [OAUTH_027]: forbidden"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}
