package clickup

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestGetUserHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"user":{"id":1,"username":"fernando"}}`))
	})

	out, err := c.GetUser(context.Background())
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/user" {
		t.Errorf("path = %q, want /user", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded type = %T, want map[string]any", out)
	}
	user, ok := m["user"].(map[string]any)
	if !ok || user["username"] != "fernando" {
		t.Errorf("user = %+v", m["user"])
	}
}

func TestGetUserReturnsAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"err":"Invalid token","ECODE":"OAUTH_001"}`))
	})

	_, err := c.GetUser(context.Background())
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

func TestListWorkspacesHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"teams":[{"id":"999","name":"Kheperi"}]}`))
	})

	out, err := c.ListWorkspaces(context.Background())
	if err != nil {
		t.Fatalf("ListWorkspaces: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/team" {
		t.Errorf("path = %q, want /team", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded type = %T, want map[string]any", out)
	}
	teams, ok := m["teams"].([]any)
	if !ok || len(teams) != 1 {
		t.Errorf("teams = %+v", m["teams"])
	}
}

func TestListWorkspacesReturnsAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"err":"Invalid token","ECODE":"OAUTH_001"}`))
	})

	_, err := c.ListWorkspaces(context.Background())
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
