package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestInviteGuestPostsExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"guest":{"id":1}}`))
	})

	out, err := c.InviteGuest(context.Background(), "111", map[string]any{"email": "a@b.com"})
	if err != nil {
		t.Fatalf("InviteGuest: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/team/111/guest" {
		t.Errorf("path = %q, want /team/111/guest", gotPath)
	}
	if gotBody["email"] != "a@b.com" {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok || m["guest"] == nil {
		t.Errorf("decoded = %+v", out)
	}
}

func TestGetGuestHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"guest":{"id":2}}`))
	})

	out, err := c.GetGuest(context.Background(), "111", "222")
	if err != nil {
		t.Fatalf("GetGuest: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/team/111/guest/222" {
		t.Errorf("path = %q, want /team/111/guest/222", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok || m["guest"] == nil {
		t.Errorf("decoded = %+v", out)
	}
}

func TestUpdateGuestPutsExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"guest":{"id":2}}`))
	})

	_, err := c.UpdateGuest(context.Background(), "111", "222", map[string]any{"can_edit_tags": true})
	if err != nil {
		t.Fatalf("UpdateGuest: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	if gotPath != "/team/111/guest/222" {
		t.Errorf("path = %q, want /team/111/guest/222", gotPath)
	}
	if gotBody["can_edit_tags"] != true {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestRemoveGuestFromWorkspaceHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	if err := c.RemoveGuestFromWorkspace(context.Background(), "111", "222"); err != nil {
		t.Fatalf("RemoveGuestFromWorkspace: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/team/111/guest/222" {
		t.Errorf("path = %q, want /team/111/guest/222", gotPath)
	}
}

func TestAddGuestToSpacePostsExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"guest":{"id":2}}`))
	})

	_, err := c.AddGuestToSpace(context.Background(), "s1", "222", map[string]any{"permission_level": "read"})
	if err != nil {
		t.Fatalf("AddGuestToSpace: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/space/s1/guest/222" {
		t.Errorf("path = %q, want /space/s1/guest/222", gotPath)
	}
	if gotBody["permission_level"] != "read" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestRemoveGuestFromSpaceHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	if err := c.RemoveGuestFromSpace(context.Background(), "s1", "222"); err != nil {
		t.Fatalf("RemoveGuestFromSpace: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/space/s1/guest/222" {
		t.Errorf("path = %q, want /space/s1/guest/222", gotPath)
	}
}

func TestAddGuestToFolderPostsExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"guest":{"id":2}}`))
	})

	_, err := c.AddGuestToFolder(context.Background(), "f1", "222", map[string]any{"permission_level": "edit"})
	if err != nil {
		t.Fatalf("AddGuestToFolder: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/folder/f1/guest/222" {
		t.Errorf("path = %q, want /folder/f1/guest/222", gotPath)
	}
	if gotBody["permission_level"] != "edit" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestRemoveGuestFromFolderHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	if err := c.RemoveGuestFromFolder(context.Background(), "f1", "222"); err != nil {
		t.Fatalf("RemoveGuestFromFolder: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/folder/f1/guest/222" {
		t.Errorf("path = %q, want /folder/f1/guest/222", gotPath)
	}
}

func TestAddGuestToListPostsExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"guest":{"id":2}}`))
	})

	_, err := c.AddGuestToList(context.Background(), "l1", "222", map[string]any{"permission_level": "comment"})
	if err != nil {
		t.Fatalf("AddGuestToList: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/list/l1/guest/222" {
		t.Errorf("path = %q, want /list/l1/guest/222", gotPath)
	}
	if gotBody["permission_level"] != "comment" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestRemoveGuestFromListHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	if err := c.RemoveGuestFromList(context.Background(), "l1", "222"); err != nil {
		t.Fatalf("RemoveGuestFromList: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/list/l1/guest/222" {
		t.Errorf("path = %q, want /list/l1/guest/222", gotPath)
	}
}

func TestAddGuestToTaskPostsExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"guest":{"id":2}}`))
	})

	_, err := c.AddGuestToTask(context.Background(), "t1", "222", map[string]any{"permission_level": "read"})
	if err != nil {
		t.Fatalf("AddGuestToTask: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/task/t1/guest/222" {
		t.Errorf("path = %q, want /task/t1/guest/222", gotPath)
	}
	if gotBody["permission_level"] != "read" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestRemoveGuestFromTaskHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	if err := c.RemoveGuestFromTask(context.Background(), "t1", "222"); err != nil {
		t.Fatalf("RemoveGuestFromTask: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/task/t1/guest/222" {
		t.Errorf("path = %q, want /task/t1/guest/222", gotPath)
	}
}

func TestGuestsAPIErrorReturnsAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"err":"Guest limit reached","ECODE":"GUEST_001"}`))
	})

	_, err := c.InviteGuest(context.Background(), "111", map[string]any{"email": "a@b.com"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusForbidden || apiErr.ECode != "GUEST_001" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
