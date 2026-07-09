package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestListListsInFolderHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath, gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"lists":[{"id":"1"}]}`))
	})

	out, err := c.ListListsInFolder(context.Background(), "folder1", true, true)
	if err != nil {
		t.Fatalf("ListListsInFolder: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/folder/folder1/list" {
		t.Errorf("path = %q, want /folder/folder1/list", gotPath)
	}
	if gotQuery != "archived=true" {
		t.Errorf("query = %q, want archived=true", gotQuery)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("out is not a map: %T", out)
	}
	if _, ok := m["lists"]; !ok {
		t.Errorf("decoded = %+v, missing lists", m)
	}
}

func TestListListsInFolderOmitsArchivedWhenNotSet(t *testing.T) {
	var gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{}`))
	})

	if _, err := c.ListListsInFolder(context.Background(), "folder1", false, false); err != nil {
		t.Fatalf("ListListsInFolder: %v", err)
	}
	if gotQuery != "" {
		t.Errorf("query = %q, want empty (archived absent)", gotQuery)
	}
}

func TestListFolderlessListsHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath, gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"lists":[]}`))
	})

	if _, err := c.ListFolderlessLists(context.Background(), "space1", true, true); err != nil {
		t.Fatalf("ListFolderlessLists: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/space/space1/list" {
		t.Errorf("path = %q, want /space/space1/list", gotPath)
	}
	if gotQuery != "archived=true" {
		t.Errorf("query = %q, want archived=true", gotQuery)
	}
}

func TestListFolderlessListsOmitsArchivedWhenNotSet(t *testing.T) {
	var gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{}`))
	})

	if _, err := c.ListFolderlessLists(context.Background(), "space1", false, false); err != nil {
		t.Fatalf("ListFolderlessLists: %v", err)
	}
	if gotQuery != "" {
		t.Errorf("query = %q, want empty", gotQuery)
	}
}

func TestGetListHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"id":"list1","name":"My List"}`))
	})

	out, err := c.GetList(context.Background(), "list1")
	if err != nil {
		t.Fatalf("GetList: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/list/list1" {
		t.Errorf("path = %q, want /list/list1", gotPath)
	}
	m := out.(map[string]any)
	if m["id"] != "list1" || m["name"] != "My List" {
		t.Errorf("decoded = %+v", m)
	}
}

func TestCreateListInFolderPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"newlist"}`))
	})

	_, err := c.CreateListInFolder(context.Background(), "folder1", map[string]any{"name": "Sprint 1"})
	if err != nil {
		t.Fatalf("CreateListInFolder: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/folder/folder1/list" {
		t.Errorf("path = %q, want /folder/folder1/list", gotPath)
	}
	if gotBody["name"] != "Sprint 1" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestCreateFolderlessListPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"newlist"}`))
	})

	_, err := c.CreateFolderlessList(context.Background(), "space1", map[string]any{"name": "Backlog"})
	if err != nil {
		t.Fatalf("CreateFolderlessList: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/space/space1/list" {
		t.Errorf("path = %q, want /space/space1/list", gotPath)
	}
	if gotBody["name"] != "Backlog" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestUpdateListPutsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"list1","name":"Renamed"}`))
	})

	_, err := c.UpdateList(context.Background(), "list1", map[string]any{"name": "Renamed"})
	if err != nil {
		t.Fatalf("UpdateList: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	if gotPath != "/list/list1" {
		t.Errorf("path = %q, want /list/list1", gotPath)
	}
	if gotBody["name"] != "Renamed" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestDeleteListHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	})

	if err := c.DeleteList(context.Background(), "list1"); err != nil {
		t.Fatalf("DeleteList: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/list/list1" {
		t.Errorf("path = %q, want /list/list1", gotPath)
	}
}

func TestAddTaskToListHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	})

	if err := c.AddTaskToList(context.Background(), "list1", "task1"); err != nil {
		t.Fatalf("AddTaskToList: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/list/list1/task/task1" {
		t.Errorf("path = %q, want /list/list1/task/task1", gotPath)
	}
}

func TestRemoveTaskFromListHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	})

	if err := c.RemoveTaskFromList(context.Background(), "list1", "task1"); err != nil {
		t.Fatalf("RemoveTaskFromList: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/list/list1/task/task1" {
		t.Errorf("path = %q, want /list/list1/task/task1", gotPath)
	}
}

func TestCreateListFromTemplateInFolderPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"list1"}`))
	})

	_, err := c.CreateListFromTemplateInFolder(context.Background(), "folder1", "tmpl1", map[string]any{"name": "From Template"})
	if err != nil {
		t.Fatalf("CreateListFromTemplateInFolder: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/folder/folder1/list_template/tmpl1" {
		t.Errorf("path = %q, want /folder/folder1/list_template/tmpl1", gotPath)
	}
	if gotBody["name"] != "From Template" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestCreateListFromTemplateInSpacePostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"list1"}`))
	})

	_, err := c.CreateListFromTemplateInSpace(context.Background(), "space1", "tmpl1", map[string]any{"name": "From Template"})
	if err != nil {
		t.Fatalf("CreateListFromTemplateInSpace: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/space/space1/list_template/tmpl1" {
		t.Errorf("path = %q, want /space/space1/list_template/tmpl1", gotPath)
	}
	if gotBody["name"] != "From Template" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestListsMethodsReturnAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"err":"List not found","ECODE":"LIST_001"}`))
	})

	_, err := c.GetList(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest || apiErr.ECode != "LIST_001" || apiErr.Err != "List not found" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
