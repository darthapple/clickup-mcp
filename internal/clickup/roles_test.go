package clickup

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestListCustomRolesHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"custom_roles":[{"id":1,"name":"QA"}]}`))
	})

	out, err := c.ListCustomRoles(context.Background(), "123")
	if err != nil {
		t.Fatalf("ListCustomRoles: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/team/123/customroles" {
		t.Errorf("path = %q, want /team/123/customroles", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok || m["custom_roles"] == nil {
		t.Errorf("decoded = %+v", out)
	}
}

func TestRolesAPIErrorReturnsAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"err":"Team not found","ECODE":"TEAM_001"}`))
	})

	_, err := c.ListCustomRoles(context.Background(), "123")
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
