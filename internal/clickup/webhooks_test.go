package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestListWebhooksHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"webhooks":[{"id":"wh1"}]}`))
	})

	out, err := c.ListWebhooks(context.Background(), "team1")
	if err != nil {
		t.Fatalf("ListWebhooks: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/team/team1/webhook" {
		t.Errorf("path = %q, want /team/team1/webhook", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded response is not a map: %T", out)
	}
	webhooks, ok := m["webhooks"].([]any)
	if !ok || len(webhooks) != 1 {
		t.Errorf("webhooks = %+v", m["webhooks"])
	}
}

func TestCreateWebhookPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"wh1"}`))
	})

	out, err := c.CreateWebhook(context.Background(), "team1", map[string]any{
		"endpoint": "https://example.com/hook",
		"events":   []any{"taskCreated", "taskUpdated"},
	})
	if err != nil {
		t.Fatalf("CreateWebhook: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/team/team1/webhook" {
		t.Errorf("path = %q, want /team/team1/webhook", gotPath)
	}
	if gotBody["endpoint"] != "https://example.com/hook" {
		t.Errorf("body = %+v", gotBody)
	}
	events, ok := gotBody["events"].([]any)
	if !ok || len(events) != 2 {
		t.Errorf("events = %+v", gotBody["events"])
	}
	m, ok := out.(map[string]any)
	if !ok || m["id"] != "wh1" {
		t.Errorf("decoded response = %+v", out)
	}
}

func TestUpdateWebhookPutsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"wh1"}`))
	})

	if _, err := c.UpdateWebhook(context.Background(), "wh1", map[string]any{"status": "active"}); err != nil {
		t.Fatalf("UpdateWebhook: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	if gotPath != "/webhook/wh1" {
		t.Errorf("path = %q, want /webhook/wh1", gotPath)
	}
	if gotBody["status"] != "active" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestDeleteWebhookHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{}`))
	})

	if err := c.DeleteWebhook(context.Background(), "wh1"); err != nil {
		t.Fatalf("DeleteWebhook: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/webhook/wh1" {
		t.Errorf("path = %q, want /webhook/wh1", gotPath)
	}
}

func TestCreateWebhookAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"err":"Invalid endpoint","ECODE":"WEBHOOK_001"}`))
	})

	_, err := c.CreateWebhook(context.Background(), "team1", map[string]any{"endpoint": "not-a-url"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest || apiErr.ECode != "WEBHOOK_001" || apiErr.Err != "Invalid endpoint" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
