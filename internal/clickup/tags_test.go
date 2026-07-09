package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestAddTaskTagHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{}`))
	})

	if err := c.AddTaskTag(context.Background(), "task1", "urgent tag"); err != nil {
		t.Fatalf("AddTaskTag: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/task/task1/tag/urgent tag" {
		t.Errorf("path = %q, want /task/task1/tag/urgent tag", gotPath)
	}
}

func TestRemoveTaskTagHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{}`))
	})

	if err := c.RemoveTaskTag(context.Background(), "task1", "urgent"); err != nil {
		t.Fatalf("RemoveTaskTag: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/task/task1/tag/urgent" {
		t.Errorf("path = %q, want /task/task1/tag/urgent", gotPath)
	}
}

func TestListSpaceTagsDecodesResponse(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"tags":[{"name":"urgent"}]}`))
	})

	out, err := c.ListSpaceTags(context.Background(), "space1")
	if err != nil {
		t.Fatalf("ListSpaceTags: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/space/space1/tag" {
		t.Errorf("path = %q, want /space/space1/tag", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded response is not a map: %T", out)
	}
	tags, ok := m["tags"].([]any)
	if !ok || len(tags) != 1 {
		t.Errorf("tags = %+v", m["tags"])
	}
}

func TestCreateSpaceTagPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"tag":{"name":"urgent"}}`))
	})

	body := map[string]any{"tag": map[string]any{"name": "urgent", "tag_fg": "#ffffff", "tag_bg": "#000000"}}
	out, err := c.CreateSpaceTag(context.Background(), "space1", body)
	if err != nil {
		t.Fatalf("CreateSpaceTag: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/space/space1/tag" {
		t.Errorf("path = %q, want /space/space1/tag", gotPath)
	}
	tag, ok := gotBody["tag"].(map[string]any)
	if !ok || tag["name"] != "urgent" {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded response is not a map: %T", out)
	}
	if _, ok := m["tag"]; !ok {
		t.Errorf("decoded response missing tag: %+v", m)
	}
}

func TestUpdateSpaceTagPutsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{}`))
	})

	body := map[string]any{"tag": map[string]any{"name": "renamed"}}
	if _, err := c.UpdateSpaceTag(context.Background(), "space1", "old tag", body); err != nil {
		t.Fatalf("UpdateSpaceTag: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	if gotPath != "/space/space1/tag/old tag" {
		t.Errorf("path = %q, want /space/space1/tag/old tag", gotPath)
	}
	tag, ok := gotBody["tag"].(map[string]any)
	if !ok || tag["name"] != "renamed" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestDeleteSpaceTagHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{}`))
	})

	if err := c.DeleteSpaceTag(context.Background(), "space1", "urgent"); err != nil {
		t.Fatalf("DeleteSpaceTag: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/space/space1/tag/urgent" {
		t.Errorf("path = %q, want /space/space1/tag/urgent", gotPath)
	}
}

func TestCreateSpaceTagAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"err":"Not authorized","ECODE":"TAG_001"}`))
	})

	_, err := c.CreateSpaceTag(context.Background(), "space1", map[string]any{"tag": map[string]any{"name": "x"}})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized || apiErr.ECode != "TAG_001" || apiErr.Err != "Not authorized" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
