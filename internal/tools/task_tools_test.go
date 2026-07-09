package tools

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupGetTask(t *testing.T) {
	t.Run("requires task_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_get_task", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires task_id, custom_task_ids and team_id", func(t *testing.T) {
		var gotPath, gotQuery string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"id":"task1","name":"Buy milk"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_get_task", map[string]any{
			"task_id":         "task1",
			"custom_task_ids": true,
			"team_id":         "555",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotPath != "/task/task1" {
			t.Errorf("path = %q, want /task/task1", gotPath)
		}
		if gotQuery != "custom_task_ids=true&team_id=555" {
			t.Errorf("query = %q, want custom_task_ids=true&team_id=555", gotQuery)
		}
		if !strings.Contains(textOf(t, res), "Buy milk") {
			t.Errorf("body = %q, want it to contain Buy milk", textOf(t, res))
		}
	})

	t.Run("defaults team_id to configured default", func(t *testing.T) {
		var gotQuery string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"id":"task1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_get_task", map[string]any{"task_id": "task1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotQuery != "team_id=999" {
			t.Errorf("query = %q, want team_id=999 (default team)", gotQuery)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"Task not found","ECODE":"TASK_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_get_task", map[string]any{"task_id": "missing"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [TASK_001]: Task not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupCreateTask(t *testing.T) {
	t.Run("requires list_id and name", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_create_task", map[string]any{"name": "Buy milk"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing list_id)")
		}

		res = callTool(t, s, "clickup_create_task", map[string]any{"list_id": "list1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing name)")
		}
	})

	t.Run("wires list_id, method and full body", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"task1","name":"Buy milk"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_create_task", map[string]any{
			"list_id":     "list1",
			"name":        "Buy milk",
			"description": "2%, please",
			"status":      "to do",
			"priority":    float64(2),
			"assignees":   []any{"u1", "u2"},
			"tags":        []any{"errand"},
			"due_date":    float64(1700000000000),
			"start_date":  float64(1690000000000),
			"parent":      "task0",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %s, want POST", gotMethod)
		}
		if gotPath != "/list/list1/task" {
			t.Errorf("path = %q, want /list/list1/task", gotPath)
		}
		if gotBody["name"] != "Buy milk" {
			t.Errorf("body[name] = %v, want Buy milk", gotBody["name"])
		}
		if gotBody["description"] != "2%, please" {
			t.Errorf("body[description] = %v", gotBody["description"])
		}
		if gotBody["status"] != "to do" {
			t.Errorf("body[status] = %v", gotBody["status"])
		}
		if gotBody["priority"] != float64(2) {
			t.Errorf("body[priority] = %v", gotBody["priority"])
		}
		if gotBody["parent"] != "task0" {
			t.Errorf("body[parent] = %v", gotBody["parent"])
		}
		assignees, ok := gotBody["assignees"].([]any)
		if !ok || len(assignees) != 2 {
			t.Errorf("body[assignees] = %v", gotBody["assignees"])
		}
	})

	t.Run("minimal body omits unsupplied optional fields", func(t *testing.T) {
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"task1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_create_task", map[string]any{"list_id": "list1", "name": "Buy milk"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if len(gotBody) != 1 {
			t.Fatalf("body = %+v, want exactly {name: Buy milk}", gotBody)
		}
	})
}

