package tools

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupListWebhooks(t *testing.T) {
	t.Run("defaults team_id and wires path", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"webhooks":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterWebhookTools(s, c)

		res := callTool(t, s, "clickup_list_webhooks", map[string]any{})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %s, want GET", gotMethod)
		}
		if gotPath != "/team/999/webhook" {
			t.Errorf("path = %q, want /team/999/webhook", gotPath)
		}
	})

	t.Run("explicit team_id overrides default", func(t *testing.T) {
		var gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"webhooks":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterWebhookTools(s, c)

		res := callTool(t, s, "clickup_list_webhooks", map[string]any{"team_id": "team42"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotPath != "/team/team42/webhook" {
			t.Errorf("path = %q, want /team/team42/webhook", gotPath)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterWebhookTools(s, c)

		res := callTool(t, s, "clickup_list_webhooks", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupCreateWebhook(t *testing.T) {
	t.Run("requires endpoint and events", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterWebhookTools(s, c)

		res := callTool(t, s, "clickup_create_webhook", map[string]any{"events": []any{"taskCreated"}})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing endpoint)")
		}
		res = callTool(t, s, "clickup_create_webhook", map[string]any{"endpoint": "https://example.com/hook"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing events)")
		}
	})

	t.Run("wires body scoping fields and method/path", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"wh1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterWebhookTools(s, c)

		res := callTool(t, s, "clickup_create_webhook", map[string]any{
			"endpoint": "https://example.com/hook",
			"events":   []any{"taskCreated", "taskUpdated"},
			"space_id": "space1",
			"list_id":  "list1",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %s, want POST", gotMethod)
		}
		if gotPath != "/team/999/webhook" {
			t.Errorf("path = %q, want /team/999/webhook", gotPath)
		}
		if gotBody["endpoint"] != "https://example.com/hook" {
			t.Errorf("body[endpoint] = %v", gotBody["endpoint"])
		}
		if gotBody["space_id"] != "space1" || gotBody["list_id"] != "list1" {
			t.Errorf("body = %+v", gotBody)
		}
		if _, present := gotBody["folder_id"]; present {
			t.Errorf("body[folder_id] present = true, want absent (not supplied)")
		}
		events, ok := gotBody["events"].([]any)
		if !ok || len(events) != 2 || events[0] != "taskCreated" {
			t.Errorf("body[events] = %v", gotBody["events"])
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"err":"bad","ECODE":"X_002"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterWebhookTools(s, c)

		res := callTool(t, s, "clickup_create_webhook", map[string]any{
			"endpoint": "https://example.com/hook",
			"events":   []any{"*"},
		})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})
}

func TestClickupUpdateWebhook(t *testing.T) {
	t.Run("requires webhook_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterWebhookTools(s, c)

		res := callTool(t, s, "clickup_update_webhook", map[string]any{"status": "suspended"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("partial update sends only supplied field", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"wh1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterWebhookTools(s, c)

		res := callTool(t, s, "clickup_update_webhook", map[string]any{
			"webhook_id": "wh1",
			"status":     "suspended",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPut {
			t.Errorf("method = %s, want PUT", gotMethod)
		}
		if gotPath != "/webhook/wh1" {
			t.Errorf("path = %q, want /webhook/wh1", gotPath)
		}
		if len(gotBody) != 1 || gotBody["status"] != "suspended" {
			t.Errorf("body = %+v, want only status=suspended", gotBody)
		}
	})

	t.Run("wires events slice", func(t *testing.T) {
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"wh1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterWebhookTools(s, c)

		res := callTool(t, s, "clickup_update_webhook", map[string]any{
			"webhook_id": "wh1",
			"events":     []any{"taskDeleted"},
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		events, ok := gotBody["events"].([]any)
		if !ok || len(events) != 1 || events[0] != "taskDeleted" {
			t.Errorf("body[events] = %v", gotBody["events"])
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterWebhookTools(s, c)

		res := callTool(t, s, "clickup_update_webhook", map[string]any{"webhook_id": "wh1", "status": "active"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupDeleteWebhook(t *testing.T) {
	t.Run("requires webhook_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterWebhookTools(s, c)

		res := callTool(t, s, "clickup_delete_webhook", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires webhook_id and reports deleted", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterWebhookTools(s, c)

		res := callTool(t, s, "clickup_delete_webhook", map[string]any{"webhook_id": "wh1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", gotMethod)
		}
		if gotPath != "/webhook/wh1" {
			t.Errorf("path = %q, want /webhook/wh1", gotPath)
		}
		var out map[string]any
		if err := json.Unmarshal([]byte(textOf(t, res)), &out); err != nil {
			t.Fatalf("decoding result: %v", err)
		}
		if out["deleted"] != true || out["webhook_id"] != "wh1" {
			t.Errorf("result = %+v", out)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterWebhookTools(s, c)

		res := callTool(t, s, "clickup_delete_webhook", map[string]any{"webhook_id": "missing"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})
}
