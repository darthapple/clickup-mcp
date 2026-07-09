package tools

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupCreateFolderFromTemplate(t *testing.T) {
	t.Run("required arg validation", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTemplateTools(s, c)
		res := callTool(t, s, "clickup_create_folder_from_template", map[string]any{"space_id": "s1", "template_id": "t1"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing name)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("argument wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"f1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTemplateTools(s, c)
		res := callTool(t, s, "clickup_create_folder_from_template", map[string]any{
			"space_id":    "s1",
			"template_id": "t1",
			"name":        "New Folder",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/space/s1/folder_template/t1" {
			t.Errorf("path = %q, want /space/s1/folder_template/t1", gotPath)
		}
		if gotBody["name"] != "New Folder" {
			t.Errorf("body = %+v", gotBody)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTemplateTools(s, c)
		res := callTool(t, s, "clickup_create_folder_from_template", map[string]any{
			"space_id": "s1", "template_id": "t1", "name": "New Folder",
		})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupCreateListFromTemplateInFolder(t *testing.T) {
	t.Run("required arg validation", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTemplateTools(s, c)
		res := callTool(t, s, "clickup_create_list_from_template_in_folder", map[string]any{"folder_id": "f1"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing template_id/name)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("argument wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"l1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTemplateTools(s, c)
		res := callTool(t, s, "clickup_create_list_from_template_in_folder", map[string]any{
			"folder_id":   "f1",
			"template_id": "t1",
			"name":        "New List",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/folder/f1/list_template/t1" {
			t.Errorf("path = %q, want /folder/f1/list_template/t1", gotPath)
		}
		if gotBody["name"] != "New List" {
			t.Errorf("body = %+v", gotBody)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTemplateTools(s, c)
		res := callTool(t, s, "clickup_create_list_from_template_in_folder", map[string]any{
			"folder_id": "f1", "template_id": "t1", "name": "New List",
		})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupCreateListFromTemplateInSpace(t *testing.T) {
	t.Run("required arg validation", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTemplateTools(s, c)
		res := callTool(t, s, "clickup_create_list_from_template_in_space", map[string]any{"space_id": "s1"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing template_id/name)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("argument wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"l1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTemplateTools(s, c)
		res := callTool(t, s, "clickup_create_list_from_template_in_space", map[string]any{
			"space_id":    "s1",
			"template_id": "t1",
			"name":        "New List",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/space/s1/list_template/t1" {
			t.Errorf("path = %q, want /space/s1/list_template/t1", gotPath)
		}
		if gotBody["name"] != "New List" {
			t.Errorf("body = %+v", gotBody)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTemplateTools(s, c)
		res := callTool(t, s, "clickup_create_list_from_template_in_space", map[string]any{
			"space_id": "s1", "template_id": "t1", "name": "New List",
		})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupCreateTaskFromTemplate(t *testing.T) {
	t.Run("required arg validation", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTemplateTools(s, c)
		res := callTool(t, s, "clickup_create_task_from_template", map[string]any{"list_id": "l1"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing template_id/name)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("argument wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"task1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTemplateTools(s, c)
		res := callTool(t, s, "clickup_create_task_from_template", map[string]any{
			"list_id":     "l1",
			"template_id": "t1",
			"name":        "New Task",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/list/l1/taskTemplate/t1" {
			t.Errorf("path = %q, want /list/l1/taskTemplate/t1", gotPath)
		}
		if gotBody["name"] != "New Task" {
			t.Errorf("body = %+v", gotBody)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTemplateTools(s, c)
		res := callTool(t, s, "clickup_create_task_from_template", map[string]any{
			"list_id": "l1", "template_id": "t1", "name": "New Task",
		})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}
