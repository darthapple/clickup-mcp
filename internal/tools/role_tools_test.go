package tools

import (
	"net/http"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupListCustomRoles(t *testing.T) {
	t.Run("argument wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"roles":[{"id":1}]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterRoleTools(s, c)
		res := callTool(t, s, "clickup_list_custom_roles", map[string]any{"team_id": "123"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %q, want GET", gotMethod)
		}
		if gotPath != "/team/123/customroles" {
			t.Errorf("path = %q, want /team/123/customroles", gotPath)
		}
	})

	t.Run("defaults team_id from config", func(t *testing.T) {
		var gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"roles":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterRoleTools(s, c)
		res := callTool(t, s, "clickup_list_custom_roles", map[string]any{})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotPath != "/team/999/customroles" {
			t.Errorf("path = %q, want /team/999/customroles (default team id)", gotPath)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterRoleTools(s, c)
		res := callTool(t, s, "clickup_list_custom_roles", map[string]any{"team_id": "123"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}
