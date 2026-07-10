package tools

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func reqWithArgs(args map[string]any) mcp.CallToolRequest {
	return mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
}

func TestParseDateTimeArg(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		wantMs  int64
		wantErr bool
	}{
		{"full datetime", "2026-07-10 09:30:00", 1783675800000, false},
		{"bare date defaults to midnight UTC", "2026-07-10", 1783641600000, false},
		{"unix epoch", "1970-01-01 00:00:00", 0, false},
		{"empty string", "", 0, true},
		{"wrong separator", "2026/07/10", 0, true},
		{"US-style date", "07/10/2026", 0, true},
		{"missing seconds", "2026-07-10 09:30", 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseDateTimeArg(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("parseDateTimeArg(%q) = %d, nil; want error", tc.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseDateTimeArg(%q) unexpected error: %v", tc.in, err)
			}
			if got != tc.wantMs {
				t.Errorf("parseDateTimeArg(%q) = %d, want %d", tc.in, got, tc.wantMs)
			}
		})
	}
}

func TestRequireDateTimeArg(t *testing.T) {
	t.Run("missing key", func(t *testing.T) {
		req := reqWithArgs(map[string]any{})
		if _, err := requireDateTimeArg(req, "start"); err == nil {
			t.Fatal("expected error for missing key")
		}
	})

	t.Run("malformed value", func(t *testing.T) {
		req := reqWithArgs(map[string]any{"start": "not-a-date"})
		if _, err := requireDateTimeArg(req, "start"); err == nil {
			t.Fatal("expected error for malformed value")
		}
	})

	t.Run("valid value", func(t *testing.T) {
		req := reqWithArgs(map[string]any{"start": "1970-01-01 00:00:01"})
		got, err := requireDateTimeArg(req, "start")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 1000 {
			t.Errorf("got %d, want 1000", got)
		}
	})
}

func TestSetDateTime(t *testing.T) {
	t.Run("absent key is a no-op", func(t *testing.T) {
		body := map[string]any{}
		req := reqWithArgs(map[string]any{})
		if err := setDateTime(body, req, "due_date"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := body["due_date"]; ok {
			t.Errorf("body[due_date] = %v, want absent", body["due_date"])
		}
	})

	t.Run("valid value stored as ms", func(t *testing.T) {
		body := map[string]any{}
		req := reqWithArgs(map[string]any{"due_date": "1970-01-01 00:00:01"})
		if err := setDateTime(body, req, "due_date"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if body["due_date"] != int64(1000) {
			t.Errorf("body[due_date] = %v, want 1000", body["due_date"])
		}
	})

	t.Run("invalid value returns error and leaves body untouched", func(t *testing.T) {
		body := map[string]any{}
		req := reqWithArgs(map[string]any{"due_date": "not-a-date"})
		if err := setDateTime(body, req, "due_date"); err == nil {
			t.Fatal("expected error for malformed value")
		}
		if _, ok := body["due_date"]; ok {
			t.Errorf("body[due_date] = %v, want absent after error", body["due_date"])
		}
	})
}

func TestFormatMsValue(t *testing.T) {
	cases := []struct {
		name     string
		v        any
		dateOnly bool
		want     string
		wantOK   bool
	}{
		{"float64 epoch, full datetime", float64(1000), false, "1970-01-01 00:00:01", true},
		{"int64 epoch, full datetime", int64(1000), false, "1970-01-01 00:00:01", true},
		{"numeric string epoch", "1690000000000", false, "2023-07-22 04:26:40", true},
		{"bare date grain", float64(1783641600000), true, "2026-07-10", true},
		{"empty string", "", false, "", false},
		{"non-numeric string", "not-a-number", false, "", false},
		{"unsupported type (bool)", true, false, "", false},
		{"nil value", nil, false, "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := formatMsValue(tc.v, tc.dateOnly)
			if ok != tc.wantOK {
				t.Fatalf("formatMsValue(%v) ok = %v, want %v", tc.v, ok, tc.wantOK)
			}
			if ok && got != tc.want {
				t.Errorf("formatMsValue(%v) = %q, want %q", tc.v, got, tc.want)
			}
		})
	}
}

func TestConvertDateTimes(t *testing.T) {
	t.Run("full datetime key at top level", func(t *testing.T) {
		v := convertDateTimes(map[string]any{"date_created": float64(1000)}, nil)
		m := v.(map[string]any)
		if m["date_created"] != "1970-01-01 00:00:01" {
			t.Errorf("date_created = %v", m["date_created"])
		}
	})

	t.Run("bare date key renders date only", func(t *testing.T) {
		v := convertDateTimes(map[string]any{"due_date": "1783641600000"}, nil)
		m := v.(map[string]any)
		if m["due_date"] != "2026-07-10" {
			t.Errorf("due_date = %v", m["due_date"])
		}
	})

	t.Run("override forces full datetime on a normally bare-date key", func(t *testing.T) {
		v := convertDateTimes(map[string]any{"start_date": int64(0)}, map[string]bool{"start_date": true})
		m := v.(map[string]any)
		if m["start_date"] != "1970-01-01 00:00:00" {
			t.Errorf("start_date = %v, want full datetime override applied", m["start_date"])
		}
	})

	t.Run("non-numeric value under an allowlisted key is left untouched", func(t *testing.T) {
		v := convertDateTimes(map[string]any{"due_date": nil}, nil)
		m := v.(map[string]any)
		if m["due_date"] != nil {
			t.Errorf("due_date = %v, want nil left untouched", m["due_date"])
		}
	})

	t.Run("non-allowlisted key is left untouched even if numeric", func(t *testing.T) {
		v := convertDateTimes(map[string]any{"duration_ms": float64(60000)}, nil)
		m := v.(map[string]any)
		if m["duration_ms"] != float64(60000) {
			t.Errorf("duration_ms = %v, want unchanged 60000", m["duration_ms"])
		}
	})

	t.Run("nested conversion inside a slice of maps", func(t *testing.T) {
		v := convertDateTimes(map[string]any{
			"tasks": []any{
				map[string]any{"date_created": float64(1000)},
				map[string]any{"date_created": float64(2000)},
			},
		}, nil)
		m := v.(map[string]any)
		tasks := m["tasks"].([]any)
		if tasks[0].(map[string]any)["date_created"] != "1970-01-01 00:00:01" {
			t.Errorf("tasks[0].date_created = %v", tasks[0].(map[string]any)["date_created"])
		}
		if tasks[1].(map[string]any)["date_created"] != "1970-01-01 00:00:02" {
			t.Errorf("tasks[1].date_created = %v", tasks[1].(map[string]any)["date_created"])
		}
	})
}
