package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestCreateAuditLogReportPostsExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"report_id":"rep1"}`))
	})

	out, err := c.CreateAuditLogReport(context.Background(), "999", map[string]any{"start_date": float64(1000)})
	if err != nil {
		t.Fatalf("CreateAuditLogReport: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotPath != "/workspaces/999/auditlogs" {
		t.Errorf("path = %q, want /workspaces/999/auditlogs", gotPath)
	}
	if gotBody["start_date"] != float64(1000) {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok || m["report_id"] != "rep1" {
		t.Errorf("decoded = %+v", out)
	}
}

func TestCreateAuditLogReportReturnsAPIErrorForNonEnterprise(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"err":"Enterprise plan required","ECODE":"AUDIT_001"}`))
	})

	_, err := c.CreateAuditLogReport(context.Background(), "999", map[string]any{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusForbidden || apiErr.ECode != "AUDIT_001" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
