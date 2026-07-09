package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestListTeamViewsHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"views":[{"id":"v1"}]}`))
	})

	out, err := c.ListTeamViews(context.Background(), "team1")
	if err != nil {
		t.Fatalf("ListTeamViews: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/team/team1/view" {
		t.Errorf("path = %q, want /team/team1/view", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded response is not a map: %T", out)
	}
	views, ok := m["views"].([]any)
	if !ok || len(views) != 1 {
		t.Errorf("views = %+v", m["views"])
	}
}

func TestListSpaceViewsHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"views":[]}`))
	})

	if _, err := c.ListSpaceViews(context.Background(), "space1"); err != nil {
		t.Fatalf("ListSpaceViews: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/space/space1/view" {
		t.Errorf("path = %q, want /space/space1/view", gotPath)
	}
}

func TestCreateSpaceViewPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"view":{"id":"v1"}}`))
	})

	out, err := c.CreateSpaceView(context.Background(), "space1", map[string]any{"name": "Board"})
	if err != nil {
		t.Fatalf("CreateSpaceView: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/space/space1/view" {
		t.Errorf("path = %q, want /space/space1/view", gotPath)
	}
	if gotBody["name"] != "Board" {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded response is not a map: %T", out)
	}
	if _, ok := m["view"]; !ok {
		t.Errorf("decoded response missing view: %+v", m)
	}
}

func TestListFolderViewsHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"views":[]}`))
	})

	if _, err := c.ListFolderViews(context.Background(), "folder1"); err != nil {
		t.Fatalf("ListFolderViews: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/folder/folder1/view" {
		t.Errorf("path = %q, want /folder/folder1/view", gotPath)
	}
}

func TestCreateFolderViewPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"view":{"id":"v1"}}`))
	})

	if _, err := c.CreateFolderView(context.Background(), "folder1", map[string]any{"name": "List view"}); err != nil {
		t.Fatalf("CreateFolderView: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/folder/folder1/view" {
		t.Errorf("path = %q, want /folder/folder1/view", gotPath)
	}
	if gotBody["name"] != "List view" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestListListViewsHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"views":[]}`))
	})

	if _, err := c.ListListViews(context.Background(), "list1"); err != nil {
		t.Fatalf("ListListViews: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/list/list1/view" {
		t.Errorf("path = %q, want /list/list1/view", gotPath)
	}
}

func TestCreateListViewPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"view":{"id":"v1"}}`))
	})

	if _, err := c.CreateListView(context.Background(), "list1", map[string]any{"name": "Calendar"}); err != nil {
		t.Fatalf("CreateListView: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/list/list1/view" {
		t.Errorf("path = %q, want /list/list1/view", gotPath)
	}
	if gotBody["name"] != "Calendar" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestGetViewHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"id":"v1"}`))
	})

	out, err := c.GetView(context.Background(), "v1")
	if err != nil {
		t.Fatalf("GetView: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/view/v1" {
		t.Errorf("path = %q, want /view/v1", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok || m["id"] != "v1" {
		t.Errorf("decoded response = %+v", out)
	}
}

func TestUpdateViewPutsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"v1"}`))
	})

	if _, err := c.UpdateView(context.Background(), "v1", map[string]any{"name": "Renamed"}); err != nil {
		t.Fatalf("UpdateView: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	if gotPath != "/view/v1" {
		t.Errorf("path = %q, want /view/v1", gotPath)
	}
	if gotBody["name"] != "Renamed" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestDeleteViewHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{}`))
	})

	if err := c.DeleteView(context.Background(), "v1"); err != nil {
		t.Fatalf("DeleteView: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/view/v1" {
		t.Errorf("path = %q, want /view/v1", gotPath)
	}
}

func TestGetViewTasksBuildsQueryWhenPageSet(t *testing.T) {
	var gotMethod, gotPath, gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"tasks":[]}`))
	})

	out, err := c.GetViewTasks(context.Background(), "v1", 2, true)
	if err != nil {
		t.Fatalf("GetViewTasks: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/view/v1/task" {
		t.Errorf("path = %q, want /view/v1/task", gotPath)
	}
	if gotQuery != "page=2" {
		t.Errorf("query = %q, want page=2", gotQuery)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded response is not a map: %T", out)
	}
	if _, ok := m["tasks"]; !ok {
		t.Errorf("decoded response missing tasks: %+v", m)
	}
}

func TestGetViewTasksOmitsPageWhenNotSet(t *testing.T) {
	var gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"tasks":[]}`))
	})

	if _, err := c.GetViewTasks(context.Background(), "v1", 0, false); err != nil {
		t.Fatalf("GetViewTasks: %v", err)
	}
	if gotQuery != "" {
		t.Errorf("query = %q, want empty", gotQuery)
	}
}

func TestGetViewAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"err":"View not found","ECODE":"VIEW_001"}`))
	})

	_, err := c.GetView(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusNotFound || apiErr.ECode != "VIEW_001" || apiErr.Err != "View not found" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
