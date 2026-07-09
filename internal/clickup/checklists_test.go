package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestCreateChecklistPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"checklist":{"id":"chk1"}}`))
	})

	out, err := c.CreateChecklist(context.Background(), "task1", map[string]any{"name": "To Do"})
	if err != nil {
		t.Fatalf("CreateChecklist: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/task/task1/checklist" {
		t.Errorf("path = %q, want /task/task1/checklist", gotPath)
	}
	if gotBody["name"] != "To Do" {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("out is not a map: %T", out)
	}
	if _, ok := m["checklist"]; !ok {
		t.Errorf("decoded = %+v, missing checklist", m)
	}
}

func TestUpdateChecklistPutsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"checklist":{"id":"chk1","name":"Renamed"}}`))
	})

	_, err := c.UpdateChecklist(context.Background(), "chk1", map[string]any{"name": "Renamed"})
	if err != nil {
		t.Fatalf("UpdateChecklist: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	if gotPath != "/checklist/chk1" {
		t.Errorf("path = %q, want /checklist/chk1", gotPath)
	}
	if gotBody["name"] != "Renamed" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestDeleteChecklistHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	})

	if err := c.DeleteChecklist(context.Background(), "chk1"); err != nil {
		t.Fatalf("DeleteChecklist: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/checklist/chk1" {
		t.Errorf("path = %q, want /checklist/chk1", gotPath)
	}
}

func TestCreateChecklistItemPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"checklist":{"id":"chk1"}}`))
	})

	_, err := c.CreateChecklistItem(context.Background(), "chk1", map[string]any{"name": "Step 1"})
	if err != nil {
		t.Fatalf("CreateChecklistItem: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/checklist/chk1/checklist_item" {
		t.Errorf("path = %q, want /checklist/chk1/checklist_item", gotPath)
	}
	if gotBody["name"] != "Step 1" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestUpdateChecklistItemPutsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"checklist":{"id":"chk1"}}`))
	})

	_, err := c.UpdateChecklistItem(context.Background(), "chk1", "item1", map[string]any{"resolved": true})
	if err != nil {
		t.Fatalf("UpdateChecklistItem: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	if gotPath != "/checklist/chk1/checklist_item/item1" {
		t.Errorf("path = %q, want /checklist/chk1/checklist_item/item1", gotPath)
	}
	if gotBody["resolved"] != true {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestDeleteChecklistItemHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	})

	if err := c.DeleteChecklistItem(context.Background(), "chk1", "item1"); err != nil {
		t.Fatalf("DeleteChecklistItem: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/checklist/chk1/checklist_item/item1" {
		t.Errorf("path = %q, want /checklist/chk1/checklist_item/item1", gotPath)
	}
}

func TestChecklistsMethodsReturnAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"err":"Checklist invalid","ECODE":"CHECKLIST_001"}`))
	})

	_, err := c.CreateChecklist(context.Background(), "task1", map[string]any{"name": ""})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest || apiErr.ECode != "CHECKLIST_001" || apiErr.Err != "Checklist invalid" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
