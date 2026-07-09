package clickup

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestListCustomTaskTypesHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"custom_items":[{"id":1,"name":"Bug"}]}`))
	})

	out, err := c.ListCustomTaskTypes(context.Background(), "999")
	if err != nil {
		t.Fatalf("ListCustomTaskTypes: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/team/999/custom_item" {
		t.Errorf("path = %q, want /team/999/custom_item", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded type = %T, want map[string]any", out)
	}
	items, ok := m["custom_items"].([]any)
	if !ok || len(items) != 1 {
		t.Errorf("custom_items = %+v", m["custom_items"])
	}
}

func TestListCustomTaskTypesReturnsAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"err":"Invalid token","ECODE":"OAUTH_001"}`))
	})

	_, err := c.ListCustomTaskTypes(context.Background(), "999")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized || apiErr.ECode != "OAUTH_001" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
