package tools

import (
	"net/http"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupListTaskMembers(t *testing.T) {
	t.Run("requires task_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterMemberTools(s, c)

		res := callTool(t, s, "clickup_list_task_members", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires task_id into path", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"members":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterMemberTools(s, c)

		res := callTool(t, s, "clickup_list_task_members", map[string]any{"task_id": "task1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %s, want GET", gotMethod)
		}
		if gotPath != "/task/task1/member" {
			t.Errorf("path = %q, want /task/task1/member", gotPath)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterMemberTools(s, c)

		res := callTool(t, s, "clickup_list_task_members", map[string]any{"task_id": "missing"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupListListMembers(t *testing.T) {
	t.Run("requires list_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterMemberTools(s, c)

		res := callTool(t, s, "clickup_list_list_members", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires list_id into path", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"members":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterMemberTools(s, c)

		res := callTool(t, s, "clickup_list_list_members", map[string]any{"list_id": "list1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %s, want GET", gotMethod)
		}
		if gotPath != "/list/list1/member" {
			t.Errorf("path = %q, want /list/list1/member", gotPath)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterMemberTools(s, c)

		res := callTool(t, s, "clickup_list_list_members", map[string]any{"list_id": "missing"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}
