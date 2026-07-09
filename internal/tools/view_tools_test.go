package tools

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupListTeamViews(t *testing.T) {
	t.Run("defaults team_id and wires path", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"views":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_list_team_views", map[string]any{})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %s, want GET", gotMethod)
		}
		if gotPath != "/team/999/view" {
			t.Errorf("path = %q, want /team/999/view", gotPath)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_list_team_views", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupListSpaceViews(t *testing.T) {
	t.Run("requires space_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_list_space_views", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires space_id into path", func(t *testing.T) {
		var gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"views":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_list_space_views", map[string]any{"space_id": "space1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotPath != "/space/space1/view" {
			t.Errorf("path = %q, want /space/space1/view", gotPath)
		}
	})
}

func TestClickupCreateSpaceView(t *testing.T) {
	t.Run("requires space_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_create_space_view", map[string]any{"name": "Board"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing space_id)")
		}
	})

	t.Run("requires name", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_create_space_view", map[string]any{"space_id": "space1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing name)")
		}
	})

	t.Run("wires body and method/path", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"view1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_create_space_view", map[string]any{
			"space_id": "space1",
			"name":     "Board",
			"type":     "board",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %s, want POST", gotMethod)
		}
		if gotPath != "/space/space1/view" {
			t.Errorf("path = %q, want /space/space1/view", gotPath)
		}
		if gotBody["name"] != "Board" || gotBody["type"] != "board" {
			t.Errorf("body = %+v", gotBody)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"err":"bad","ECODE":"X_002"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_create_space_view", map[string]any{"space_id": "space1", "name": "Board"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})
}

func TestClickupListFolderViews(t *testing.T) {
	t.Run("requires folder_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_list_folder_views", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires folder_id into path", func(t *testing.T) {
		var gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"views":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_list_folder_views", map[string]any{"folder_id": "folder1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotPath != "/folder/folder1/view" {
			t.Errorf("path = %q, want /folder/folder1/view", gotPath)
		}
	})
}

func TestClickupCreateFolderView(t *testing.T) {
	t.Run("requires folder_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_create_folder_view", map[string]any{"name": "Timeline"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing folder_id)")
		}
	})

	t.Run("requires name", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_create_folder_view", map[string]any{"folder_id": "folder1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing name)")
		}
	})

	t.Run("wires body and method/path", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"view1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_create_folder_view", map[string]any{
			"folder_id": "folder1",
			"name":      "Timeline",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %s, want POST", gotMethod)
		}
		if gotPath != "/folder/folder1/view" {
			t.Errorf("path = %q, want /folder/folder1/view", gotPath)
		}
		if gotBody["name"] != "Timeline" {
			t.Errorf("body[name] = %v, want Timeline", gotBody["name"])
		}
		if _, present := gotBody["type"]; present {
			t.Errorf("body[type] present = true, want absent (not supplied)")
		}
	})
}

func TestClickupListListViews(t *testing.T) {
	t.Run("requires list_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_list_list_views", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires list_id into path", func(t *testing.T) {
		var gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"views":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_list_list_views", map[string]any{"list_id": "list1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotPath != "/list/list1/view" {
			t.Errorf("path = %q, want /list/list1/view", gotPath)
		}
	})
}

func TestClickupCreateListView(t *testing.T) {
	t.Run("requires list_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_create_list_view", map[string]any{"name": "Sprint board"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing list_id)")
		}
	})

	t.Run("requires name", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_create_list_view", map[string]any{"list_id": "list1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing name)")
		}
	})

	t.Run("wires body and method/path", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"view1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_create_list_view", map[string]any{
			"list_id": "list1",
			"name":    "Sprint board",
			"type":    "board",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %s, want POST", gotMethod)
		}
		if gotPath != "/list/list1/view" {
			t.Errorf("path = %q, want /list/list1/view", gotPath)
		}
		if gotBody["name"] != "Sprint board" || gotBody["type"] != "board" {
			t.Errorf("body = %+v", gotBody)
		}
	})
}

func TestClickupGetView(t *testing.T) {
	t.Run("requires view_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_get_view", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires view_id into path", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"id":"view1","name":"Board"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_get_view", map[string]any{"view_id": "view1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %s, want GET", gotMethod)
		}
		if gotPath != "/view/view1" {
			t.Errorf("path = %q, want /view/view1", gotPath)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_get_view", map[string]any{"view_id": "missing"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupUpdateView(t *testing.T) {
	t.Run("requires view_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_update_view", map[string]any{"name": "Renamed"})
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
			_, _ = w.Write([]byte(`{"id":"view1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_update_view", map[string]any{"view_id": "view1", "name": "Renamed"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPut {
			t.Errorf("method = %s, want PUT", gotMethod)
		}
		if gotPath != "/view/view1" {
			t.Errorf("path = %q, want /view/view1", gotPath)
		}
		if len(gotBody) != 1 || gotBody["name"] != "Renamed" {
			t.Errorf("body = %+v, want only name=Renamed", gotBody)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_update_view", map[string]any{"view_id": "view1", "name": "x"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupDeleteView(t *testing.T) {
	t.Run("requires view_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_delete_view", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires view_id and reports deleted", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_delete_view", map[string]any{"view_id": "view1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", gotMethod)
		}
		if gotPath != "/view/view1" {
			t.Errorf("path = %q, want /view/view1", gotPath)
		}
		var out map[string]any
		if err := json.Unmarshal([]byte(textOf(t, res)), &out); err != nil {
			t.Fatalf("decoding result: %v", err)
		}
		if out["deleted"] != true || out["view_id"] != "view1" {
			t.Errorf("result = %+v", out)
		}
	})
}

func TestClickupGetViewTasks(t *testing.T) {
	t.Run("requires view_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_get_view_tasks", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("omits page query when not supplied", func(t *testing.T) {
		var gotPath string
		var gotQuery url.Values
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotQuery = r.URL.Query()
			_, _ = w.Write([]byte(`{"tasks":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_get_view_tasks", map[string]any{"view_id": "view1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotPath != "/view/view1/task" {
			t.Errorf("path = %q, want /view/view1/task", gotPath)
		}
		if gotQuery.Get("page") != "" {
			t.Errorf("page query = %q, want empty when not supplied", gotQuery.Get("page"))
		}
	})

	t.Run("wires page query when supplied, including page 0", func(t *testing.T) {
		var gotQuery url.Values
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotQuery = r.URL.Query()
			_, _ = w.Write([]byte(`{"tasks":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterViewTools(s, c)

		res := callTool(t, s, "clickup_get_view_tasks", map[string]any{"view_id": "view1", "page": float64(0)})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotQuery.Get("page") != "0" {
			t.Errorf("page query = %q, want 0 (explicit page 0 must still be sent)", gotQuery.Get("page"))
		}

		res = callTool(t, s, "clickup_get_view_tasks", map[string]any{"view_id": "view1", "page": float64(3)})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotQuery.Get("page") != "3" {
			t.Errorf("page query = %q, want 3", gotQuery.Get("page"))
		}
	})
}
