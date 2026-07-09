package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestInviteWorkspaceUserPostsExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"team":{"id":"123"}}`))
	})

	out, err := c.InviteWorkspaceUser(context.Background(), "123", map[string]any{"email": "a@b.com", "admin": false})
	if err != nil {
		t.Fatalf("InviteWorkspaceUser: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/team/123/user" {
		t.Errorf("path = %q, want /team/123/user", gotPath)
	}
	if gotBody["email"] != "a@b.com" || gotBody["admin"] != false {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok || m["team"] == nil {
		t.Errorf("decoded = %+v", out)
	}
}

func TestUpdateWorkspaceUserPutsExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"team":{"id":"123"}}`))
	})

	_, err := c.UpdateWorkspaceUser(context.Background(), "123", "456", map[string]any{"admin": true})
	if err != nil {
		t.Fatalf("UpdateWorkspaceUser: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	if gotPath != "/team/123/user/456" {
		t.Errorf("path = %q, want /team/123/user/456", gotPath)
	}
	if gotBody["admin"] != true {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestRemoveWorkspaceUserHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	if err := c.RemoveWorkspaceUser(context.Background(), "123", "456"); err != nil {
		t.Fatalf("RemoveWorkspaceUser: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/team/123/user/456" {
		t.Errorf("path = %q, want /team/123/user/456", gotPath)
	}
}

func TestUsersAPIErrorReturnsAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"err":"Not authorized","ECODE":"OAUTH_017"}`))
	})

	_, err := c.InviteWorkspaceUser(context.Background(), "123", map[string]any{"email": "a@b.com"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized || apiErr.ECode != "OAUTH_017" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
