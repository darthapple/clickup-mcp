package tools

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupInviteWorkspaceUser(t *testing.T) {
	t.Run("required arg validation", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterUserTools(s, c)
		res := callTool(t, s, "clickup_invite_workspace_user", map[string]any{})
		if !res.IsError {
			t.Error("IsError = false, want true (missing email)")
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
			_, _ = w.Write([]byte(`{"id":"u1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterUserTools(s, c)
		res := callTool(t, s, "clickup_invite_workspace_user", map[string]any{
			"team_id": "123",
			"email":   "a@b.com",
			"admin":   true,
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/team/123/user" {
			t.Errorf("path = %q, want /team/123/user", gotPath)
		}
		if gotBody["email"] != "a@b.com" || gotBody["admin"] != true {
			t.Errorf("body = %+v", gotBody)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterUserTools(s, c)
		res := callTool(t, s, "clickup_invite_workspace_user", map[string]any{"email": "a@b.com"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupUpdateWorkspaceUser(t *testing.T) {
	t.Run("required arg validation", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterUserTools(s, c)
		res := callTool(t, s, "clickup_update_workspace_user", map[string]any{})
		if !res.IsError {
			t.Error("IsError = false, want true (missing user_id)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("partial update semantics", func(t *testing.T) {
		var gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"u1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterUserTools(s, c)
		res := callTool(t, s, "clickup_update_workspace_user", map[string]any{
			"team_id":  "123",
			"user_id":  "u1",
			"username": "bob",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotPath != "/team/123/user/u1" {
			t.Errorf("path = %q, want /team/123/user/u1", gotPath)
		}
		if len(gotBody) != 1 {
			t.Errorf("body = %+v, want exactly one field", gotBody)
		}
		if gotBody["username"] != "bob" {
			t.Errorf("body[username] = %v, want bob", gotBody["username"])
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterUserTools(s, c)
		res := callTool(t, s, "clickup_update_workspace_user", map[string]any{"user_id": "u1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupRemoveWorkspaceUser(t *testing.T) {
	t.Run("required arg validation", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterUserTools(s, c)
		res := callTool(t, s, "clickup_remove_workspace_user", map[string]any{})
		if !res.IsError {
			t.Error("IsError = false, want true (missing user_id)")
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
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterUserTools(s, c)
		res := callTool(t, s, "clickup_remove_workspace_user", map[string]any{
			"team_id": "123",
			"user_id": "u1",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/team/123/user/u1" {
			t.Errorf("path = %q, want /team/123/user/u1", gotPath)
		}
		if !strings.Contains(textOf(t, res), `"removed":true`) {
			t.Errorf("text = %q, want removed:true", textOf(t, res))
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterUserTools(s, c)
		res := callTool(t, s, "clickup_remove_workspace_user", map[string]any{"user_id": "u1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}
