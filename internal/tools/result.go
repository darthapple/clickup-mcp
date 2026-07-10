// Package tools registers every ClickUp REST endpoint as an MCP tool on top
// of the internal/clickup client.
package tools

import (
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"clickup-mcp/internal/clickup"
)

// JSONResult renders v (typically the raw decoded ClickUp API response) as
// the tool call's result, converting every recognized Unix-ms date/time
// field (see dateTimeKeys in datetime.go) to a human-readable UTC string
// first. overrides names keys that should render as full datetime for this
// call even though they're normally bare-date (dateOnlyKeys) — see
// convertDateTimes.
func JSONResult(v any, overrides ...string) (*mcp.CallToolResult, error) {
	ov := make(map[string]bool, len(overrides))
	for _, k := range overrides {
		ov[k] = true
	}
	return mcp.NewToolResultJSON(convertDateTimes(v, ov))
}

// ErrorResult renders err as a failed-but-not-crashed tool result: ClickUp
// API errors get a readable status/code/message, anything else falls back to
// err.Error(). The returned Go error is always nil so callers can write
// `return tools.ErrorResult(err)` directly from a handler.
func ErrorResult(err error) (*mcp.CallToolResult, error) {
	var apiErr *clickup.APIError
	if errors.As(err, &apiErr) {
		return mcp.NewToolResultError(fmt.Sprintf("ClickUp API error %d [%s]: %s", apiErr.StatusCode, apiErr.ECode, apiErr.Err)), nil
	}
	return mcp.NewToolResultError(err.Error()), nil
}
