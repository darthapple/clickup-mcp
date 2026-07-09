package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func int64Ptr(v int64) *int64 { return &v }
func boolPtr(v bool) *bool    { return &v }

func TestListTimeEntriesBuildsFullQuery(t *testing.T) {
	var gotMethod, gotPath, gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"data":[]}`))
	})

	filters := TimeEntryFilters{
		StartDate:  int64Ptr(1000),
		EndDate:    int64Ptr(2000),
		Assignee:   "42",
		SpaceID:    "s1",
		FolderID:   "f1",
		ListID:     "l1",
		TaskID:     "t1",
		IncludeAll: boolPtr(true),
	}

	out, err := c.ListTimeEntries(context.Background(), "team1", filters)
	if err != nil {
		t.Fatalf("ListTimeEntries: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/team/team1/time_entries" {
		t.Errorf("path = %q, want /team/team1/time_entries", gotPath)
	}
	want := "assignee=42&end_date=2000&folder_id=f1&include_task_tags=true&list_id=l1&space_id=s1&start_date=1000&task_id=t1"
	if gotQuery != want {
		t.Errorf("query = %q, want %q", gotQuery, want)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded response is not a map: %T", out)
	}
	if _, ok := m["data"]; !ok {
		t.Errorf("decoded response missing data: %+v", m)
	}
}

func TestListTimeEntriesOmitsUnsetFilters(t *testing.T) {
	var gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"data":[]}`))
	})

	if _, err := c.ListTimeEntries(context.Background(), "team1", TimeEntryFilters{}); err != nil {
		t.Fatalf("ListTimeEntries: %v", err)
	}
	if gotQuery != "" {
		t.Errorf("query = %q, want empty", gotQuery)
	}
}

func TestCreateTimeEntryPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"data":{"id":"te1"}}`))
	})

	out, err := c.CreateTimeEntry(context.Background(), "team1", map[string]any{"duration": float64(3600000)})
	if err != nil {
		t.Fatalf("CreateTimeEntry: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/team/team1/time_entries" {
		t.Errorf("path = %q, want /team/team1/time_entries", gotPath)
	}
	if gotBody["duration"] != float64(3600000) {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded response is not a map: %T", out)
	}
	if _, ok := m["data"]; !ok {
		t.Errorf("decoded response missing data: %+v", m)
	}
}

func TestGetTimeEntryHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":{"id":"te1"}}`))
	})

	if _, err := c.GetTimeEntry(context.Background(), "team1", "te1"); err != nil {
		t.Fatalf("GetTimeEntry: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/team/team1/time_entries/te1" {
		t.Errorf("path = %q, want /team/team1/time_entries/te1", gotPath)
	}
}

func TestUpdateTimeEntryPutsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"data":{"id":"te1"}}`))
	})

	if _, err := c.UpdateTimeEntry(context.Background(), "team1", "te1", map[string]any{"description": "updated"}); err != nil {
		t.Fatalf("UpdateTimeEntry: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	if gotPath != "/team/team1/time_entries/te1" {
		t.Errorf("path = %q, want /team/team1/time_entries/te1", gotPath)
	}
	if gotBody["description"] != "updated" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestDeleteTimeEntryHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{}`))
	})

	if err := c.DeleteTimeEntry(context.Background(), "team1", "te1"); err != nil {
		t.Fatalf("DeleteTimeEntry: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/team/team1/time_entries/te1" {
		t.Errorf("path = %q, want /team/team1/time_entries/te1", gotPath)
	}
}

func TestGetTimeEntryHistoryHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":[]}`))
	})

	if _, err := c.GetTimeEntryHistory(context.Background(), "team1", "te1"); err != nil {
		t.Fatalf("GetTimeEntryHistory: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/team/team1/time_entries/te1/history" {
		t.Errorf("path = %q, want /team/team1/time_entries/te1/history", gotPath)
	}
}

func TestGetCurrentTimeEntryBuildsQuery(t *testing.T) {
	var gotMethod, gotPath, gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"data":null}`))
	})

	if _, err := c.GetCurrentTimeEntry(context.Background(), "team1", "42"); err != nil {
		t.Fatalf("GetCurrentTimeEntry: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/team/team1/time_entries/current" {
		t.Errorf("path = %q, want /team/team1/time_entries/current", gotPath)
	}
	if gotQuery != "assignee=42" {
		t.Errorf("query = %q, want assignee=42", gotQuery)
	}
}

func TestGetCurrentTimeEntryOmitsAssigneeWhenEmpty(t *testing.T) {
	var gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"data":null}`))
	})

	if _, err := c.GetCurrentTimeEntry(context.Background(), "team1", ""); err != nil {
		t.Fatalf("GetCurrentTimeEntry: %v", err)
	}
	if gotQuery != "" {
		t.Errorf("query = %q, want empty", gotQuery)
	}
}

func TestStartTimeEntryPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"data":{"id":"te1"}}`))
	})

	if _, err := c.StartTimeEntry(context.Background(), "team1", "task1", map[string]any{"description": "working"}); err != nil {
		t.Fatalf("StartTimeEntry: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/team/team1/time_entries/start/task1" {
		t.Errorf("path = %q, want /team/team1/time_entries/start/task1", gotPath)
	}
	if gotBody["description"] != "working" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestStopTimeEntryHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":{"id":"te1"}}`))
	})

	if _, err := c.StopTimeEntry(context.Background(), "team1"); err != nil {
		t.Fatalf("StopTimeEntry: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/team/team1/time_entries/stop" {
		t.Errorf("path = %q, want /team/team1/time_entries/stop", gotPath)
	}
}

func TestListTimeEntryTagsHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":[]}`))
	})

	if _, err := c.ListTimeEntryTags(context.Background(), "team1"); err != nil {
		t.Fatalf("ListTimeEntryTags: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/team/team1/time_entries/tags" {
		t.Errorf("path = %q, want /team/team1/time_entries/tags", gotPath)
	}
}

func TestAddTimeEntryTagsPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{}`))
	})

	body := map[string]any{"time_entry_ids": []any{"te1", "te2"}, "tags": []any{"billable"}}
	if err := c.AddTimeEntryTags(context.Background(), "team1", body); err != nil {
		t.Fatalf("AddTimeEntryTags: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/team/team1/time_entries/tags" {
		t.Errorf("path = %q, want /team/team1/time_entries/tags", gotPath)
	}
	ids, ok := gotBody["time_entry_ids"].([]any)
	if !ok || len(ids) != 2 {
		t.Errorf("time_entry_ids = %+v", gotBody["time_entry_ids"])
	}
}

func TestRenameTimeEntryTagPutsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{}`))
	})

	body := map[string]any{"name": "billable", "new_name": "Billable"}
	if err := c.RenameTimeEntryTag(context.Background(), "team1", body); err != nil {
		t.Fatalf("RenameTimeEntryTag: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	if gotPath != "/team/team1/time_entries/tags" {
		t.Errorf("path = %q, want /team/team1/time_entries/tags", gotPath)
	}
	if gotBody["new_name"] != "Billable" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestRemoveTimeEntryTagsDeletesExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{}`))
	})

	body := map[string]any{"time_entry_ids": []any{"te1"}, "tags": []any{"billable"}}
	if err := c.RemoveTimeEntryTags(context.Background(), "team1", body); err != nil {
		t.Fatalf("RemoveTimeEntryTags: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/team/team1/time_entries/tags" {
		t.Errorf("path = %q, want /team/team1/time_entries/tags", gotPath)
	}
	tags, ok := gotBody["tags"].([]any)
	if !ok || len(tags) != 1 {
		t.Errorf("tags = %+v", gotBody["tags"])
	}
}

func TestListTimeEntriesAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"err":"Not permitted","ECODE":"TIME_001"}`))
	})

	_, err := c.ListTimeEntries(context.Background(), "team1", TimeEntryFilters{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusForbidden || apiErr.ECode != "TIME_001" || apiErr.Err != "Not permitted" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
