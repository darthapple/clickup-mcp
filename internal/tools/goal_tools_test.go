package tools

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupListGoals(t *testing.T) {
	t.Run("defaults_team_id", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("method = %s, want GET", r.Method)
			}
			if r.URL.Path != "/team/999/goal" {
				t.Errorf("path = %s, want /team/999/goal (default team from config)", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"goals":[{"id":"g1"}]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_list_goals", map[string]any{})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if !strings.Contains(textOf(t, res), "g1") {
			t.Errorf("result = %q, want it to contain g1", textOf(t, res))
		}
	})

	t.Run("explicit_team_id", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/team/team42/goal" {
				t.Errorf("path = %s, want /team/team42/goal", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"goals":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_list_goals", map[string]any{"team_id": "team42"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})
}

func TestClickupCreateGoal(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_create_goal", map[string]any{})
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
			if r.URL.Path != "/team/999/goal" {
				t.Errorf("path = %s, want /team/999/goal", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["name"] != "Q3 Goal" {
				t.Errorf("name = %v, want Q3 Goal", body["name"])
			}
			if body["due_date"] != float64(123000) {
				t.Errorf("due_date = %v, want 123000", body["due_date"])
			}
			if body["multiple_owners"] != true {
				t.Errorf("multiple_owners = %v, want true", body["multiple_owners"])
			}
			owners, ok := body["owners"].([]any)
			if !ok || len(owners) != 2 || owners[0] != "u1" || owners[1] != "u2" {
				t.Errorf("owners = %v, want [u1 u2]", body["owners"])
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"goal1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_create_goal", map[string]any{
			"name":            "Q3 Goal",
			"due_date":        "1970-01-01 00:02:03",
			"multiple_owners": true,
			"owners":          []any{"u1", "u2"},
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if !strings.Contains(textOf(t, res), "goal1") {
			t.Errorf("result = %q, want it to contain goal1", textOf(t, res))
		}
	})

	t.Run("only_name_in_body_when_optional_fields_omitted", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if len(body) != 1 {
				t.Errorf("body = %v, want only name set", body)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"goal1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_create_goal", map[string]any{"name": "Q3 Goal"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})
}

func TestClickupGetGoal(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_get_goal", map[string]any{})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing goal_id")
		}
	})

	t.Run("wiring", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("method = %s, want GET", r.Method)
			}
			if r.URL.Path != "/goal/goal1" {
				t.Errorf("path = %s, want /goal/goal1", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"goal1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_get_goal", map[string]any{"goal_id": "goal1"})
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
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_get_goal", map[string]any{"goal_id": "goal1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("error text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupUpdateGoal(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_update_goal", map[string]any{"name": "x"})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing goal_id")
		}
	})

	t.Run("partial_update_only_color", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("method = %s, want PUT", r.Method)
			}
			if r.URL.Path != "/goal/goal1" {
				t.Errorf("path = %s, want /goal/goal1", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if len(body) != 1 {
				t.Errorf("body = %v, want only color set", body)
			}
			if body["color"] != "#ff0000" {
				t.Errorf("color = %v, want #ff0000", body["color"])
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"goal1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_update_goal", map[string]any{"goal_id": "goal1", "color": "#ff0000"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})

	t.Run("wiring_add_rem_owners", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			addOwners, ok := body["add_owners"].([]any)
			if !ok || len(addOwners) != 1 || addOwners[0] != "u1" {
				t.Errorf("add_owners = %v, want [u1]", body["add_owners"])
			}
			remOwners, ok := body["rem_owners"].([]any)
			if !ok || len(remOwners) != 1 || remOwners[0] != "u2" {
				t.Errorf("rem_owners = %v, want [u2]", body["rem_owners"])
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"goal1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_update_goal", map[string]any{
			"goal_id":    "goal1",
			"add_owners": []any{"u1"},
			"rem_owners": []any{"u2"},
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
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_update_goal", map[string]any{"goal_id": "goal1", "name": "x"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("error text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupDeleteGoal(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_delete_goal", map[string]any{})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing goal_id")
		}
	})

	t.Run("wiring", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("method = %s, want DELETE", r.Method)
			}
			if r.URL.Path != "/goal/goal1" {
				t.Errorf("path = %s, want /goal/goal1", r.URL.Path)
			}
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_delete_goal", map[string]any{"goal_id": "goal1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})
}

func TestClickupCreateKeyResult(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_create_key_result", map[string]any{"goal_id": "goal1", "name": "KR1"})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing required type")
		}
	})

	t.Run("wiring", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("method = %s, want POST", r.Method)
			}
			if r.URL.Path != "/goal/goal1/key_result" {
				t.Errorf("path = %s, want /goal/goal1/key_result", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["name"] != "KR1" {
				t.Errorf("name = %v, want KR1", body["name"])
			}
			if body["type"] != "number" {
				t.Errorf("type = %v, want number", body["type"])
			}
			if body["steps_start"] != float64(0) {
				t.Errorf("steps_start = %v, want 0", body["steps_start"])
			}
			if body["steps_end"] != float64(100) {
				t.Errorf("steps_end = %v, want 100", body["steps_end"])
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"kr1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_create_key_result", map[string]any{
			"goal_id":     "goal1",
			"name":        "KR1",
			"type":        "number",
			"steps_start": 0,
			"steps_end":   100,
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if !strings.Contains(textOf(t, res), "kr1") {
			t.Errorf("result = %q, want it to contain kr1", textOf(t, res))
		}
	})

	t.Run("only_required_fields_when_optional_omitted", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if len(body) != 2 {
				t.Errorf("body = %v, want only name and type set", body)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"kr1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_create_key_result", map[string]any{
			"goal_id": "goal1",
			"name":    "KR1",
			"type":    "number",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})
}

func TestClickupUpdateKeyResult(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_update_key_result", map[string]any{"note": "x"})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing key_result_id")
		}
	})

	t.Run("partial_update_only_steps_current", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("method = %s, want PUT", r.Method)
			}
			if r.URL.Path != "/key_result/kr1" {
				t.Errorf("path = %s, want /key_result/kr1", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if len(body) != 1 {
				t.Errorf("body = %v, want only steps_current set", body)
			}
			if body["steps_current"] != float64(42) {
				t.Errorf("steps_current = %v, want 42", body["steps_current"])
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"kr1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_update_key_result", map[string]any{"key_result_id": "kr1", "steps_current": 42})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})

	t.Run("partial_update_only_note", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if len(body) != 1 {
				t.Errorf("body = %v, want only note set", body)
			}
			if body["note"] != "progress note" {
				t.Errorf("note = %v, want progress note", body["note"])
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"kr1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_update_key_result", map[string]any{"key_result_id": "kr1", "note": "progress note"})
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
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_update_key_result", map[string]any{"key_result_id": "kr1", "note": "x"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("error text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupDeleteKeyResult(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_delete_key_result", map[string]any{})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing key_result_id")
		}
	})

	t.Run("wiring", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("method = %s, want DELETE", r.Method)
			}
			if r.URL.Path != "/key_result/kr1" {
				t.Errorf("path = %s, want /key_result/kr1", r.URL.Path)
			}
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGoalTools(s, c)
		res := callTool(t, s, "clickup_delete_key_result", map[string]any{"key_result_id": "kr1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})
}
