package tools

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupTaskTag(t *testing.T) {
	cases := []struct {
		name       string
		tool       string
		wantMethod string
	}{
		{"add", "clickup_add_task_tag", http.MethodPost},
		{"remove", "clickup_remove_task_tag", http.MethodDelete},
	}

	for _, tc := range cases {
		t.Run(tc.name+"/missing_required_arg", func(t *testing.T) {
			hit := false
			c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				hit = true
				t.Errorf("server should not have been called")
			})
			s := server.NewMCPServer("test", "1.0.0")
			RegisterTagTools(s, c)
			res := callTool(t, s, tc.tool, map[string]any{"task_id": "task1"})
			if !res.IsError {
				t.Error("IsError = false, want true")
			}
			if hit {
				t.Error("server was hit despite missing tag_name")
			}
		})

		t.Run(tc.name+"/wiring", func(t *testing.T) {
			c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tc.wantMethod {
					t.Errorf("method = %s, want %s", r.Method, tc.wantMethod)
				}
				if r.URL.Path != "/task/task1/tag/urgent" {
					t.Errorf("path = %s, want /task/task1/tag/urgent", r.URL.Path)
				}
				w.WriteHeader(http.StatusNoContent)
			})
			s := server.NewMCPServer("test", "1.0.0")
			RegisterTagTools(s, c)
			res := callTool(t, s, tc.tool, map[string]any{"task_id": "task1", "tag_name": "urgent"})
			if res.IsError {
				t.Fatalf("IsError = true, want false: %s", textOf(t, res))
			}
			if !strings.Contains(textOf(t, res), "urgent") {
				t.Errorf("result = %q, want it to contain tag_name", textOf(t, res))
			}
		})

		t.Run(tc.name+"/error_passthrough", func(t *testing.T) {
			c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
			})
			s := server.NewMCPServer("test", "1.0.0")
			RegisterTagTools(s, c)
			res := callTool(t, s, tc.tool, map[string]any{"task_id": "task1", "tag_name": "urgent"})
			if !res.IsError {
				t.Fatal("IsError = false, want true")
			}
			want := "ClickUp API error 404 [X_001]: not found"
			if textOf(t, res) != want {
				t.Errorf("error text = %q, want %q", textOf(t, res), want)
			}
		})
	}
}

func TestClickupListSpaceTags(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTagTools(s, c)
		res := callTool(t, s, "clickup_list_space_tags", map[string]any{})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing space_id")
		}
	})

	t.Run("wiring", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("method = %s, want GET", r.Method)
			}
			if r.URL.Path != "/space/space1/tag" {
				t.Errorf("path = %s, want /space/space1/tag", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"tags":[{"name":"urgent"}]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTagTools(s, c)
		res := callTool(t, s, "clickup_list_space_tags", map[string]any{"space_id": "space1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if !strings.Contains(textOf(t, res), "urgent") {
			t.Errorf("result = %q, want it to contain urgent", textOf(t, res))
		}
	})
}

func TestClickupCreateSpaceTag(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTagTools(s, c)
		res := callTool(t, s, "clickup_create_space_tag", map[string]any{"space_id": "space1"})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing name")
		}
	})

	t.Run("wiring_wraps_body_under_tag", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("method = %s, want POST", r.Method)
			}
			if r.URL.Path != "/space/space1/tag" {
				t.Errorf("path = %s, want /space/space1/tag", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			tag, ok := body["tag"].(map[string]any)
			if !ok {
				t.Fatalf("body[tag] = %v, want a nested object", body["tag"])
			}
			if tag["name"] != "urgent" {
				t.Errorf("tag.name = %v, want urgent", tag["name"])
			}
			if tag["tag_fg"] != "#000000" {
				t.Errorf("tag.tag_fg = %v, want #000000", tag["tag_fg"])
			}
			if tag["tag_bg"] != "#ffffff" {
				t.Errorf("tag.tag_bg = %v, want #ffffff", tag["tag_bg"])
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"tag":{"name":"urgent"}}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTagTools(s, c)
		res := callTool(t, s, "clickup_create_space_tag", map[string]any{
			"space_id": "space1",
			"name":     "urgent",
			"tag_fg":   "#000000",
			"tag_bg":   "#ffffff",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})

	t.Run("only_name_set_when_colors_omitted", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			tag, ok := body["tag"].(map[string]any)
			if !ok {
				t.Fatalf("body[tag] = %v, want a nested object", body["tag"])
			}
			if len(tag) != 1 {
				t.Errorf("tag = %v, want only name set", tag)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"tag":{"name":"urgent"}}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTagTools(s, c)
		res := callTool(t, s, "clickup_create_space_tag", map[string]any{"space_id": "space1", "name": "urgent"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})
}

func TestClickupUpdateSpaceTag(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTagTools(s, c)
		res := callTool(t, s, "clickup_update_space_tag", map[string]any{"space_id": "space1"})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing tag_name")
		}
	})

	t.Run("partial_update_only_name", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("method = %s, want PUT", r.Method)
			}
			if r.URL.Path != "/space/space1/tag/urgent" {
				t.Errorf("path = %s, want /space/space1/tag/urgent", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			tag, ok := body["tag"].(map[string]any)
			if !ok {
				t.Fatalf("body[tag] = %v, want a nested object", body["tag"])
			}
			if len(tag) != 1 || tag["name"] != "renamed" {
				t.Errorf("tag = %v, want only name=renamed", tag)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"tag":{"name":"renamed"}}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTagTools(s, c)
		res := callTool(t, s, "clickup_update_space_tag", map[string]any{
			"space_id": "space1",
			"tag_name": "urgent",
			"name":     "renamed",
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
		RegisterTagTools(s, c)
		res := callTool(t, s, "clickup_update_space_tag", map[string]any{
			"space_id": "space1",
			"tag_name": "urgent",
			"name":     "renamed",
		})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("error text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupDeleteSpaceTag(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTagTools(s, c)
		res := callTool(t, s, "clickup_delete_space_tag", map[string]any{"space_id": "space1"})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing tag_name")
		}
	})

	t.Run("wiring", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("method = %s, want DELETE", r.Method)
			}
			if r.URL.Path != "/space/space1/tag/urgent" {
				t.Errorf("path = %s, want /space/space1/tag/urgent", r.URL.Path)
			}
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTagTools(s, c)
		res := callTool(t, s, "clickup_delete_space_tag", map[string]any{"space_id": "space1", "tag_name": "urgent"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})
}
