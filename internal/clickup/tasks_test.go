package clickup

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetTaskHitsExpectedPath(t *testing.T) {
	var gotPath, gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"id":"abc123"}`))
	})

	if _, err := c.GetTask(context.Background(), "abc123", true, "999"); err != nil {
		t.Fatalf("GetTask: %v", err)
	}
	if gotPath != "/task/abc123" {
		t.Errorf("path = %q, want /task/abc123", gotPath)
	}
	if gotQuery != "custom_task_ids=true&team_id=999" {
		t.Errorf("query = %q", gotQuery)
	}
}

func TestCreateTaskPostsExpectedBody(t *testing.T) {
	var gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"new"}`))
	})

	_, err := c.CreateTask(context.Background(), "list1", map[string]any{"name": "Buy milk", "priority": float64(2)})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	if gotPath != "/list/list1/task" {
		t.Errorf("path = %q, want /list/list1/task", gotPath)
	}
	if gotBody["name"] != "Buy milk" || gotBody["priority"] != float64(2) {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestSearchTasksBuildsArrayQueryParams(t *testing.T) {
	var gotQuery map[string][]string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = map[string][]string(r.URL.Query())
		_, _ = w.Write([]byte(`{"tasks":[]}`))
	})

	_, err := c.SearchTasks(context.Background(), "999", SearchTasksOptions{
		SpaceIDs: []string{"s1", "s2"},
		TaskQueryFilters: TaskQueryFilters{
			Statuses: []string{"open", "in progress"},
		},
	})
	if err != nil {
		t.Fatalf("SearchTasks: %v", err)
	}
	if got := gotQuery["space_ids[]"]; len(got) != 2 || got[0] != "s1" || got[1] != "s2" {
		t.Errorf("space_ids[] = %v", got)
	}
	if got := gotQuery["statuses[]"]; len(got) != 2 {
		t.Errorf("statuses[] = %v", got)
	}
}
