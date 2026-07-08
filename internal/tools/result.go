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
// the tool call's result.
func JSONResult(v any) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultJSON(v)
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
