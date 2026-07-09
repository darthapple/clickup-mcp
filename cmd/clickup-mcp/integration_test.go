package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/mcptest"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
	"clickup-mcp/internal/config"
	"clickup-mcp/internal/tools"
)

// newIntegrationServer builds a real MCP server (tools.RegisterAll applied,
// exactly as main.go does) backed by a fake ClickUp API, and drives it
// through an in-process client over real stdio pipes via mcptest — the full
// RegisterAll -> server.MCPServer -> JSON-RPC -> handler path, rather than
// invoking a tool handler directly (see internal/tools' callTool helper,
// which bypasses the transport). The network boundary is faked, so this
// proves our own layers are wired together correctly; it does not prove
// anything about ClickUp's real API behavior — see e2e_test.go (build tag
// e2e) for that.
func newIntegrationServer(t *testing.T, handler http.HandlerFunc) *mcptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	c := clickup.NewClient(&config.Config{
		APIToken:    "pk_test_token",
		TeamID:      "999",
		BaseURLv2:   srv.URL,
		BaseURLv3:   srv.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  0,
	})

	s := server.NewMCPServer("clickup-mcp-integration", "test")
	tools.RegisterAll(s, c)

	registered := s.ListTools()
	serverTools := make([]server.ServerTool, 0, len(registered))
	for _, st := range registered {
		serverTools = append(serverTools, *st)
	}

	mts, err := mcptest.NewServer(t, serverTools...)
	if err != nil {
		t.Fatalf("mcptest.NewServer: %v", err)
	}
	t.Cleanup(mts.Close)
	return mts
}

// callToolText calls name on mts and returns its text content plus whether
// the result is an error result.
func callToolText(t *testing.T, mts *mcptest.Server, name string, args map[string]any) (string, bool) {
	t.Helper()
	var req mcp.CallToolRequest
	req.Params.Name = name
	req.Params.Arguments = args

	res, err := mts.Client().CallTool(context.Background(), req)
	if err != nil {
		t.Fatalf("CallTool(%s): %v", name, err)
	}
	for _, c := range res.Content {
		if tc, ok := c.(mcp.TextContent); ok {
			return tc.Text, res.IsError
		}
	}
	return "", res.IsError
}

func TestIntegrationToolsListIsNonEmpty(t *testing.T) {
	mts := newIntegrationServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	})
	listed, err := mts.Client().ListTools(context.Background(), mcp.ListToolsRequest{})
	if err != nil {
		t.Fatalf("ListTools: %v", err)
	}
	if len(listed.Tools) == 0 {
		t.Fatal("expected at least one tool, got none")
	}
}

func TestIntegrationSimpleGet(t *testing.T) {
	mts := newIntegrationServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/task/abc" {
			t.Errorf("path = %s, want /task/abc", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"abc","name":"Test task"}`))
	})

	text, isErr := callToolText(t, mts, "clickup_get_task", map[string]any{"task_id": "abc"})
	if isErr {
		t.Fatalf("unexpected error: %s", text)
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(text), &out); err != nil {
		t.Fatalf("decoding result: %v", err)
	}
	if out["id"] != "abc" {
		t.Errorf("id = %v, want abc", out["id"])
	}
}

func TestIntegrationPostWithBody(t *testing.T) {
	var gotMethod string
	var gotBody map[string]any
	mts := newIntegrationServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"new1"}`))
	})

	text, isErr := callToolText(t, mts, "clickup_create_task", map[string]any{
		"list_id": "list1",
		"name":    "Buy milk",
	})
	if isErr {
		t.Fatalf("unexpected error: %s", text)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotBody["name"] != "Buy milk" {
		t.Errorf("body[name] = %v, want %q", gotBody["name"], "Buy milk")
	}
}

func TestIntegrationArrayQueryParamSearch(t *testing.T) {
	mts := newIntegrationServer(t, func(w http.ResponseWriter, r *http.Request) {
		got := r.URL.Query()["statuses[]"]
		if len(got) != 2 || got[0] != "open" || got[1] != "in progress" {
			t.Errorf("statuses[] = %v", got)
		}
		_, _ = w.Write([]byte(`{"tasks":[]}`))
	})

	_, isErr := callToolText(t, mts, "clickup_search_tasks", map[string]any{
		"statuses": []any{"open", "in progress"},
	})
	if isErr {
		t.Fatal("unexpected error from clickup_search_tasks")
	}
}

func TestIntegrationGeneratedGuestScopeTool(t *testing.T) {
	mts := newIntegrationServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/list/list1/guest/g1" {
			t.Errorf("path = %s, want /list/list1/guest/g1", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"guest":{}}`))
	})

	_, isErr := callToolText(t, mts, "clickup_add_guest_to_list", map[string]any{
		"list_id":  "list1",
		"guest_id": "g1",
	})
	if isErr {
		t.Fatal("unexpected error from clickup_add_guest_to_list")
	}
}

func TestIntegrationErrorPassthrough(t *testing.T) {
	mts := newIntegrationServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"err":"Task not found","ECODE":"TASK_001"}`))
	})

	text, isErr := callToolText(t, mts, "clickup_get_task", map[string]any{"task_id": "missing"})
	if !isErr {
		t.Fatal("expected an error result")
	}
	want := "ClickUp API error 404 [TASK_001]: Task not found"
	if text != want {
		t.Errorf("error text = %q, want %q", text, want)
	}
}
