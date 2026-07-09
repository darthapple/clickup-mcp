package clickup

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestListTaskMembersHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"members":[{"id":1,"username":"alice"}]}`))
	})

	out, err := c.ListTaskMembers(context.Background(), "task1")
	if err != nil {
		t.Fatalf("ListTaskMembers: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/task/task1/member" {
		t.Errorf("path = %q, want /task/task1/member", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok || m["members"] == nil {
		t.Errorf("decoded = %+v", out)
	}
}

func TestListListMembersHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"members":[{"id":2,"username":"bob"}]}`))
	})

	out, err := c.ListListMembers(context.Background(), "list1")
	if err != nil {
		t.Fatalf("ListListMembers: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/list/list1/member" {
		t.Errorf("path = %q, want /list/list1/member", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok || m["members"] == nil {
		t.Errorf("decoded = %+v", out)
	}
}

func TestMembersAPIErrorReturnsAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"err":"Task not found","ECODE":"TASK_002"}`))
	})

	_, err := c.ListTaskMembers(context.Background(), "task1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusNotFound || apiErr.ECode != "TASK_002" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
