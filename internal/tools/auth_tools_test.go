package tools

import (
	"net/http"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupGetUser(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"user":{"id":1,"username":"fernando"}}`))
	})
	s := server.NewMCPServer("test", "1.0.0")
	RegisterAuthTools(s, c)

	res := callTool(t, s, "clickup_get_user", map[string]any{})
	if res.IsError {
		t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/user" {
		t.Errorf("path = %q, want /user", gotPath)
	}
	if got := textOf(t, res); !strings.Contains(got, "fernando") {
		t.Errorf("body = %q, want it to contain fernando", got)
	}
}

func TestClickupGetUserErrorPassthrough(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"err":"Invalid token","ECODE":"OAUTH_001"}`))
	})
	s := server.NewMCPServer("test", "1.0.0")
	RegisterAuthTools(s, c)

	res := callTool(t, s, "clickup_get_user", map[string]any{})
	if !res.IsError {
		t.Fatal("IsError = false, want true")
	}
	text := textOf(t, res)
	want := "ClickUp API error 401 [OAUTH_001]: Invalid token"
	if text != want {
		t.Errorf("error text = %q, want %q", text, want)
	}
}

func TestClickupListWorkspaces(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"teams":[{"id":"999","name":"Kheperi"}]}`))
	})
	s := server.NewMCPServer("test", "1.0.0")
	RegisterAuthTools(s, c)

	res := callTool(t, s, "clickup_list_workspaces", map[string]any{})
	if res.IsError {
		t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/team" {
		t.Errorf("path = %q, want /team", gotPath)
	}
	if got := textOf(t, res); !strings.Contains(got, "Kheperi") {
		t.Errorf("body = %q, want it to contain Kheperi", got)
	}
}

func TestClickupListWorkspacesErrorPassthrough(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
	})
	s := server.NewMCPServer("test", "1.0.0")
	RegisterAuthTools(s, c)

	res := callTool(t, s, "clickup_list_workspaces", map[string]any{})
	if !res.IsError {
		t.Fatal("IsError = false, want true")
	}
	text := textOf(t, res)
	want := "ClickUp API error 404 [X_001]: not found"
	if text != want {
		t.Errorf("error text = %q, want %q", text, want)
	}
}
