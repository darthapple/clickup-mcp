package tools

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupAddTaskDependency(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDependencyTools(s, c)
		res := callTool(t, s, "clickup_add_task_dependency", map[string]any{"depends_on": "task2"})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing task_id")
		}
	})

	t.Run("wiring_depends_on", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("method = %s, want POST", r.Method)
			}
			if r.URL.Path != "/task/task1/dependency" {
				t.Errorf("path = %s, want /task/task1/dependency", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if len(body) != 1 {
				t.Errorf("body = %v, want only depends_on set", body)
			}
			if body["depends_on"] != "task2" {
				t.Errorf("depends_on = %v, want task2", body["depends_on"])
			}
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDependencyTools(s, c)
		res := callTool(t, s, "clickup_add_task_dependency", map[string]any{"task_id": "task1", "depends_on": "task2"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})

	t.Run("wiring_dependency_of", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if len(body) != 1 {
				t.Errorf("body = %v, want only dependency_of set", body)
			}
			if body["dependency_of"] != "task3" {
				t.Errorf("dependency_of = %v, want task3", body["dependency_of"])
			}
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDependencyTools(s, c)
		res := callTool(t, s, "clickup_add_task_dependency", map[string]any{"task_id": "task1", "dependency_of": "task3"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})

	t.Run("error_passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDependencyTools(s, c)
		res := callTool(t, s, "clickup_add_task_dependency", map[string]any{"task_id": "task1", "depends_on": "task2"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("error text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupRemoveTaskDependency(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDependencyTools(s, c)
		res := callTool(t, s, "clickup_remove_task_dependency", map[string]any{"depends_on": "task2"})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing task_id")
		}
	})

	t.Run("wiring_uses_query_params_not_body", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("method = %s, want DELETE", r.Method)
			}
			if r.URL.Path != "/task/task1/dependency" {
				t.Errorf("path = %s, want /task/task1/dependency", r.URL.Path)
			}
			q := r.URL.Query()
			if q.Get("depends_on") != "task2" {
				t.Errorf("depends_on query = %q, want task2", q.Get("depends_on"))
			}
			if q.Get("dependency_of") != "task3" {
				t.Errorf("dependency_of query = %q, want task3", q.Get("dependency_of"))
			}
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDependencyTools(s, c)
		res := callTool(t, s, "clickup_remove_task_dependency", map[string]any{
			"task_id":       "task1",
			"depends_on":    "task2",
			"dependency_of": "task3",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})

	t.Run("only_supplied_query_param_is_set", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			if q.Get("depends_on") != "task2" {
				t.Errorf("depends_on query = %q, want task2", q.Get("depends_on"))
			}
			if q.Has("dependency_of") {
				t.Errorf("dependency_of query should not be set, got %q", q.Get("dependency_of"))
			}
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDependencyTools(s, c)
		res := callTool(t, s, "clickup_remove_task_dependency", map[string]any{
			"task_id":    "task1",
			"depends_on": "task2",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})

	t.Run("error_passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterDependencyTools(s, c)
		res := callTool(t, s, "clickup_remove_task_dependency", map[string]any{"task_id": "task1", "depends_on": "task2"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("error text = %q, want %q", textOf(t, res), want)
		}
	})
}
