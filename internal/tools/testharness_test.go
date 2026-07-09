package tools

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
	"clickup-mcp/internal/config"
)

// newTestClient builds a *clickup.Client pointed at an httptest.Server
// running handler, mirroring clickup's own (package-internal) testClient
// helper for use from the tools package.
func newTestClient(t *testing.T, handler http.HandlerFunc) (*clickup.Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c := clickup.NewClient(&config.Config{
		APIToken:    "pk_test_token",
		TeamID:      "999",
		BaseURLv2:   srv.URL,
		BaseURLv3:   srv.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  0,
	})
	return c, srv
}

// callTool looks up name on s and invokes its handler directly (bypassing
// the JSON-RPC/stdio transport), failing the test if the tool isn't
// registered.
func callTool(t *testing.T, s *server.MCPServer, name string, args map[string]any) *mcp.CallToolResult {
	t.Helper()
	st := s.GetTool(name)
	if st == nil {
		t.Fatalf("tool %q is not registered", name)
	}
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      name,
			Arguments: args,
		},
	}
	res, err := st.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("tool %q handler returned non-nil error: %v", name, err)
	}
	return res
}

// textOf is defined in result_test.go and reused across this package's
// tool-handler tests.
