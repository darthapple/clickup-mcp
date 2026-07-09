package tools

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupListFields(t *testing.T) {
	cases := []struct {
		name     string
		tool     string
		idArg    string
		idVal    string
		wantPath string
	}{
		{"list", "clickup_list_list_fields", "list_id", "list1", "/list/list1/field"},
		{"folder", "clickup_list_folder_fields", "folder_id", "folder1", "/folder/folder1/field"},
		{"space", "clickup_list_space_fields", "space_id", "space1", "/space/space1/field"},
	}

	for _, tc := range cases {
		t.Run(tc.name+"/missing_required_arg", func(t *testing.T) {
			hit := false
			c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				hit = true
				t.Errorf("server should not have been called")
			})
			s := server.NewMCPServer("test", "1.0.0")
			RegisterCustomFieldTools(s, c)
			res := callTool(t, s, tc.tool, map[string]any{})
			if !res.IsError {
				t.Error("IsError = false, want true")
			}
			if hit {
				t.Error("server was hit despite missing required arg")
			}
		})

		t.Run(tc.name+"/wiring", func(t *testing.T) {
			c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("method = %s, want GET", r.Method)
				}
				if r.URL.Path != tc.wantPath {
					t.Errorf("path = %s, want %s", r.URL.Path, tc.wantPath)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"fields":[{"id":"f1"}]}`))
			})
			s := server.NewMCPServer("test", "1.0.0")
			RegisterCustomFieldTools(s, c)
			res := callTool(t, s, tc.tool, map[string]any{tc.idArg: tc.idVal})
			if res.IsError {
				t.Fatalf("IsError = true, want false: %s", textOf(t, res))
			}
			if !strings.Contains(textOf(t, res), "f1") {
				t.Errorf("result = %q, want it to contain f1", textOf(t, res))
			}
		})
	}
}

func TestClickupListWorkspaceFields(t *testing.T) {
	t.Run("defaults_team_id", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/team/999/field" {
				t.Errorf("path = %s, want /team/999/field (default team from config)", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"fields":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterCustomFieldTools(s, c)
		res := callTool(t, s, "clickup_list_workspace_fields", map[string]any{})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})

	t.Run("explicit_team_id", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/team/team42/field" {
				t.Errorf("path = %s, want /team/team42/field", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"fields":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterCustomFieldTools(s, c)
		res := callTool(t, s, "clickup_list_workspace_fields", map[string]any{"team_id": "team42"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})
}

func TestClickupSetTaskCustomField(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterCustomFieldTools(s, c)
		res := callTool(t, s, "clickup_set_task_custom_field", map[string]any{"task_id": "task1", "field_id": "f1"})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing required value_json")
		}
	})

	t.Run("wiring_string_value", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("method = %s, want POST", r.Method)
			}
			if r.URL.Path != "/task/task1/field/f1" {
				t.Errorf("path = %s, want /task/task1/field/f1", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["value"] != "hello" {
				t.Errorf("value = %v, want hello", body["value"])
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"f1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterCustomFieldTools(s, c)
		res := callTool(t, s, "clickup_set_task_custom_field", map[string]any{
			"task_id":    "task1",
			"field_id":   "f1",
			"value_json": `"hello"`,
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})

	t.Run("wiring_number_value", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["value"] != float64(42) {
				t.Errorf("value = %v, want 42", body["value"])
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"f1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterCustomFieldTools(s, c)
		res := callTool(t, s, "clickup_set_task_custom_field", map[string]any{
			"task_id":    "task1",
			"field_id":   "f1",
			"value_json": `42`,
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})

	t.Run("wiring_array_value", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			arr, ok := body["value"].([]any)
			if !ok || len(arr) != 2 || arr[0] != "uuid1" || arr[1] != "uuid2" {
				t.Errorf("value = %v, want [uuid1 uuid2]", body["value"])
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"f1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterCustomFieldTools(s, c)
		res := callTool(t, s, "clickup_set_task_custom_field", map[string]any{
			"task_id":    "task1",
			"field_id":   "f1",
			"value_json": `["uuid1","uuid2"]`,
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})

	t.Run("invalid_json_produces_error_not_panic", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called for invalid JSON")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterCustomFieldTools(s, c)
		res := callTool(t, s, "clickup_set_task_custom_field", map[string]any{
			"task_id":    "task1",
			"field_id":   "f1",
			"value_json": `{not valid json`,
		})
		if !res.IsError {
			t.Fatal("IsError = false, want true for invalid JSON")
		}
		if hit {
			t.Error("server was hit despite invalid value_json")
		}
		if !strings.Contains(textOf(t, res), "value_json") {
			t.Errorf("error text = %q, want it to mention value_json", textOf(t, res))
		}
	})

	t.Run("error_passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterCustomFieldTools(s, c)
		res := callTool(t, s, "clickup_set_task_custom_field", map[string]any{
			"task_id":    "task1",
			"field_id":   "f1",
			"value_json": `1`,
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

func TestClickupRemoveTaskCustomField(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterCustomFieldTools(s, c)
		res := callTool(t, s, "clickup_remove_task_custom_field", map[string]any{"task_id": "task1"})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing field_id")
		}
	})

	t.Run("wiring", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("method = %s, want DELETE", r.Method)
			}
			if r.URL.Path != "/task/task1/field/f1" {
				t.Errorf("path = %s, want /task/task1/field/f1", r.URL.Path)
			}
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterCustomFieldTools(s, c)
		res := callTool(t, s, "clickup_remove_task_custom_field", map[string]any{"task_id": "task1", "field_id": "f1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})
}
