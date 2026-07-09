package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestAddTaskDependencyPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{}`))
	})

	err := c.AddTaskDependency(context.Background(), "task1", map[string]any{"depends_on": "task2"})
	if err != nil {
		t.Fatalf("AddTaskDependency: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/task/task1/dependency" {
		t.Errorf("path = %q, want /task/task1/dependency", gotPath)
	}
	if gotBody["depends_on"] != "task2" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestRemoveTaskDependencyBuildsQuery(t *testing.T) {
	var gotMethod, gotPath, gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{}`))
	})

	err := c.RemoveTaskDependency(context.Background(), "task1", "task2", "task3")
	if err != nil {
		t.Fatalf("RemoveTaskDependency: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/task/task1/dependency" {
		t.Errorf("path = %q, want /task/task1/dependency", gotPath)
	}
	if gotQuery != "dependency_of=task3&depends_on=task2" {
		t.Errorf("query = %q", gotQuery)
	}
}

func TestRemoveTaskDependencyOmitsEmptyParams(t *testing.T) {
	var gotQuery url.Values
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		_, _ = w.Write([]byte(`{}`))
	})

	if err := c.RemoveTaskDependency(context.Background(), "task1", "task2", ""); err != nil {
		t.Fatalf("RemoveTaskDependency: %v", err)
	}
	if gotQuery.Get("depends_on") != "task2" {
		t.Errorf("depends_on = %q, want task2", gotQuery.Get("depends_on"))
	}
	if _, ok := gotQuery["dependency_of"]; ok {
		t.Errorf("dependency_of present = %v, want absent", gotQuery["dependency_of"])
	}
}

func TestAddTaskDependencyAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"err":"Cannot add dependency","ECODE":"DEP_001"}`))
	})

	err := c.AddTaskDependency(context.Background(), "task1", map[string]any{"depends_on": "task2"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest || apiErr.ECode != "DEP_001" || apiErr.Err != "Cannot add dependency" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
