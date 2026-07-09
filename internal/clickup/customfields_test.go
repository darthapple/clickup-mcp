package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestListListFieldsHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"fields":[{"id":"f1"}]}`))
	})

	out, err := c.ListListFields(context.Background(), "list1")
	if err != nil {
		t.Fatalf("ListListFields: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/list/list1/field" {
		t.Errorf("path = %q, want /list/list1/field", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("out is not a map: %T", out)
	}
	if _, ok := m["fields"]; !ok {
		t.Errorf("decoded = %+v, missing fields", m)
	}
}

func TestListFolderFieldsHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"fields":[]}`))
	})

	if _, err := c.ListFolderFields(context.Background(), "folder1"); err != nil {
		t.Fatalf("ListFolderFields: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/folder/folder1/field" {
		t.Errorf("path = %q, want /folder/folder1/field", gotPath)
	}
}

func TestListSpaceFieldsHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"fields":[]}`))
	})

	if _, err := c.ListSpaceFields(context.Background(), "space1"); err != nil {
		t.Fatalf("ListSpaceFields: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/space/space1/field" {
		t.Errorf("path = %q, want /space/space1/field", gotPath)
	}
}

func TestListWorkspaceFieldsHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"fields":[]}`))
	})

	if _, err := c.ListWorkspaceFields(context.Background(), "team1"); err != nil {
		t.Fatalf("ListWorkspaceFields: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/team/team1/field" {
		t.Errorf("path = %q, want /team/team1/field", gotPath)
	}
}

func TestSetTaskCustomFieldPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"task1"}`))
	})

	_, err := c.SetTaskCustomField(context.Background(), "task1", "field1", map[string]any{"value": "high"})
	if err != nil {
		t.Fatalf("SetTaskCustomField: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/task/task1/field/field1" {
		t.Errorf("path = %q, want /task/task1/field/field1", gotPath)
	}
	if gotBody["value"] != "high" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestRemoveTaskCustomFieldHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	})

	if err := c.RemoveTaskCustomField(context.Background(), "task1", "field1"); err != nil {
		t.Fatalf("RemoveTaskCustomField: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/task/task1/field/field1" {
		t.Errorf("path = %q, want /task/task1/field/field1", gotPath)
	}
}

func TestCustomFieldsMethodsReturnAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"err":"Invalid field value","ECODE":"FIELD_001"}`))
	})

	_, err := c.SetTaskCustomField(context.Background(), "task1", "field1", map[string]any{"value": "bad"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest || apiErr.ECode != "FIELD_001" || apiErr.Err != "Invalid field value" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
