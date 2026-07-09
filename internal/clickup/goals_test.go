package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestListGoalsHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"goals":[{"id":"g1"}]}`))
	})

	out, err := c.ListGoals(context.Background(), "team1")
	if err != nil {
		t.Fatalf("ListGoals: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/team/team1/goal" {
		t.Errorf("path = %q, want /team/team1/goal", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded response is not a map: %T", out)
	}
	goals, ok := m["goals"].([]any)
	if !ok || len(goals) != 1 {
		t.Errorf("goals = %+v", m["goals"])
	}
}

func TestCreateGoalPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"g1"}`))
	})

	out, err := c.CreateGoal(context.Background(), "team1", map[string]any{"name": "Grow revenue"})
	if err != nil {
		t.Fatalf("CreateGoal: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/team/team1/goal" {
		t.Errorf("path = %q, want /team/team1/goal", gotPath)
	}
	if gotBody["name"] != "Grow revenue" {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok || m["id"] != "g1" {
		t.Errorf("decoded response = %+v", out)
	}
}

func TestGetGoalHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"id":"g1"}`))
	})

	out, err := c.GetGoal(context.Background(), "g1")
	if err != nil {
		t.Fatalf("GetGoal: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/goal/g1" {
		t.Errorf("path = %q, want /goal/g1", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok || m["id"] != "g1" {
		t.Errorf("decoded response = %+v", out)
	}
}

func TestUpdateGoalPutsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"g1"}`))
	})

	if _, err := c.UpdateGoal(context.Background(), "g1", map[string]any{"name": "New name"}); err != nil {
		t.Fatalf("UpdateGoal: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	if gotPath != "/goal/g1" {
		t.Errorf("path = %q, want /goal/g1", gotPath)
	}
	if gotBody["name"] != "New name" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestDeleteGoalHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{}`))
	})

	if err := c.DeleteGoal(context.Background(), "g1"); err != nil {
		t.Fatalf("DeleteGoal: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/goal/g1" {
		t.Errorf("path = %q, want /goal/g1", gotPath)
	}
}

func TestCreateKeyResultPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"kr1"}`))
	})

	out, err := c.CreateKeyResult(context.Background(), "g1", map[string]any{"name": "Sign 10 deals"})
	if err != nil {
		t.Fatalf("CreateKeyResult: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/goal/g1/key_result" {
		t.Errorf("path = %q, want /goal/g1/key_result", gotPath)
	}
	if gotBody["name"] != "Sign 10 deals" {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok || m["id"] != "kr1" {
		t.Errorf("decoded response = %+v", out)
	}
}

func TestUpdateKeyResultPutsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"kr1"}`))
	})

	if _, err := c.UpdateKeyResult(context.Background(), "kr1", map[string]any{"steps_current": float64(5)}); err != nil {
		t.Fatalf("UpdateKeyResult: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	if gotPath != "/key_result/kr1" {
		t.Errorf("path = %q, want /key_result/kr1", gotPath)
	}
	if gotBody["steps_current"] != float64(5) {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestDeleteKeyResultHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{}`))
	})

	if err := c.DeleteKeyResult(context.Background(), "kr1"); err != nil {
		t.Fatalf("DeleteKeyResult: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/key_result/kr1" {
		t.Errorf("path = %q, want /key_result/kr1", gotPath)
	}
}

func TestGetGoalAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"err":"Goal not found","ECODE":"GOAL_001"}`))
	})

	_, err := c.GetGoal(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusNotFound || apiErr.ECode != "GOAL_001" || apiErr.Err != "Goal not found" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
