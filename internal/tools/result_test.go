package tools

import (
	"errors"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"clickup-mcp/internal/clickup"
)

func textOf(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	for _, c := range res.Content {
		if tc, ok := c.(mcp.TextContent); ok {
			return tc.Text
		}
	}
	t.Fatal("no text content in result")
	return ""
}

func TestJSONResult(t *testing.T) {
	res, err := JSONResult(map[string]any{"id": "42"})
	if err != nil {
		t.Fatalf("JSONResult: %v", err)
	}
	if res.IsError {
		t.Error("IsError = true, want false")
	}
}

func TestErrorResultWithAPIError(t *testing.T) {
	apiErr := &clickup.APIError{StatusCode: 404, Err: "Task not found", ECode: "TASK_001"}
	res, err := ErrorResult(apiErr)
	if err != nil {
		t.Fatalf("ErrorResult returned non-nil error: %v", err)
	}
	if !res.IsError {
		t.Error("IsError = false, want true")
	}
	text := textOf(t, res)
	if !strings.Contains(text, "404") || !strings.Contains(text, "TASK_001") || !strings.Contains(text, "Task not found") {
		t.Errorf("error text = %q, want it to mention status/code/message", text)
	}
}

func TestErrorResultWithGenericError(t *testing.T) {
	res, err := ErrorResult(errors.New("boom"))
	if err != nil {
		t.Fatalf("ErrorResult returned non-nil error: %v", err)
	}
	if !res.IsError {
		t.Error("IsError = false, want true")
	}
	if textOf(t, res) != "boom" {
		t.Errorf("error text = %q, want %q", textOf(t, res), "boom")
	}
}
