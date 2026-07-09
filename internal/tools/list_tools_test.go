package tools

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupListLists(t *testing.T) {
	t.Run("requires folder_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_list_lists", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires folder_id and archived", func(t *testing.T) {
		var gotMethod, gotPath, gotQuery string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"lists":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_list_lists", map[string]any{"folder_id": "folder1", "archived": true})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %s, want GET", gotMethod)
		}
		if gotPath != "/folder/folder1/list" {
			t.Errorf("path = %q, want /folder/folder1/list", gotPath)
		}
		if gotQuery != "archived=true" {
			t.Errorf("query = %q, want archived=true", gotQuery)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_list_lists", map[string]any{"folder_id": "folder1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupListFolderlessLists(t *testing.T) {
	t.Run("requires space_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_list_folderless_lists", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires space_id into path", func(t *testing.T) {
		var gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"lists":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_list_folderless_lists", map[string]any{"space_id": "space1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotPath != "/space/space1/list" {
			t.Errorf("path = %q, want /space/space1/list", gotPath)
		}
	})
}

func TestClickupGetList(t *testing.T) {
	t.Run("requires list_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_get_list", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires list_id into path", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"id":"list1","name":"Backlog"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_get_list", map[string]any{"list_id": "list1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %s, want GET", gotMethod)
		}
		if gotPath != "/list/list1" {
			t.Errorf("path = %q, want /list/list1", gotPath)
		}
		if !strings.Contains(textOf(t, res), "Backlog") {
			t.Errorf("body = %q, want it to contain Backlog", textOf(t, res))
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_get_list", map[string]any{"list_id": "missing"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupCreateListInFolder(t *testing.T) {
	t.Run("requires folder_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_create_list_in_folder", map[string]any{"name": "Sprint 1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires folder_id, method and body", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"list1","name":"Sprint 1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_create_list_in_folder", map[string]any{
			"folder_id": "folder1",
			"name":      "Sprint 1",
			"content":   "First sprint",
			"priority":  float64(3),
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %s, want POST", gotMethod)
		}
		if gotPath != "/folder/folder1/list" {
			t.Errorf("path = %q, want /folder/folder1/list", gotPath)
		}
		if gotBody["name"] != "Sprint 1" || gotBody["content"] != "First sprint" || gotBody["priority"] != float64(3) {
			t.Errorf("body = %+v", gotBody)
		}
		if _, present := gotBody["assignee"]; present {
			t.Errorf("body[assignee] present = true, want absent (not supplied)")
		}
	})
}

func TestClickupCreateFolderlessList(t *testing.T) {
	t.Run("requires space_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_create_folderless_list", map[string]any{"name": "Sprint 1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires space_id into path", func(t *testing.T) {
		var gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"id":"list1","name":"Sprint 1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_create_folderless_list", map[string]any{"space_id": "space1", "name": "Sprint 1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotPath != "/space/space1/list" {
			t.Errorf("path = %q, want /space/space1/list", gotPath)
		}
	})
}

func TestClickupUpdateList(t *testing.T) {
	t.Run("requires list_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_update_list", map[string]any{"name": "Renamed"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("partial update sends only supplied field", func(t *testing.T) {
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"list1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_update_list", map[string]any{
			"list_id": "list1",
			"status":  "green",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if len(gotBody) != 1 {
			t.Fatalf("body = %+v, want exactly one field", gotBody)
		}
		if gotBody["status"] != "green" {
			t.Errorf("body[status] = %v, want green", gotBody["status"])
		}
	})

	t.Run("wires archived and method/path", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"list1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_update_list", map[string]any{"list_id": "list1", "archived": true})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPut {
			t.Errorf("method = %s, want PUT", gotMethod)
		}
		if gotPath != "/list/list1" {
			t.Errorf("path = %q, want /list/list1", gotPath)
		}
		if gotBody["archived"] != true {
			t.Errorf("body[archived] = %v, want true", gotBody["archived"])
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_update_list", map[string]any{"list_id": "list1", "name": "x"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupDeleteList(t *testing.T) {
	t.Run("requires list_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_delete_list", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires list_id and reports deleted", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_delete_list", map[string]any{"list_id": "list1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", gotMethod)
		}
		if gotPath != "/list/list1" {
			t.Errorf("path = %q, want /list/list1", gotPath)
		}
		if !strings.Contains(textOf(t, res), `"deleted":true`) {
			t.Errorf("body = %q, want it to report deleted:true", textOf(t, res))
		}
	})
}

func TestClickupAddTaskToList(t *testing.T) {
	t.Run("requires list_id and task_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_add_task_to_list", map[string]any{"list_id": "list1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing task_id)")
		}

		res = callTool(t, s, "clickup_add_task_to_list", map[string]any{"task_id": "task1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing list_id)")
		}
	})

	t.Run("wires list_id and task_id into path", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusOK)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_add_task_to_list", map[string]any{"list_id": "list1", "task_id": "task1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %s, want POST", gotMethod)
		}
		if gotPath != "/list/list1/task/task1" {
			t.Errorf("path = %q, want /list/list1/task/task1", gotPath)
		}
		if !strings.Contains(textOf(t, res), `"added":true`) {
			t.Errorf("body = %q, want it to report added:true", textOf(t, res))
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_add_task_to_list", map[string]any{"list_id": "list1", "task_id": "task1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupRemoveTaskFromList(t *testing.T) {
	t.Run("requires list_id and task_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_remove_task_from_list", map[string]any{"list_id": "list1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing task_id)")
		}
	})

	t.Run("wires list_id and task_id into path", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusOK)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterListTools(s, c)

		res := callTool(t, s, "clickup_remove_task_from_list", map[string]any{"list_id": "list1", "task_id": "task1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", gotMethod)
		}
		if gotPath != "/list/list1/task/task1" {
			t.Errorf("path = %q, want /list/list1/task/task1", gotPath)
		}
		if !strings.Contains(textOf(t, res), `"removed":true`) {
			t.Errorf("body = %q, want it to report removed:true", textOf(t, res))
		}
	})
}
