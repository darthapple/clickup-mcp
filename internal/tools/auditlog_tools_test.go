package tools

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupCreateAuditLogReport(t *testing.T) {
	t.Run("argument wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"report1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterAuditLogTools(s, c)
		res := callTool(t, s, "clickup_create_audit_log_report", map[string]any{
			"team_id":     "123",
			"start_date":  "1970-01-01 00:00:01",
			"end_date":    "1970-01-01 00:00:02",
			"event_types": []any{"taskCreated", "taskDeleted"},
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/workspaces/123/auditlogs" {
			t.Errorf("path = %q, want /workspaces/123/auditlogs", gotPath)
		}
		if gotBody["start_date"] != float64(1000) || gotBody["end_date"] != float64(2000) {
			t.Errorf("body = %+v", gotBody)
		}
		types, ok := gotBody["event_types"].([]any)
		if !ok || len(types) != 2 {
			t.Errorf("body[event_types] = %+v", gotBody["event_types"])
		}
	})

	t.Run("no optional args yields empty body", func(t *testing.T) {
		var gotBody map[string]any
		bodyRaw := []byte{}
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			buf := make([]byte, 1024)
			n, _ := r.Body.Read(buf)
			bodyRaw = buf[:n]
			if len(bodyRaw) > 0 {
				_ = json.Unmarshal(bodyRaw, &gotBody)
			}
			_, _ = w.Write([]byte(`{"id":"report1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterAuditLogTools(s, c)
		res := callTool(t, s, "clickup_create_audit_log_report", map[string]any{"team_id": "123"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if len(gotBody) != 0 {
			t.Errorf("body = %+v, want empty", gotBody)
		}
	})

	t.Run("error passthrough (expected on non-Enterprise workspaces)", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(403)
			_, _ = w.Write([]byte(`{"err":"feature not available","ECODE":"AUDIT_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterAuditLogTools(s, c)
		res := callTool(t, s, "clickup_create_audit_log_report", map[string]any{"team_id": "123"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 403 [AUDIT_001]: feature not available"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}
