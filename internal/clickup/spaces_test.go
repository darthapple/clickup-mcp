package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestListSpacesHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath, gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"spaces":[{"id":"1"}]}`))
	})

	out, err := c.ListSpaces(context.Background(), "team1", true, true)
	if err != nil {
		t.Fatalf("ListSpaces: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/team/team1/space" {
		t.Errorf("path = %q, want /team/team1/space", gotPath)
	}
	if gotQuery != "archived=true" {
		t.Errorf("query = %q, want archived=true", gotQuery)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("out is not a map: %T", out)
	}
	if _, ok := m["spaces"]; !ok {
		t.Errorf("decoded = %+v, missing spaces", m)
	}
}

func TestListSpacesOmitsArchivedWhenNotSet(t *testing.T) {
	var gotQuery string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{}`))
	})

	if _, err := c.ListSpaces(context.Background(), "team1", false, false); err != nil {
		t.Fatalf("ListSpaces: %v", err)
	}
	if gotQuery != "" {
		t.Errorf("query = %q, want empty (archived absent)", gotQuery)
	}
}

func TestCreateSpacePostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"newspace"}`))
	})

	_, err := c.CreateSpace(context.Background(), "team1", map[string]any{"name": "Engineering", "multiple_assignees": true})
	if err != nil {
		t.Fatalf("CreateSpace: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/team/team1/space" {
		t.Errorf("path = %q, want /team/team1/space", gotPath)
	}
	if gotBody["name"] != "Engineering" || gotBody["multiple_assignees"] != true {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestGetSpaceHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"id":"space1","name":"My Space"}`))
	})

	out, err := c.GetSpace(context.Background(), "space1")
	if err != nil {
		t.Fatalf("GetSpace: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/space/space1" {
		t.Errorf("path = %q, want /space/space1", gotPath)
	}
	m := out.(map[string]any)
	if m["id"] != "space1" || m["name"] != "My Space" {
		t.Errorf("decoded = %+v", m)
	}
}

func TestUpdateSpacePutsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"space1","name":"Renamed"}`))
	})

	_, err := c.UpdateSpace(context.Background(), "space1", map[string]any{"name": "Renamed"})
	if err != nil {
		t.Fatalf("UpdateSpace: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	if gotPath != "/space/space1" {
		t.Errorf("path = %q, want /space/space1", gotPath)
	}
	if gotBody["name"] != "Renamed" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestDeleteSpaceHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	})

	if err := c.DeleteSpace(context.Background(), "space1"); err != nil {
		t.Fatalf("DeleteSpace: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/space/space1" {
		t.Errorf("path = %q, want /space/space1", gotPath)
	}
}

func TestSpacesMethodsReturnAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"err":"Space access denied","ECODE":"SPACE_001"}`))
	})

	_, err := c.GetSpace(context.Background(), "space1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized || apiErr.ECode != "SPACE_001" || apiErr.Err != "Space access denied" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
