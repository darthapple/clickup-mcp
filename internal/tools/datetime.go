package tools

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// dateLayout and dateTimeLayout are the two accepted/produced UTC string
// formats for every date/time value crossing the MCP boundary — a bare
// calendar date, or a full datetime with second resolution. Sub-second
// precision is never preserved.
const (
	dateLayout     = "2006-01-02"
	dateTimeLayout = "2006-01-02 15:04:05"
)

// dateTimeKeys is the curated, best-effort set of ClickUp response field
// names known to hold Unix-ms timestamps. It is not exhaustive — ClickUp's
// API surface is large and not every resource's raw JSON has been audited —
// but covers every date field this codebase currently touches or passes
// through. Extend it as new date-bearing fields turn up in raw passthrough
// responses.
var dateTimeKeys = map[string]bool{
	"date_created":  true,
	"date_updated":  true,
	"date_closed":   true,
	"date_done":     true,
	"date_added":    true,
	"date_joined":   true,
	"date_invited":  true,
	"archived_date": true,
	"resolved_date": true,
	"date":          true,
	"at":            true,
	"start":         true,
	"end":           true,
	"end_date":      true,
	"due_date":      true,
	"start_date":    true,
}

// bareDateKeys is the subset of dateTimeKeys that render/accept a bare
// YYYY-MM-DD date instead of a full YYYY-MM-DD HH:MM:SS datetime, because
// the underlying ClickUp concept (a Task or Goal's due_date/start_date) is a
// calendar date, not a precise moment. Every other key in dateTimeKeys
// always renders as a full datetime.
var bareDateKeys = map[string]bool{
	"due_date":   true,
	"start_date": true,
}

// parseDateTimeArg parses a caller-supplied date/time string into a Unix ms
// epoch timestamp (UTC). Accepts either a full "YYYY-MM-DD HH:MM:SS" or a
// bare "YYYY-MM-DD" (midnight UTC assumed).
func parseDateTimeArg(s string) (int64, error) {
	if t, err := time.Parse(dateTimeLayout, s); err == nil {
		return t.UTC().UnixMilli(), nil
	}
	if t, err := time.Parse(dateLayout, s); err == nil {
		return t.UTC().UnixMilli(), nil
	}
	return 0, fmt.Errorf("must be a UTC date/time in \"YYYY-MM-DD HH:MM:SS\" or \"YYYY-MM-DD\" format, got %q", s)
}

// requireDateTimeArg reads key as a required string argument and parses it
// into Unix ms. Used where the parsed value is needed directly rather than
// stored into a request body map (e.g. report date ranges, time-entry
// filters).
func requireDateTimeArg(req mcp.CallToolRequest, key string) (int64, error) {
	s, err := req.RequireString(key)
	if err != nil {
		return 0, err
	}
	ms, err := parseDateTimeArg(s)
	if err != nil {
		return 0, fmt.Errorf("%s %w", key, err)
	}
	return ms, nil
}

// setDateTime parses key as a date/time string and stores it in body as a
// Unix ms epoch number, if the caller supplied it. Mirrors setFloat's
// omit-if-absent shape (body.go) but returns an error since, unlike a plain
// float, parsing this value can fail.
func setDateTime(body map[string]any, req mcp.CallToolRequest, key string) error {
	if !hasArg(req, key) {
		return nil
	}
	ms, err := parseDateTimeArg(req.GetString(key, ""))
	if err != nil {
		return fmt.Errorf("%s %w", key, err)
	}
	body[key] = ms
	return nil
}

// convertDateTimes recursively walks v — either the raw decoded ClickUp
// JSON passthrough or a hand-built map[string]any/[]any — and rewrites
// every value found under a dateTimeKeys key from a Unix-ms epoch (ClickUp
// encodes these as either a JSON number or a numeric string depending on
// endpoint) to a human-readable UTC string. overrides forces specific keys
// to render as full datetime even if they're in bareDateKeys, for the rare
// case where the same key name means something else in a given tool's
// response (e.g. "start_date" is a Task's calendar date, but a time-report
// tool's own echoed query-range boundary under the same key needs full
// precision). Values that aren't parseable as an epoch (null, non-numeric)
// are left untouched rather than corrupted.
func convertDateTimes(v any, overrides map[string]bool) any {
	switch val := v.(type) {
	case map[string]any:
		for k, sub := range val {
			if dateTimeKeys[k] {
				bare := bareDateKeys[k] && !overrides[k]
				if formatted, ok := formatMsValue(sub, bare); ok {
					val[k] = formatted
					continue
				}
			}
			val[k] = convertDateTimes(sub, overrides)
		}
		return val
	case []any:
		for i, item := range val {
			val[i] = convertDateTimes(item, overrides)
		}
		return val
	default:
		return v
	}
}

// formatMsValue renders a raw epoch-ms value as a UTC date or datetime
// string. v may be a float64/json number (raw ClickUp JSON decodes numbers
// this way), an int64 (hand-built response maps), or a numeric string
// (ClickUp encodes some ms fields as JSON strings to avoid float precision
// loss). Returns ok=false for anything else (null, non-numeric string),
// leaving the original value untouched.
func formatMsValue(v any, dateOnly bool) (string, bool) {
	var ms int64
	switch n := v.(type) {
	case string:
		if n == "" {
			return "", false
		}
		parsed, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return "", false
		}
		ms = parsed
	case float64:
		ms = int64(n)
	case int64:
		ms = n
	default:
		return "", false
	}
	t := time.UnixMilli(ms).UTC()
	if dateOnly {
		return t.Format(dateLayout), true
	}
	return t.Format(dateTimeLayout), true
}
