package tools

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

// neverCalled returns an http.HandlerFunc that fails the test if it is ever
// invoked, for asserting that required-arg validation short-circuits before
// any HTTP call is made.
func neverCalled(t *testing.T) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected HTTP call: %s %s", r.Method, r.URL.Path)
	}
}

func decodeJSONBody(t *testing.T, r *http.Request) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		t.Fatalf("decoding request body: %v", err)
	}
	return body
}

func TestClickupListTimeEntries(t *testing.T) {
	t.Run("default team and no filters", func(t *testing.T) {
		var gotPath, gotMethod string
		var gotQuery url.Values
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotQuery = r.URL.Query()
			_, _ = w.Write([]byte(`{"data":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_list_time_entries", map[string]any{})
		if res.IsError {
			t.Fatalf("unexpected error: %s", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %s, want GET", gotMethod)
		}
		if gotPath != "/team/999/time_entries" {
			t.Errorf("path = %q, want /team/999/time_entries", gotPath)
		}
		if len(gotQuery) != 0 {
			t.Errorf("query = %v, want empty", gotQuery)
		}
	})

	t.Run("filters wiring", func(t *testing.T) {
		var gotPath string
		var gotQuery url.Values
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotQuery = r.URL.Query()
			_, _ = w.Write([]byte(`{"data":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_list_time_entries", map[string]any{
			"team_id":    "888",
			"start_date": "1970-01-01 00:00:01",
			"end_date":   "1970-01-01 00:00:02",
			"assignee":   "u1",
			"space_id":   "sp1",
			"folder_id":  "f1",
			"list_id":    "l1",
			"task_id":    "t1",
		})
		if res.IsError {
			t.Fatalf("unexpected error: %s", textOf(t, res))
		}
		if gotPath != "/team/888/time_entries" {
			t.Errorf("path = %q, want /team/888/time_entries", gotPath)
		}
		want := map[string]string{
			"start_date": "1000",
			"end_date":   "2000",
			"assignee":   "u1",
			"space_id":   "sp1",
			"folder_id":  "f1",
			"list_id":    "l1",
			"task_id":    "t1",
		}
		for k, v := range want {
			if got := gotQuery.Get(k); got != v {
				t.Errorf("query[%s] = %q, want %q", k, got, v)
			}
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"TIME_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_list_time_entries", map[string]any{})
		if !res.IsError {
			t.Fatal("expected IsError = true")
		}
		text := textOf(t, res)
		want := "ClickUp API error 404 [TIME_001]: not found"
		if text != want {
			t.Errorf("error text = %q, want %q", text, want)
		}
	})
}

func TestClickupCreateTimeEntry(t *testing.T) {
	t.Run("missing start", func(t *testing.T) {
		c, _ := newTestClient(t, neverCalled(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_create_time_entry", map[string]any{"duration": float64(1000)})
		if !res.IsError {
			t.Fatal("expected IsError = true for missing start")
		}
	})

	t.Run("missing duration", func(t *testing.T) {
		c, _ := newTestClient(t, neverCalled(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_create_time_entry", map[string]any{"start": "1970-01-01 00:00:01"})
		if !res.IsError {
			t.Fatal("expected IsError = true for missing duration")
		}
	})

	t.Run("full wiring", func(t *testing.T) {
		var gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotBody = decodeJSONBody(t, r)
			_, _ = w.Write([]byte(`{"data":{"id":"e1"}}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_create_time_entry", map[string]any{
			"team_id":     "777",
			"task_id":     "t1",
			"start":       "1970-01-01 00:00:01",
			"duration":    float64(5000),
			"description": "worked",
			"billable":    true,
			"assignee":    "u1",
			"tags":        []any{"a", "b"},
		})
		if res.IsError {
			t.Fatalf("unexpected error: %s", textOf(t, res))
		}
		if gotPath != "/team/777/time_entries" {
			t.Errorf("path = %q, want /team/777/time_entries", gotPath)
		}
		if gotBody["tid"] != "t1" {
			t.Errorf("body[tid] = %v, want t1", gotBody["tid"])
		}
		if gotBody["description"] != "worked" {
			t.Errorf("body[description] = %v", gotBody["description"])
		}
		if gotBody["billable"] != true {
			t.Errorf("body[billable] = %v", gotBody["billable"])
		}
		if gotBody["assignee"] != "u1" {
			t.Errorf("body[assignee] = %v", gotBody["assignee"])
		}
		if v, ok := gotBody["start"].(float64); !ok || v != 1000 {
			t.Errorf("body[start] = %v", gotBody["start"])
		}
		if v, ok := gotBody["duration"].(float64); !ok || v != 5000 {
			t.Errorf("body[duration] = %v", gotBody["duration"])
		}
	})

	t.Run("minimal wiring omits optional fields", func(t *testing.T) {
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotBody = decodeJSONBody(t, r)
			_, _ = w.Write([]byte(`{"data":{"id":"e1"}}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_create_time_entry", map[string]any{
			"start":    "1970-01-01 00:00:01",
			"duration": float64(5000),
		})
		if res.IsError {
			t.Fatalf("unexpected error: %s", textOf(t, res))
		}
		for _, k := range []string{"tid", "description", "billable", "assignee", "tags"} {
			if _, ok := gotBody[k]; ok {
				t.Errorf("body[%s] = %v, want absent", k, gotBody[k])
			}
		}
	})
}

func TestClickupGetTimeEntry(t *testing.T) {
	t.Run("missing timer_id", func(t *testing.T) {
		c, _ := newTestClient(t, neverCalled(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_get_time_entry", map[string]any{})
		if !res.IsError {
			t.Fatal("expected IsError = true for missing timer_id")
		}
	})

	t.Run("success", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"data":{"id":"e1"}}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_get_time_entry", map[string]any{"timer_id": "e1"})
		if res.IsError {
			t.Fatalf("unexpected error: %s", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %s, want GET", gotMethod)
		}
		if gotPath != "/team/999/time_entries/e1" {
			t.Errorf("path = %q, want /team/999/time_entries/e1", gotPath)
		}
	})
}

func TestClickupUpdateTimeEntry(t *testing.T) {
	t.Run("missing timer_id", func(t *testing.T) {
		c, _ := newTestClient(t, neverCalled(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_update_time_entry", map[string]any{"description": "x"})
		if !res.IsError {
			t.Fatal("expected IsError = true for missing timer_id")
		}
	})

	t.Run("partial update sends only supplied fields", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody = decodeJSONBody(t, r)
			_, _ = w.Write([]byte(`{"data":{"id":"e1"}}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_update_time_entry", map[string]any{
			"timer_id":    "e1",
			"description": "updated",
		})
		if res.IsError {
			t.Fatalf("unexpected error: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPut {
			t.Errorf("method = %s, want PUT", gotMethod)
		}
		if gotPath != "/team/999/time_entries/e1" {
			t.Errorf("path = %q, want /team/999/time_entries/e1", gotPath)
		}
		if len(gotBody) != 1 || gotBody["description"] != "updated" {
			t.Errorf("body = %+v, want only description", gotBody)
		}
	})

	t.Run("full update sends all supplied fields", func(t *testing.T) {
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotBody = decodeJSONBody(t, r)
			_, _ = w.Write([]byte(`{"data":{"id":"e1"}}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_update_time_entry", map[string]any{
			"timer_id":    "e1",
			"start":       "1970-01-01 00:00:00",
			"duration":    float64(20),
			"description": "d",
			"billable":    false,
		})
		if res.IsError {
			t.Fatalf("unexpected error: %s", textOf(t, res))
		}
		if len(gotBody) != 4 {
			t.Errorf("body = %+v, want 4 keys", gotBody)
		}
	})
}

func TestClickupDeleteTimeEntry(t *testing.T) {
	t.Run("missing timer_id", func(t *testing.T) {
		c, _ := newTestClient(t, neverCalled(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_delete_time_entry", map[string]any{})
		if !res.IsError {
			t.Fatal("expected IsError = true for missing timer_id")
		}
	})

	t.Run("success", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_delete_time_entry", map[string]any{"timer_id": "e1"})
		if res.IsError {
			t.Fatalf("unexpected error: %s", textOf(t, res))
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", gotMethod)
		}
		if gotPath != "/team/999/time_entries/e1" {
			t.Errorf("path = %q, want /team/999/time_entries/e1", gotPath)
		}
		var out map[string]any
		if err := json.Unmarshal([]byte(textOf(t, res)), &out); err != nil {
			t.Fatalf("decoding result: %v", err)
		}
		if out["deleted"] != true || out["timer_id"] != "e1" {
			t.Errorf("result = %+v", out)
		}
	})
}

func TestClickupGetTimeEntryHistory(t *testing.T) {
	t.Run("missing timer_id", func(t *testing.T) {
		c, _ := newTestClient(t, neverCalled(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_get_time_entry_history", map[string]any{})
		if !res.IsError {
			t.Fatal("expected IsError = true for missing timer_id")
		}
	})

	t.Run("success", func(t *testing.T) {
		var gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"data":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_get_time_entry_history", map[string]any{"timer_id": "e1"})
		if res.IsError {
			t.Fatalf("unexpected error: %s", textOf(t, res))
		}
		if gotPath != "/team/999/time_entries/e1/history" {
			t.Errorf("path = %q, want /team/999/time_entries/e1/history", gotPath)
		}
	})
}

func TestClickupGetCurrentTimeEntry(t *testing.T) {
	t.Run("no assignee", func(t *testing.T) {
		var gotPath string
		var gotQuery url.Values
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotQuery = r.URL.Query()
			_, _ = w.Write([]byte(`{"data":null}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_get_current_time_entry", map[string]any{})
		if res.IsError {
			t.Fatalf("unexpected error: %s", textOf(t, res))
		}
		if gotPath != "/team/999/time_entries/current" {
			t.Errorf("path = %q, want /team/999/time_entries/current", gotPath)
		}
		if gotQuery.Get("assignee") != "" {
			t.Errorf("assignee query = %q, want empty", gotQuery.Get("assignee"))
		}
	})

	t.Run("with assignee", func(t *testing.T) {
		var gotQuery url.Values
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotQuery = r.URL.Query()
			_, _ = w.Write([]byte(`{"data":null}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_get_current_time_entry", map[string]any{"assignee": "u9"})
		if res.IsError {
			t.Fatalf("unexpected error: %s", textOf(t, res))
		}
		if gotQuery.Get("assignee") != "u9" {
			t.Errorf("assignee query = %q, want u9", gotQuery.Get("assignee"))
		}
	})
}

func TestClickupStartTimeEntry(t *testing.T) {
	t.Run("missing task_id", func(t *testing.T) {
		c, _ := newTestClient(t, neverCalled(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_start_time_entry", map[string]any{})
		if !res.IsError {
			t.Fatal("expected IsError = true for missing task_id")
		}
	})

	t.Run("success", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody = decodeJSONBody(t, r)
			_, _ = w.Write([]byte(`{"data":{"id":"e1"}}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_start_time_entry", map[string]any{
			"task_id":     "t1",
			"description": "starting",
			"tags":        []any{"tag1"},
		})
		if res.IsError {
			t.Fatalf("unexpected error: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %s, want POST", gotMethod)
		}
		if gotPath != "/team/999/time_entries/start/t1" {
			t.Errorf("path = %q, want /team/999/time_entries/start/t1", gotPath)
		}
		if gotBody["description"] != "starting" {
			t.Errorf("body[description] = %v", gotBody["description"])
		}
	})
}

func TestClickupStopTimeEntry(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":{"id":"e1"}}`))
	})
	s := server.NewMCPServer("test", "1.0.0")
	RegisterTimeTrackingTools(s, c)
	res := callTool(t, s, "clickup_stop_time_entry", map[string]any{})
	if res.IsError {
		t.Fatalf("unexpected error: %s", textOf(t, res))
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/team/999/time_entries/stop" {
		t.Errorf("path = %q, want /team/999/time_entries/stop", gotPath)
	}
}

func TestClickupListTimeEntryTags(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":[]}`))
	})
	s := server.NewMCPServer("test", "1.0.0")
	RegisterTimeTrackingTools(s, c)
	res := callTool(t, s, "clickup_list_time_entry_tags", map[string]any{})
	if res.IsError {
		t.Fatalf("unexpected error: %s", textOf(t, res))
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/team/999/time_entries/tags" {
		t.Errorf("path = %q, want /team/999/time_entries/tags", gotPath)
	}
}

func TestClickupAddTimeEntryTags(t *testing.T) {
	t.Run("missing time_entry_ids", func(t *testing.T) {
		c, _ := newTestClient(t, neverCalled(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_add_time_entry_tags", map[string]any{"tags": []any{"a"}})
		if !res.IsError {
			t.Fatal("expected IsError = true for missing time_entry_ids")
		}
	})

	t.Run("missing tags", func(t *testing.T) {
		c, _ := newTestClient(t, neverCalled(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_add_time_entry_tags", map[string]any{"time_entry_ids": []any{"e1"}})
		if !res.IsError {
			t.Fatal("expected IsError = true for missing tags")
		}
	})

	t.Run("success wraps tag names as objects", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody = decodeJSONBody(t, r)
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_add_time_entry_tags", map[string]any{
			"time_entry_ids": []any{"e1", "e2"},
			"tags":           []any{"urgent", "billed"},
		})
		if res.IsError {
			t.Fatalf("unexpected error: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %s, want POST", gotMethod)
		}
		if gotPath != "/team/999/time_entries/tags" {
			t.Errorf("path = %q, want /team/999/time_entries/tags", gotPath)
		}
		tags, ok := gotBody["tags"].([]any)
		if !ok || len(tags) != 2 {
			t.Fatalf("body[tags] = %v", gotBody["tags"])
		}
		first, ok := tags[0].(map[string]any)
		if !ok || first["name"] != "urgent" {
			t.Errorf("tags[0] = %v, want {name: urgent}", tags[0])
		}
	})
}

func TestClickupRenameTimeEntryTag(t *testing.T) {
	t.Run("missing name", func(t *testing.T) {
		c, _ := newTestClient(t, neverCalled(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_rename_time_entry_tag", map[string]any{"new_name": "x"})
		if !res.IsError {
			t.Fatal("expected IsError = true for missing name")
		}
	})

	t.Run("missing new_name", func(t *testing.T) {
		c, _ := newTestClient(t, neverCalled(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_rename_time_entry_tag", map[string]any{"name": "x"})
		if !res.IsError {
			t.Fatal("expected IsError = true for missing new_name")
		}
	})

	t.Run("success", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody = decodeJSONBody(t, r)
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_rename_time_entry_tag", map[string]any{
			"name":     "old",
			"new_name": "new",
			"tag_bg":   "#fff",
			"tag_fg":   "#000",
		})
		if res.IsError {
			t.Fatalf("unexpected error: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPut {
			t.Errorf("method = %s, want PUT", gotMethod)
		}
		if gotPath != "/team/999/time_entries/tags" {
			t.Errorf("path = %q, want /team/999/time_entries/tags", gotPath)
		}
		if gotBody["name"] != "old" || gotBody["new_name"] != "new" || gotBody["tag_bg"] != "#fff" || gotBody["tag_fg"] != "#000" {
			t.Errorf("body = %+v", gotBody)
		}
	})
}

func TestClickupRemoveTimeEntryTags(t *testing.T) {
	t.Run("missing time_entry_ids", func(t *testing.T) {
		c, _ := newTestClient(t, neverCalled(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_remove_time_entry_tags", map[string]any{"tags": []any{"a"}})
		if !res.IsError {
			t.Fatal("expected IsError = true for missing time_entry_ids")
		}
	})

	t.Run("success", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody = decodeJSONBody(t, r)
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterTimeTrackingTools(s, c)
		res := callTool(t, s, "clickup_remove_time_entry_tags", map[string]any{
			"time_entry_ids": []any{"e1"},
			"tags":           []any{"urgent"},
		})
		if res.IsError {
			t.Fatalf("unexpected error: %s", textOf(t, res))
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", gotMethod)
		}
		if gotPath != "/team/999/time_entries/tags" {
			t.Errorf("path = %q, want /team/999/time_entries/tags", gotPath)
		}
		ids, ok := gotBody["time_entry_ids"].([]any)
		if !ok || len(ids) != 1 || ids[0] != "e1" {
			t.Errorf("body[time_entry_ids] = %v", gotBody["time_entry_ids"])
		}
	})
}
