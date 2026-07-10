package tools

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestFormatDuration(t *testing.T) {
	cases := []struct {
		ms   int64
		want string
	}{
		{0, "00:00:00:00"},
		{1000, "00:00:00:01"},
		{61000, "00:00:01:01"},
		{3661000, "00:01:01:01"},
		// 1 day (86400000ms) + 1hr1min1sec (3661000ms) = 90061000ms, spanning
		// multiple days.
		{90061000, "01:01:01:01"},
		{-5000, "00:00:00:00"}, // negative clamps to zero
	}
	for _, tc := range cases {
		if got := formatDuration(tc.ms); got != tc.want {
			t.Errorf("formatDuration(%d) = %q, want %q", tc.ms, got, tc.want)
		}
	}
}

func TestParseMsField(t *testing.T) {
	cases := []struct {
		name string
		m    map[string]any
		key  string
		want int64
	}{
		{"missing key", map[string]any{}, "duration", 0},
		{"empty string", map[string]any{"duration": ""}, "duration", 0},
		{"non-numeric string", map[string]any{"duration": "not-a-number"}, "duration", 0},
		{"zero", map[string]any{"duration": "0"}, "duration", 0},
		{"ordinary value", map[string]any{"duration": "12345"}, "duration", 12345},
		{"large epoch-scale value", map[string]any{"start": "1700000000000"}, "start", 1700000000000},
		{"wrong type is not a string", map[string]any{"duration": float64(500)}, "duration", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := parseMsField(tc.m, tc.key); got != tc.want {
				t.Errorf("parseMsField(%v, %q) = %d, want %d", tc.m, tc.key, got, tc.want)
			}
		})
	}
}

func TestClickupGetListTimeReport(t *testing.T) {
	t.Run("requires list_id, start_date, end_date", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterReportTools(s, c)

		res := callTool(t, s, "clickup_get_list_time_report", map[string]any{
			"start_date": float64(0),
			"end_date":   float64(1000),
		})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing list_id)")
		}

		res = callTool(t, s, "clickup_get_list_time_report", map[string]any{
			"list_id":  "list1",
			"end_date": float64(1000),
		})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing start_date)")
		}

		res = callTool(t, s, "clickup_get_list_time_report", map[string]any{
			"list_id":    "list1",
			"start_date": float64(0),
		})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing end_date)")
		}
	})

	t.Run("aggregates tasks and time entries, surfacing unmatched entries", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/list/list1/task":
				_, _ = w.Write([]byte(`{
					"tasks": [
						{"id": "task1", "name": "Task One", "status": {"status": "open"}},
						{"id": "task2", "name": "Task Two", "status": {"status": "done"}}
					],
					"last_page": true
				}`))
			case "/list/list1/member":
				_, _ = w.Write([]byte(`{
					"members": [
						{"id": 1, "username": "alice"},
						{"id": 2, "username": "bob"}
					]
				}`))
			case "/team/999/time_entries":
				if got := r.URL.Query().Get("assignee"); got != "1,2" {
					t.Errorf("assignee query param = %q, want %q (all list members)", got, "1,2")
				}
				_, _ = w.Write([]byte(`{
					"data": [
						{"id": "e1", "task": {"id": "task1", "name": "Task One"}, "start": "1000", "end": "61000", "duration": "60000", "description": "work"},
						{"id": "e2", "task": {"id": "taskX", "name": "Ghost Task"}, "start": "2000", "end": "32000", "duration": "30000", "description": "orphan"}
					]
				}`))
			default:
				t.Errorf("unexpected request path: %s", r.URL.Path)
				w.WriteHeader(http.StatusNotFound)
			}
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterReportTools(s, c)

		res := callTool(t, s, "clickup_get_list_time_report", map[string]any{
			"list_id":    "list1",
			"start_date": float64(0),
			"end_date":   float64(100000),
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}

		var out map[string]any
		if err := json.Unmarshal([]byte(textOf(t, res)), &out); err != nil {
			t.Fatalf("decoding result: %v", err)
		}

		if out["total_duration_ms"] != float64(60000) {
			t.Errorf("total_duration_ms = %v, want 60000", out["total_duration_ms"])
		}

		tasks, ok := out["tasks"].([]any)
		if !ok || len(tasks) != 2 {
			t.Fatalf("tasks = %v, want 2 entries", out["tasks"])
		}
		byID := map[string]map[string]any{}
		for _, tk := range tasks {
			tm := tk.(map[string]any)
			byID[tm["task_id"].(string)] = tm
		}
		task1, ok := byID["task1"]
		if !ok {
			t.Fatalf("task1 missing from tasks: %+v", tasks)
		}
		if task1["duration_ms"] != float64(60000) {
			t.Errorf("task1 duration_ms = %v, want 60000", task1["duration_ms"])
		}
		entries1, _ := task1["entries"].([]any)
		if len(entries1) != 1 {
			t.Errorf("task1 entries = %v, want 1", task1["entries"])
		}

		task2, ok := byID["task2"]
		if !ok {
			t.Fatalf("task2 (no time entries) missing from tasks: %+v", tasks)
		}
		if task2["duration_ms"] != float64(0) {
			t.Errorf("task2 duration_ms = %v, want 0 (task with no tracked time must still appear)", task2["duration_ms"])
		}
		entries2, _ := task2["entries"].([]any)
		if len(entries2) != 0 {
			t.Errorf("task2 entries = %v, want empty", task2["entries"])
		}

		unmatched, ok := out["unmatched_entries"].([]any)
		if !ok || len(unmatched) != 1 {
			t.Fatalf("unmatched_entries = %v, want 1 entry (taskX doesn't belong to any task in the list)", out["unmatched_entries"])
		}
		um := unmatched[0].(map[string]any)
		if um["task_id"] != "taskX" || um["duration_ms"] != float64(30000) {
			t.Errorf("unmatched entry = %+v, want task_id=taskX duration_ms=30000", um)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterReportTools(s, c)

		res := callTool(t, s, "clickup_get_list_time_report", map[string]any{
			"list_id": "missing", "start_date": float64(0), "end_date": float64(1),
		})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})
}

