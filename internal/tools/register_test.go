package tools

import (
	"testing"

	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

// registerFuncs lists every Register*Tools function invoked (directly or
// transitively) by RegisterAll. Kept in sync manually with register.go; a
// missing entry here would make the cross-function overlap check in
// TestRegisterAllHasNoDuplicateNames pass trivially rather than catch a
// regression, so if RegisterAll gains a new resource area, add it here too.
var registerFuncs = []func(*server.MCPServer, *clickup.Client){
	RegisterAuthTools,
	RegisterSpaceTools,
	RegisterFolderTools,
	RegisterListTools,
	RegisterTaskTools,
	RegisterCommentTools,
	RegisterChecklistTools,
	RegisterCustomFieldTools,
	RegisterDependencyTools,
	RegisterTagTools,
	RegisterGoalTools,
	RegisterTimeTrackingTools,
	RegisterViewTools,
	RegisterReportTools,
	RegisterWebhookTools,
	RegisterGuestTools,
	RegisterMemberTools,
	RegisterUserTools,
	RegisterRoleTools,
	RegisterTemplateTools,
	RegisterAttachmentTools,
	RegisterAuditLogTools,
	RegisterChatTools,
	RegisterCustomTaskTypeTools,
	RegisterSharedHierarchyTools,
	RegisterDocTools,
}

// expectedToolCount is the total number of distinct MCP tools RegisterAll
// registers, as of this test's last update. It's a deliberate regression
// guard: if a future change accidentally drops (or silently collides) a
// tool registration, this baseline catches it. Update it deliberately
// whenever tools are intentionally added or removed.
const expectedToolCount = 144

func TestRegisterAllDoesNotPanic(t *testing.T) {
	s := server.NewMCPServer("test", "1.0.0")
	RegisterAll(s, nil)
}

func TestRegisterAllToolCountBaseline(t *testing.T) {
	s := server.NewMCPServer("test", "1.0.0")
	RegisterAll(s, nil)
	got := len(s.ListTools())
	if got != expectedToolCount {
		t.Errorf("RegisterAll registered %d tools, want %d (update expectedToolCount if this change is intentional)", got, expectedToolCount)
	}
}

// TestRegisterAllHasNoDuplicateNames guards against two different
// Register*Tools functions accidentally registering the same tool name
// (which s.AddTool would silently allow, overwriting the earlier one with
// no error). It does this by registering each resource area on its own
// fresh server, summing their individual tool counts, and checking that sum
// against the combined count from a single RegisterAll call: if any name
// collided across functions, the combined count would be lower than the sum.
func TestRegisterAllHasNoDuplicateNames(t *testing.T) {
	var sum int
	for _, fn := range registerFuncs {
		s := server.NewMCPServer("test", "1.0.0")
		fn(s, nil)
		sum += len(s.ListTools())
	}

	combined := server.NewMCPServer("test", "1.0.0")
	RegisterAll(combined, nil)
	got := len(combined.ListTools())

	if got != sum {
		t.Errorf("combined tool count = %d, want %d (sum of each resource area registered independently) — a tool name likely collides across two files", got, sum)
	}
}
