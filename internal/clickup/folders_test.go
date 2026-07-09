package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestListFoldersHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath, gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"folders":[{"id":"1"}]}`))
	})

	out, err := c.ListFolders(context.Background(), "space1", true, true)
	if err != nil {
		t.Fatalf("ListFolders: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/space/space1/folder" {
		t.Errorf("path = %q, want /space/space1/folder", gotPath)
	}
	if gotQuery != "archived=true" {
		t.Errorf("query = %q, want archived=true", gotQuery)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("out is not a map: %T", out)
	}
	if _, ok := m["folders"]; !ok {
		t.Errorf("decoded = %+v, missing folders", m)
	}
}

func TestListFoldersOmitsArchivedWhenNotSet(t *testing.T) {
	var gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{}`))
	})

	if _, err := c.ListFolders(context.Background(), "space1", false, false); err != nil {
		t.Fatalf("ListFolders: %v", err)
	}
	if gotQuery != "" {
		t.Errorf("query = %q, want empty (archived absent)", gotQuery)
	}
}

func TestGetFolderHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"id":"folder1","name":"My Folder"}`))
	})

	out, err := c.GetFolder(context.Background(), "folder1")
	if err != nil {
		t.Fatalf("GetFolder: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/folder/folder1" {
		t.Errorf("path = %q, want /folder/folder1", gotPath)
	}
	m := out.(map[string]any)
	if m["id"] != "folder1" || m["name"] != "My Folder" {
		t.Errorf("decoded = %+v", m)
	}
}

func TestCreateFolderPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"newfolder"}`))
	})

	_, err := c.CreateFolder(context.Background(), "space1", map[string]any{"name": "Q1 Planning"})
	if err != nil {
		t.Fatalf("CreateFolder: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/space/space1/folder" {
		t.Errorf("path = %q, want /space/space1/folder", gotPath)
	}
	if gotBody["name"] != "Q1 Planning" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestUpdateFolderPutsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"folder1","name":"Renamed"}`))
	})

	_, err := c.UpdateFolder(context.Background(), "folder1", map[string]any{"name": "Renamed"})
	if err != nil {
		t.Fatalf("UpdateFolder: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	if gotPath != "/folder/folder1" {
		t.Errorf("path = %q, want /folder/folder1", gotPath)
	}
	if gotBody["name"] != "Renamed" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestDeleteFolderHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	})

	if err := c.DeleteFolder(context.Background(), "folder1"); err != nil {
		t.Fatalf("DeleteFolder: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/folder/folder1" {
		t.Errorf("path = %q, want /folder/folder1", gotPath)
	}
}

func TestCreateFolderFromTemplatePostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"folder1"}`))
	})

	_, err := c.CreateFolderFromTemplate(context.Background(), "space1", "tmpl1", map[string]any{"name": "From Template"})
	if err != nil {
		t.Fatalf("CreateFolderFromTemplate: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/space/space1/folder_template/tmpl1" {
		t.Errorf("path = %q, want /space/space1/folder_template/tmpl1", gotPath)
	}
	if gotBody["name"] != "From Template" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestFoldersMethodsReturnAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"err":"Folder access denied","ECODE":"FOLDER_001"}`))
	})

	_, err := c.GetFolder(context.Background(), "folder1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusForbidden || apiErr.ECode != "FOLDER_001" || apiErr.Err != "Folder access denied" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
