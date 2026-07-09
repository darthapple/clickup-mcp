package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestCreateTaskFromTemplatePostsExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"newtask123"}`))
	})

	out, err := c.CreateTaskFromTemplate(context.Background(), "list1", "tmpl1", map[string]any{"name": "From template"})
	if err != nil {
		t.Fatalf("CreateTaskFromTemplate: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/list/list1/taskTemplate/tmpl1" {
		t.Errorf("path = %q, want /list/list1/taskTemplate/tmpl1", gotPath)
	}
	if gotBody["name"] != "From template" {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok || m["id"] != "newtask123" {
		t.Errorf("decoded = %+v", out)
	}
}

func TestTemplatesAPIErrorReturnsAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"err":"Template not found","ECODE":"TEMPLATE_001"}`))
	})

	_, err := c.CreateTaskFromTemplate(context.Background(), "list1", "tmpl1", map[string]any{"name": "x"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest || apiErr.ECode != "TEMPLATE_001" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