func TestClickupUpdateTask(t *testing.T) {
	t.Run("requires task_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_update_task", map[string]any{"status": "done"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("partial update sends only supplied field", func(t *testing.T) {
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"task1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_update_task", map[string]any{"task_id": "task1", "status": "done"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if len(gotBody) != 1 {
			t.Fatalf("body = %+v, want exactly one field", gotBody)
		}
		if gotBody["status"] != "done" {
			t.Errorf("body[status] = %v, want done", gotBody["status"])
		}
	})

	t.Run("wires method and path", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"id":"task1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_update_task", map[string]any{"task_id": "task1", "archived": true})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPut {
			t.Errorf("method = %s, want PUT", gotMethod)
		}
		if gotPath != "/task/task1" {
			t.Errorf("path = %q, want /task/task1", gotPath)
		}
	})

	t.Run("combines assignees_add/rem into an assignees object", func(t *testing.T) {
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"task1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_update_task", map[string]any{
			"task_id":       "task1",
			"assignees_add": []any{"u1"},
			"assignees_rem": []any{"u2"},
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		assignees, ok := gotBody["assignees"].(map[string]any)
		if !ok {
			t.Fatalf("body[assignees] = %v (%T), want map", gotBody["assignees"], gotBody["assignees"])
		}
		add, ok := assignees["add"].([]any)
		if !ok || len(add) != 1 || add[0] != "u1" {
			t.Errorf("assignees[add] = %v", assignees["add"])
		}
		rem, ok := assignees["rem"].([]any)
		if !ok || len(rem) != 1 || rem[0] != "u2" {
			t.Errorf("assignees[rem] = %v", assignees["rem"])
		}
	})

	t.Run("only assignees_add supplied still wires assignees object", func(t *testing.T) {
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"task1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_update_task", map[string]any{
			"task_id":       "task1",
			"assignees_add": []any{"u1"},
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		assignees, ok := gotBody["assignees"].(map[string]any)
		if !ok {
			t.Fatalf("body[assignees] = %v (%T), want map", gotBody["assignees"], gotBody["assignees"])
		}
		if _, present := assignees["rem"]; present {
			t.Errorf("assignees[rem] present = true, want absent")
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"Task not found","ECODE":"TASK_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_update_task", map[string]any{"task_id": "task1", "status": "done"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [TASK_001]: Task not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupDeleteTask(t *testing.T) {
	t.Run("requires task_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_delete_task", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires task_id and reports deleted", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_delete_task", map[string]any{"task_id": "task1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", gotMethod)
		}
		if gotPath != "/task/task1" {
			t.Errorf("path = %q, want /task/task1", gotPath)
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
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_delete_task", map[string]any{"task_id": "task1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 403 [OAUTH_027]: forbidden"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupListTasks(t *testing.T) {
	t.Run("requires list_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_list_tasks", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires filters into query params", func(t *testing.T) {
		var gotPath string
		var gotQuery map[string][]string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotQuery = map[string][]string(r.URL.Query())
			_, _ = w.Write([]byte(`{"tasks":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_list_tasks", map[string]any{
			"list_id":        "list1",
			"archived":       true,
			"include_closed": true,
			"subtasks":       false,
			"page":           float64(2),
			"order_by":       "created",
			"reverse":        true,
			"statuses":       []any{"open", "in progress"},
			"assignees":      []any{"u1"},
			"tags":           []any{"bug"},
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotPath != "/list/list1/task" {
			t.Errorf("path = %q, want /list/list1/task", gotPath)
		}
		if got := gotQuery["archived"]; len(got) != 1 || got[0] != "true" {
			t.Errorf("archived = %v", got)
		}
		if got := gotQuery["subtasks"]; len(got) != 1 || got[0] != "false" {
			t.Errorf("subtasks = %v, want explicit false", got)
		}
		if got := gotQuery["page"]; len(got) != 1 || got[0] != "2" {
			t.Errorf("page = %v", got)
		}
		if got := gotQuery["order_by"]; len(got) != 1 || got[0] != "created" {
			t.Errorf("order_by = %v", got)
		}
		if got := gotQuery["statuses[]"]; len(got) != 2 {
			t.Errorf("statuses[] = %v", got)
		}
		if got := gotQuery["assignees[]"]; len(got) != 1 || got[0] != "u1" {
			t.Errorf("assignees[] = %v", got)
		}
		if got := gotQuery["tags[]"]; len(got) != 1 || got[0] != "bug" {
			t.Errorf("tags[] = %v", got)
		}
	})

	t.Run("omits unset filters", func(t *testing.T) {
		var gotQuery string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"tasks":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_list_tasks", map[string]any{"list_id": "list1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotQuery != "" {
			t.Errorf("query = %q, want empty", gotQuery)
		}
	})
}

func TestClickupSearchTasks(t *testing.T) {
	t.Run("defaults team_id to configured default", func(t *testing.T) {
		var gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"tasks":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_search_tasks", map[string]any{})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotPath != "/team/999/task" {
			t.Errorf("path = %q, want /team/999/task (default team)", gotPath)
		}
	})

	t.Run("explicit team_id overrides default", func(t *testing.T) {
		var gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"tasks":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_search_tasks", map[string]any{"team_id": "777"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotPath != "/team/777/task" {
			t.Errorf("path = %q, want /team/777/task", gotPath)
		}
	})

	t.Run("wires space_ids/list_ids/folder_ids arrays", func(t *testing.T) {
		var gotQuery map[string][]string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotQuery = map[string][]string(r.URL.Query())
			_, _ = w.Write([]byte(`{"tasks":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_search_tasks", map[string]any{
			"space_ids":  []any{"s1", "s2"},
			"list_ids":   []any{"l1"},
			"folder_ids": []any{"f1"},
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if got := gotQuery["space_ids[]"]; len(got) != 2 {
			t.Errorf("space_ids[] = %v", got)
		}
		if got := gotQuery["list_ids[]"]; len(got) != 1 || got[0] != "l1" {
			t.Errorf("list_ids[] = %v", got)
		}
		if got := gotQuery["project_ids[]"]; len(got) != 1 || got[0] != "f1" {
			t.Errorf("project_ids[] = %v, want folder_ids wired to project_ids[]", got)
		}
	})
}

func TestClickupGetTaskTimeInStatus(t *testing.T) {
	t.Run("requires task_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_get_task_time_in_status", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires task_id into path", func(t *testing.T) {
		var gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"current_status":{"status":"in progress"}}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_get_task_time_in_status", map[string]any{"task_id": "task1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotPath != "/task/task1/time_in_status" {
			t.Errorf("path = %q, want /task/task1/time_in_status", gotPath)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"Task not found","ECODE":"TASK_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_get_task_time_in_status", map[string]any{"task_id": "task1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [TASK_001]: Task not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupGetBulkTimeInStatus(t *testing.T) {
	t.Run("requires task_ids", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_get_bulk_time_in_status", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires task_ids into query", func(t *testing.T) {
		var gotPath string
		var gotQuery map[string][]string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotQuery = map[string][]string(r.URL.Query())
			_, _ = w.Write([]byte(`{}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_get_bulk_time_in_status", map[string]any{"task_ids": []any{"t1", "t2"}})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotPath != "/task/bulk_time_in_status" {
			t.Errorf("path = %q, want /task/bulk_time_in_status", gotPath)
		}
		if got := gotQuery["task_ids[]"]; len(got) != 2 || got[0] != "t1" || got[1] != "t2" {
			t.Errorf("task_ids[] = %v", got)
		}
	})
}

func TestClickupAddTaskLink(t *testing.T) {
	t.Run("requires task_id and links_to", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_add_task_link", map[string]any{"task_id": "task1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing links_to)")
		}

		res = callTool(t, s, "clickup_add_task_link", map[string]any{"links_to": "task2"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing task_id)")
		}
	})

	t.Run("wires task_id and links_to into path", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"id":"task1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_add_task_link", map[string]any{"task_id": "task1", "links_to": "task2"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %s, want POST", gotMethod)
		}
		if gotPath != "/task/task1/link/task2" {
			t.Errorf("path = %q, want /task/task1/link/task2", gotPath)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"Task not found","ECODE":"TASK_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_add_task_link", map[string]any{"task_id": "task1", "links_to": "task2"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [TASK_001]: Task not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupRemoveTaskLink(t *testing.T) {
	t.Run("requires task_id and links_to", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_remove_task_link", map[string]any{"task_id": "task1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing links_to)")
		}
	})

	t.Run("wires task_id and links_to into path", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"id":"task1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTaskTools(s, c)

		res := callTool(t, s, "clickup_remove_task_link", map[string]any{"task_id": "task1", "links_to": "task2"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", gotMethod)
		}
		if gotPath != "/task/task1/link/task2" {
			t.Errorf("path = %q, want /task/task1/link/task2", gotPath)
		}
	})
}
