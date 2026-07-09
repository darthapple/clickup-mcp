package tools

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupListFolders(t *testing.T) {
	t.Run("requires space_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterFolderTools(s, c)

		res := callTool(t, s, "clickup_list_folders", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires space_id and archived", func(t *testing.T) {
		var gotMethod, gotPath, gotQuery string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"folders":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterFolderTools(s, c)

		res := callTool(t, s, "clickup_list_folders", map[string]any{"space_id": "space1", "archived": true})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %s, want GET", gotMethod)
		}
		if gotPath != "/space/space1/folder" {
			t.Errorf("path = %q, want /space/space1/folder", gotPath)
		}
		if gotQuery != "archived=true" {
			t.Errorf("query = %q, want archived=true", gotQuery)
		}
	})

	t.Run("omits archived query param when not supplied", func(t *testing.T) {
		var gotQuery string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"folders":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterFolderTools(s, c)

		res := callTool(t, s, "clickup_list_folders", map[string]any{"space_id": "space1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotQuery != "" {
			t.Errorf("query = %q, want empty", gotQuery)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterFolderTools(s, c)

		res := callTool(t, s, "clickup_list_folders", map[string]any{"space_id": "space1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupGetFolder(t *testing.T) {
	t.Run("requires folder_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterFolderTools(s, c)

		res := callTool(t, s, "clickup_get_folder", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires folder_id into path", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"id":"folder1","name":"Sprints"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterFolderTools(s, c)

		res := callTool(t, s, "clickup_get_folder", map[string]any{"folder_id": "folder1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %s, want GET", gotMethod)
		}
		if gotPath != "/folder/folder1" {
			t.Errorf("path = %q, want /folder/folder1", gotPath)
		}
		if !strings.Contains(textOf(t, res), "Sprints") {
			t.Errorf("body = %q, want it to contain Sprints", textOf(t, res))
		}
	})
}

func TestClickupCreateFolder(t *testing.T) {
	t.Run("requires space_id and name", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterFolderTools(s, c)

		res := callTool(t, s, "clickup_create_folder", map[string]any{"space_id": "space1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing name)")
		}

		res = callTool(t, s, "clickup_create_folder", map[string]any{"name": "Sprints"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing space_id)")
		}
	})

	t.Run("wires space_id, method and body", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"folder1","name":"Sprints"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterFolderTools(s, c)

		res := callTool(t, s, "clickup_create_folder", map[string]any{"space_id": "space1", "name": "Sprints"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %s, want POST", gotMethod)
		}
		if gotPath != "/space/space1/folder" {
			t.Errorf("path = %q, want /space/space1/folder", gotPath)
		}
		if len(gotBody) != 1 || gotBody["name"] != "Sprints" {
			t.Errorf("body = %+v, want exactly {name: Sprints}", gotBody)
		}
	})
}

func TestClickupUpdateFolder(t *testing.T) {
	t.Run("requires folder_id and name", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterFolderTools(s, c)

		res := callTool(t, s, "clickup_update_folder", map[string]any{"folder_id": "folder1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing name)")
		}

		res = callTool(t, s, "clickup_update_folder", map[string]any{"name": "Renamed"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing folder_id)")
		}
	})

	t.Run("wires folder_id, method and body", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"folder1","name":"Renamed"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterFolderTools(s, c)

		res := callTool(t, s, "clickup_update_folder", map[string]any{"folder_id": "folder1", "name": "Renamed"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPut {
			t.Errorf("method = %s, want PUT", gotMethod)
		}
		if gotPath != "/folder/folder1" {
			t.Errorf("path = %q, want /folder/folder1", gotPath)
		}
		if len(gotBody) != 1 || gotBody["name"] != "Renamed" {
			t.Errorf("body = %+v, want exactly {name: Renamed}", gotBody)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterFolderTools(s, c)

		res := callTool(t, s, "clickup_update_folder", map[string]any{"folder_id": "folder1", "name": "Renamed"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupDeleteFolder(t *testing.T) {
	t.Run("requires folder_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterFolderTools(s, c)

		res := callTool(t, s, "clickup_delete_folder", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires folder_id and reports deleted", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterFolderTools(s, c)

		res := callTool(t, s, "clickup_delete_folder", map[string]any{"folder_id": "folder1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", gotMethod)
		}
		if gotPath != "/folder/folder1" {
			t.Errorf("path = %q, want /folder/folder1", gotPath)
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
		RegisterFolderTools(s, c)

		res := callTool(t, s, "clickup_delete_folder", map[string]any{"folder_id": "folder1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 403 [OAUTH_027]: forbidden"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}
