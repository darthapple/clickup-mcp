package tools

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupCreateChecklist(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChecklistTools(s, c)
		res := callTool(t, s, "clickup_create_checklist", map[string]any{"task_id": "task1"})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing required name")
		}
	})

	t.Run("wiring", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("method = %s, want POST", r.Method)
			}
			if r.URL.Path != "/task/task1/checklist" {
				t.Errorf("path = %s, want /task/task1/checklist", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["name"] != "Checklist A" {
				t.Errorf("name = %v, want Checklist A", body["name"])
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"cl1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChecklistTools(s, c)
		res := callTool(t, s, "clickup_create_checklist", map[string]any{"task_id": "task1", "name": "Checklist A"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if !strings.Contains(textOf(t, res), "cl1") {
			t.Errorf("result = %q, want it to contain cl1", textOf(t, res))
		}
	})

	t.Run("error_passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChecklistTools(s, c)
		res := callTool(t, s, "clickup_create_checklist", map[string]any{"task_id": "task1", "name": "x"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("error text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupUpdateChecklist(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChecklistTools(s, c)
		res := callTool(t, s, "clickup_update_checklist", map[string]any{"name": "x"})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing checklist_id")
		}
	})

	t.Run("partial_update_only_position", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("method = %s, want PUT", r.Method)
			}
			if r.URL.Path != "/checklist/cl1" {
				t.Errorf("path = %s, want /checklist/cl1", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if len(body) != 1 {
				t.Errorf("body = %v, want only position set", body)
			}
			if body["position"] != float64(2) {
				t.Errorf("position = %v, want 2", body["position"])
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"cl1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChecklistTools(s, c)
		res := callTool(t, s, "clickup_update_checklist", map[string]any{"checklist_id": "cl1", "position": 2})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})

	t.Run("partial_update_only_name", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if len(body) != 1 {
				t.Errorf("body = %v, want only name set", body)
			}
			if body["name"] != "Renamed" {
				t.Errorf("name = %v, want Renamed", body["name"])
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"cl1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChecklistTools(s, c)
		res := callTool(t, s, "clickup_update_checklist", map[string]any{"checklist_id": "cl1", "name": "Renamed"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})
}

func TestClickupDeleteChecklist(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChecklistTools(s, c)
		res := callTool(t, s, "clickup_delete_checklist", map[string]any{})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing checklist_id")
		}
	})

	t.Run("wiring", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("method = %s, want DELETE", r.Method)
			}
			if r.URL.Path != "/checklist/cl1" {
				t.Errorf("path = %s, want /checklist/cl1", r.URL.Path)
			}
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChecklistTools(s, c)
		res := callTool(t, s, "clickup_delete_checklist", map[string]any{"checklist_id": "cl1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})
}

func TestClickupChecklistItems(t *testing.T) {
	t.Run("create/missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChecklistTools(s, c)
		res := callTool(t, s, "clickup_create_checklist_item", map[string]any{"checklist_id": "cl1"})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing required name")
		}
	})

	t.Run("create/wiring", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("method = %s, want POST", r.Method)
			}
			if r.URL.Path != "/checklist/cl1/checklist_item" {
				t.Errorf("path = %s, want /checklist/cl1/checklist_item", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["name"] != "Item 1" {
				t.Errorf("name = %v, want Item 1", body["name"])
			}
			if body["assignee"] != "user1" {
				t.Errorf("assignee = %v, want user1", body["assignee"])
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"item1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChecklistTools(s, c)
		res := callTool(t, s, "clickup_create_checklist_item", map[string]any{
			"checklist_id": "cl1",
			"name":         "Item 1",
			"assignee":     "user1",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if !strings.Contains(textOf(t, res), "item1") {
			t.Errorf("result = %q, want it to contain item1", textOf(t, res))
		}
	})

	t.Run("update/missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChecklistTools(s, c)
		res := callTool(t, s, "clickup_update_checklist_item", map[string]any{"checklist_id": "cl1"})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing checklist_item_id")
		}
	})

	t.Run("update/partial_update_only_resolved", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("method = %s, want PUT", r.Method)
			}
			if r.URL.Path != "/checklist/cl1/checklist_item/item1" {
				t.Errorf("path = %s, want /checklist/cl1/checklist_item/item1", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if len(body) != 1 {
				t.Errorf("body = %v, want only resolved set", body)
			}
			if body["resolved"] != true {
				t.Errorf("resolved = %v, want true", body["resolved"])
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"item1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChecklistTools(s, c)
		res := callTool(t, s, "clickup_update_checklist_item", map[string]any{
			"checklist_id":      "cl1",
			"checklist_item_id": "item1",
			"resolved":          true,
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})

	t.Run("update/error_passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChecklistTools(s, c)
		res := callTool(t, s, "clickup_update_checklist_item", map[string]any{
			"checklist_id":      "cl1",
			"checklist_item_id": "item1",
			"name":              "x",
		})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("error text = %q, want %q", textOf(t, res), want)
		}
	})

	t.Run("delete/missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChecklistTools(s, c)
		res := callTool(t, s, "clickup_delete_checklist_item", map[string]any{"checklist_id": "cl1"})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing checklist_item_id")
		}
	})

	t.Run("delete/wiring", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("method = %s, want DELETE", r.Method)
			}
			if r.URL.Path != "/checklist/cl1/checklist_item/item1" {
				t.Errorf("path = %s, want /checklist/cl1/checklist_item/item1", r.URL.Path)
			}
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChecklistTools(s, c)
		res := callTool(t, s, "clickup_delete_checklist_item", map[string]any{
			"checklist_id":      "cl1",
			"checklist_item_id": "item1",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})
}
