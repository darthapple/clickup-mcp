package clickup

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestGetSharedHierarchyHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"shared":{"tasks":[],"lists":[],"folders":[]}}`))
	})

	out, err := c.GetSharedHierarchy(context.Background(), "999")
	if err != nil {
		t.Fatalf("GetSharedHierarchy: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/team/999/shared" {
		t.Errorf("path = %q, want /team/999/shared", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded type = %T, want map[string]any", out)
	}
	if _, ok := m["shared"]; !ok {
		t.Errorf("decoded = %+v, want key %q", m, "shared")
	}
}

func TestGetSharedHierarchyReturnsAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"err":"Team not found","ECODE":"TEAM_001"}`))
	})

	_, err := c.GetSharedHierarchy(context.Background(), "999")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusNotFound || apiErr.ECode != "TEAM_001" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
