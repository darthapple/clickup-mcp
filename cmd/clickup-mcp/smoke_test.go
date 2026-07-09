//go:build smoke

// Smoke test: a fast, read-only sanity check against the real ClickUp API —
// "is the deployed binary able to authenticate and talk to ClickUp at all,"
// not a thorough correctness check (see e2e_test.go, build tag e2e, for
// that). Safe to run against any real workspace with zero setup beyond
// credentials, since it never creates/updates/deletes anything. Excluded
// from `go test ./...`; run explicitly from the module root with:
//
//	go build -o bin/clickup-mcp ./cmd/clickup-mcp
//	go test -tags smoke ./cmd/clickup-mcp/... -run TestSmoke -v
//
// Requires CLICKUP_API_TOKEN and CLICKUP_TEAM_ID in the environment.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestSmoke(t *testing.T) {
	if os.Getenv("CLICKUP_API_TOKEN") == "" || os.Getenv("CLICKUP_TEAM_ID") == "" {
		t.Skip("CLICKUP_API_TOKEN / CLICKUP_TEAM_ID not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c, err := client.NewStdioMCPClient("../../bin/clickup-mcp", os.Environ())
	if err != nil {
		t.Fatalf("start server: %v", err)
	}
	defer c.Close()

	if _, err := c.Initialize(ctx, mcp.InitializeRequest{}); err != nil {
		t.Fatalf("initialize: %v", err)
	}

	toolsRes, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	t.Logf("server exposes %d tools", len(toolsRes.Tools))
	for _, tool := range toolsRes.Tools {
		t.Logf("  - %s", tool.Name)
	}

	teamID := os.Getenv("CLICKUP_TEAM_ID")

	user := callJSON(ctx, t, c, "clickup_get_user", nil)
	userObj, _ := user["user"].(map[string]any)
	// GetUser's id comes back as a raw JSON number, unlike the string IDs
	// used everywhere else in ClickUp's API.
	if id, ok := userObj["id"].(float64); !ok || id == 0 {
		t.Errorf("clickup_get_user: user.id = %v, want a non-zero number", userObj["id"])
	}
	if _, ok := userObj["username"].(string); !ok {
		t.Errorf("clickup_get_user: user.username missing or not a string: %+v", userObj)
	}

	workspaces := callJSON(ctx, t, c, "clickup_list_workspaces", nil)
	teams, _ := workspaces["teams"].([]any)
	foundTeam := false
	for _, tm := range teams {
		tmMap, _ := tm.(map[string]any)
		if tmMap["id"] == teamID {
			foundTeam = true
			break
		}
	}
	if !foundTeam {
		t.Errorf("clickup_list_workspaces: CLICKUP_TEAM_ID %q not found among accessible teams: %+v", teamID, teams)
	}

	// Walk the real hierarchy read-only (no create/update/delete here, since
	// this runs against the user's live workspace) to exercise Phase 1+2
	// list/get endpoints end to end.
	spaces := callJSON(ctx, t, c, "clickup_list_spaces", map[string]any{"team_id": teamID})
	spaceID := firstID(spaces, "spaces")
	if spaceID == "" {
		t.Log("no spaces found, skipping folder/list/task chain")
		return
	}

	folders := callJSON(ctx, t, c, "clickup_list_folders", map[string]any{"space_id": spaceID})
	lists := callJSON(ctx, t, c, "clickup_list_folderless_lists", map[string]any{"space_id": spaceID})

	listID := firstID(lists, "lists")
	if listID == "" {
		if folderID := firstID(folders, "folders"); folderID != "" {
			folderLists := callJSON(ctx, t, c, "clickup_list_lists", map[string]any{"folder_id": folderID})
			listID = firstID(folderLists, "lists")
		}
	}
	if listID == "" {
		t.Log("no lists found, skipping task chain")
		return
	}

	callAndLog(ctx, t, c, "clickup_list_list_fields", map[string]any{"list_id": listID})
	tasks := callJSON(ctx, t, c, "clickup_list_tasks", map[string]any{"list_id": listID})
	taskID := firstID(tasks, "tasks")
	if taskID == "" {
		t.Log("no tasks found in list, skipping task detail chain")
		return
	}

	callAndLog(ctx, t, c, "clickup_get_task", map[string]any{"task_id": taskID})
	callAndLog(ctx, t, c, "clickup_list_task_comments", map[string]any{"task_id": taskID})
}

// callJSON calls a tool and decodes its text content as JSON.
func callJSON(ctx context.Context, t *testing.T, c *client.Client, name string, args map[string]any) map[string]any {
	t.Helper()
	req := mcp.CallToolRequest{}
	req.Params.Name = name
	req.Params.Arguments = args

	res, err := c.CallTool(ctx, req)
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	if res.IsError {
		t.Errorf("%s returned an error result", name)
		return nil
	}
	for _, content := range res.Content {
		if tc, ok := content.(mcp.TextContent); ok {
			var decoded map[string]any
			if json.Unmarshal([]byte(tc.Text), &decoded) == nil {
				return decoded
			}
		}
	}
	return nil
}

// firstID returns the "id" field of the first element under key in m, if any.
func firstID(m map[string]any, key string) string {
	arr, _ := m[key].([]any)
	if len(arr) == 0 {
		return ""
	}
	first, _ := arr[0].(map[string]any)
	id, _ := first["id"].(string)
	return id
}

func callAndLog(ctx context.Context, t *testing.T, c *client.Client, name string, args map[string]any) {
	t.Helper()
	req := mcp.CallToolRequest{}
	req.Params.Name = name
	req.Params.Arguments = args

	res, err := c.CallTool(ctx, req)
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	if res.IsError {
		t.Errorf("%s returned an error result", name)
	}
	for _, content := range res.Content {
		if tc, ok := content.(mcp.TextContent); ok {
			var pretty any
			if json.Unmarshal([]byte(tc.Text), &pretty) == nil {
				b, _ := json.MarshalIndent(pretty, "", "  ")
				fmt.Printf("=== %s ===\n%s\n", name, b)
			} else {
				fmt.Printf("=== %s ===\n%s\n", name, tc.Text)
			}
		}
	}
}
