package tools

import (
	"encoding/json"
	"runtime"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupGetServerVersion(t *testing.T) {
	t.Run("reports ServerVersion and the Go runtime version", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterVersionTools(s, c)

		old := ServerVersion
		ServerVersion = "v1.2.3"
		t.Cleanup(func() { ServerVersion = old })

		res := callTool(t, s, "clickup_get_server_version", map[string]any{})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}

		var out map[string]any
		if err := json.Unmarshal([]byte(textOf(t, res)), &out); err != nil {
			t.Fatalf("decoding result: %v", err)
		}
		if out["version"] != "v1.2.3" {
			t.Errorf("version = %v, want v1.2.3", out["version"])
		}
		if out["go_version"] != runtime.Version() {
			t.Errorf("go_version = %v, want %v", out["go_version"], runtime.Version())
		}
	})
}