func TestClickupGetUserTimeReport(t *testing.T) {
	t.Run("requires user_id, start_date, end_date", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterReportTools(s, c)

		res := callTool(t, s, "clickup_get_user_time_report", map[string]any{
			"start_date": float64(0), "end_date": float64(1),
		})
		if !res.IsError {
			t.Fatal("IsError = false, want true (missing user_id)")
		}
	})

	t.Run("aggregates entries by_task, sums total, and degrades gracefully on failed list lookup", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/team/999/time_entries":
				_, _ = w.Write([]byte(`{
					"data": [
						{"id": "e1", "task": {"id": "task1", "name": "Task One"}, "task_location": {"list_id": "list1", "folder_id": "folder1"}, "start": "1000", "end": "11000", "duration": "10000", "description": "a"},
						{"id": "e2", "task": {"id": "task1", "name": "Task One"}, "task_location": {"list_id": "list1", "folder_id": "folder1"}, "start": "20000", "end": "25000", "duration": "5000", "description": "b"},
						{"id": "e3", "task": {"id": "task2", "name": "Task Two"}, "task_location": {"list_id": "list2", "folder_id": "folder2"}, "start": "0", "end": "7000", "duration": "7000", "description": "c"}
					]
				}`))
			case "/list/list1":
				_, _ = w.Write([]byte(`{"name": "List One", "folder": {"name": "Folder One"}, "space": {"id": "space1", "name": "Space One"}}`))
			case "/list/list2":
				// Simulates a deleted/inaccessible list: lookup fails, report
				// must still succeed with empty names rather than erroring.
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"err":"not found","ECODE":"LIST_001"}`))
			default:
				t.Errorf("unexpected request path: %s", r.URL.Path)
				w.WriteHeader(http.StatusNotFound)
			}
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterReportTools(s, c)

		res := callTool(t, s, "clickup_get_user_time_report", map[string]any{
			"user_id": "user1", "start_date": float64(0), "end_date": float64(100000),
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}

		var out map[string]any
		if err := json.Unmarshal([]byte(textOf(t, res)), &out); err != nil {
			t.Fatalf("decoding result: %v", err)
		}

		if out["total_duration_ms"] != float64(22000) {
			t.Errorf("total_duration_ms = %v, want 22000 (10000+5000+7000)", out["total_duration_ms"])
		}

		byTask, ok := out["by_task"].([]any)
		if !ok || len(byTask) != 2 {
			t.Fatalf("by_task = %v, want 2 tasks", out["by_task"])
		}
		byTaskID := map[string]map[string]any{}
		for _, bt := range byTask {
			m := bt.(map[string]any)
			byTaskID[m["task_id"].(string)] = m
		}

		task1, ok := byTaskID["task1"]
		if !ok {
			t.Fatalf("task1 missing from by_task: %+v", byTask)
		}
		if task1["duration_ms"] != float64(15000) {
			t.Errorf("task1 aggregated duration_ms = %v, want 15000 (10000+5000)", task1["duration_ms"])
		}
		if task1["list_name"] != "List One" || task1["folder_name"] != "Folder One" || task1["space_name"] != "Space One" {
			t.Errorf("task1 location = %+v, want resolved list/folder/space names", task1)
		}

		task2, ok := byTaskID["task2"]
		if !ok {
			t.Fatalf("task2 missing from by_task: %+v", byTask)
		}
		if task2["duration_ms"] != float64(7000) {
			t.Errorf("task2 duration_ms = %v, want 7000", task2["duration_ms"])
		}
		// list2's lookup failed (404): names must degrade to empty strings,
		// not abort the whole report.
		if task2["list_name"] != "" || task2["folder_name"] != "" || task2["space_name"] != "" {
			t.Errorf("task2 location = %+v, want empty names after failed list lookup", task2)
		}

		entries, ok := out["entries"].([]any)
		if !ok || len(entries) != 3 {
			t.Fatalf("entries = %v, want 3", out["entries"])
		}
	})

	t.Run("error passthrough when the primary time-entries fetch fails", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"err":"boom","ECODE":"X_500"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterReportTools(s, c)

		res := callTool(t, s, "clickup_get_user_time_report", map[string]any{
			"user_id": "user1", "start_date": float64(0), "end_date": float64(1),
		})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})
}
