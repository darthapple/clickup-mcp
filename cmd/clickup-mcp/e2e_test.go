//go:build e2e

// Real end-to-end tests: these spawn the actual compiled binary over real
// stdio and call its tools against the real ClickUp API — the only test
// tier in this repo that can tell us whether our assumptions about ClickUp's
// actual behavior (not just our own code's internal wiring) are correct.
//
// Unlike smoke_test.go (read-only, safe against any workspace), this suite
// creates/mutates/deletes real data, so it needs a dedicated sandbox: it
// provisions its own disposable List inside CLICKUP_E2E_SPACE_ID at the
// start of the run and deletes it (cascading to everything created inside
// it) at the end.
//
// Run with:
//
//	go build -o bin/clickup-mcp ./cmd/clickup-mcp
//	set -a; source ../../.env; set +a
//	go test -tags e2e ./cmd/clickup-mcp/... -run TestE2E -v
//
// Requires CLICKUP_API_TOKEN, CLICKUP_TEAM_ID, and CLICKUP_E2E_SPACE_ID.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	mcpclient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"

	"clickup-mcp/internal/clickup"
	"clickup-mcp/internal/config"
)

// e2eListID is the disposable sandbox List provisioned by TestMain for the
// duration of this run; every test creates its fixtures inside it.
var e2eListID string

func TestMain(m *testing.M) {
	token := os.Getenv("CLICKUP_API_TOKEN")
	teamID := os.Getenv("CLICKUP_TEAM_ID")
	spaceID := os.Getenv("CLICKUP_E2E_SPACE_ID")
	if token == "" || teamID == "" || spaceID == "" {
		fmt.Println("skipping e2e suite: CLICKUP_API_TOKEN, CLICKUP_TEAM_ID, and CLICKUP_E2E_SPACE_ID must all be set")
		os.Exit(0)
	}

	c := clickup.NewClient(&config.Config{
		APIToken:    token,
		TeamID:      teamID,
		BaseURLv2:   "https://api.clickup.com/api/v2",
		BaseURLv3:   "https://api.clickup.com/api/v3",
		HTTPTimeout: 30 * time.Second,
		MaxRetries:  4,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	listName := fmt.Sprintf("clickup-mcp-e2e-%d", time.Now().UnixNano())
	raw, err := c.CreateFolderlessList(ctx, spaceID, map[string]any{"name": listName})
	if err != nil {
		fmt.Println("e2e setup: creating sandbox list:", err)
		os.Exit(1)
	}
	listMap, _ := raw.(map[string]any)
	e2eListID, _ = listMap["id"].(string)
	if e2eListID == "" {
		fmt.Printf("e2e setup: could not read id from CreateFolderlessList response: %+v\n", raw)
		os.Exit(1)
	}

	code := m.Run()

	teardownCtx, teardownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer teardownCancel()
	if err := c.DeleteList(teardownCtx, e2eListID); err != nil {
		fmt.Printf("e2e teardown: deleting sandbox list %s: %v\n", e2eListID, err)
	}

	os.Exit(code)
}

// newE2EClient spawns the compiled binary as a real MCP stdio subprocess
// against the real ClickUp API.
func newE2EClient(t *testing.T) *mcpclient.Client {
	t.Helper()
	ctx := context.Background()

	c, err := mcpclient.NewStdioMCPClient("../../bin/clickup-mcp", os.Environ())
	if err != nil {
		t.Fatalf("start server: %v", err)
	}
	t.Cleanup(func() { c.Close() })

	if _, err := c.Initialize(ctx, mcp.InitializeRequest{}); err != nil {
		t.Fatalf("initialize: %v", err)
	}
	return c
}

// callE2ETool calls name against the real API and decodes its text content
// as JSON. Fails the test if the call errors or returns an MCP error result.
func callE2ETool(ctx context.Context, t *testing.T, c *mcpclient.Client, name string, args map[string]any) map[string]any {
	t.Helper()
	req := mcp.CallToolRequest{}
	req.Params.Name = name
	req.Params.Arguments = args

	res, err := c.CallTool(ctx, req)
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	if res.IsError {
		t.Fatalf("%s returned an error result: %s", name, textContent(res))
	}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(textContent(res)), &decoded); err != nil {
		t.Fatalf("%s: decoding result: %v", name, err)
	}
	return decoded
}

// callE2EToolExpectError calls name and asserts it returns an MCP error
// result (rather than succeeding), returning the error text.
func callE2EToolExpectError(ctx context.Context, t *testing.T, c *mcpclient.Client, name string, args map[string]any) string {
	t.Helper()
	req := mcp.CallToolRequest{}
	req.Params.Name = name
	req.Params.Arguments = args

	res, err := c.CallTool(ctx, req)
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	if !res.IsError {
		t.Fatalf("%s: expected an error result, got success: %s", name, textContent(res))
	}
	return textContent(res)
}

func textContent(res *mcp.CallToolResult) string {
	for _, c := range res.Content {
		if tc, ok := c.(mcp.TextContent); ok {
			return tc.Text
		}
	}
	return ""
}

// TestE2ETaskLifecycle exercises create -> get -> update -> delete -> get
// against the real ClickUp API: we choose the input, so we know the
// expected output, and we assert it against ClickUp's actual response
// instead of a fixture we wrote ourselves.
func TestE2ETaskLifecycle(t *testing.T) {
	ctx := context.Background()
	c := newE2EClient(t)

	created := callE2ETool(ctx, t, c, "clickup_create_task", map[string]any{
		"list_id":     e2eListID,
		"name":        "e2e task lifecycle",
		"description": "created by TestE2ETaskLifecycle",
		"priority":    float64(3),
	})
	taskID, _ := created["id"].(string)
	if taskID == "" {
		t.Fatalf("create response has no id: %+v", created)
	}
	t.Cleanup(func() {
		req := mcp.CallToolRequest{}
		req.Params.Name = "clickup_delete_task"
		req.Params.Arguments = map[string]any{"task_id": taskID}
		_, _ = c.CallTool(context.Background(), req)
	})

	if created["name"] != "e2e task lifecycle" {
		t.Errorf("created name = %v, want %q", created["name"], "e2e task lifecycle")
	}

	got := callE2ETool(ctx, t, c, "clickup_get_task", map[string]any{"task_id": taskID})
	if got["name"] != "e2e task lifecycle" {
		t.Errorf("get name = %v, want %q", got["name"], "e2e task lifecycle")
	}
	if got["description"] != "created by TestE2ETaskLifecycle" {
		t.Errorf("get description = %v", got["description"])
	}

	updated := callE2ETool(ctx, t, c, "clickup_update_task", map[string]any{
		"task_id": taskID,
		"name":    "e2e task lifecycle (updated)",
	})
	if updated["name"] != "e2e task lifecycle (updated)" {
		t.Errorf("update response name = %v", updated["name"])
	}

	gotAfterUpdate := callE2ETool(ctx, t, c, "clickup_get_task", map[string]any{"task_id": taskID})
	if gotAfterUpdate["name"] != "e2e task lifecycle (updated)" {
		t.Errorf("get after update name = %v, want the updated name", gotAfterUpdate["name"])
	}

	callE2ETool(ctx, t, c, "clickup_delete_task", map[string]any{"task_id": taskID})

	errText := callE2EToolExpectError(ctx, t, c, "clickup_get_task", map[string]any{"task_id": taskID})
	t.Logf("get after delete correctly errored: %s", errText)
}

// TestE2ECommentLifecycle exercises create -> list -> delete for task
// comments against the real API.
func TestE2ECommentLifecycle(t *testing.T) {
	ctx := context.Background()
	c := newE2EClient(t)

	task := callE2ETool(ctx, t, c, "clickup_create_task", map[string]any{
		"list_id": e2eListID,
		"name":    "e2e comment lifecycle",
	})
	taskID, _ := task["id"].(string)
	if taskID == "" {
		t.Fatalf("create task response has no id: %+v", task)
	}
	t.Cleanup(func() {
		req := mcp.CallToolRequest{}
		req.Params.Name = "clickup_delete_task"
		req.Params.Arguments = map[string]any{"task_id": taskID}
		_, _ = c.CallTool(context.Background(), req)
	})

	commentText := "e2e comment lifecycle test comment"
	callE2ETool(ctx, t, c, "clickup_create_task_comment", map[string]any{
		"task_id":      taskID,
		"comment_text": commentText,
	})

	listed := callE2ETool(ctx, t, c, "clickup_list_task_comments", map[string]any{"task_id": taskID})
	comments, _ := listed["comments"].([]any)
	if len(comments) == 0 {
		t.Fatalf("expected at least one comment, got none: %+v", listed)
	}
	first, _ := comments[0].(map[string]any)
	commentID, _ := first["id"].(string)
	if commentID == "" {
		t.Fatalf("comment has no id: %+v", first)
	}
	textField, _ := first["comment_text"].(string)
	if textField != commentText {
		t.Errorf("comment_text = %q, want %q", textField, commentText)
	}

	callE2ETool(ctx, t, c, "clickup_delete_comment", map[string]any{"comment_id": commentID})

	listedAfterDelete := callE2ETool(ctx, t, c, "clickup_list_task_comments", map[string]any{"task_id": taskID})
	remaining, _ := listedAfterDelete["comments"].([]any)
	for _, rc := range remaining {
		rcm, _ := rc.(map[string]any)
		if rcm["id"] == commentID {
			t.Errorf("comment %s still present after delete", commentID)
		}
	}
}

// TestE2EChecklistLifecycle exercises create checklist -> create item ->
// resolve item -> delete checklist against the real API, asserting the
// resulting state through clickup_get_task.
func TestE2EChecklistLifecycle(t *testing.T) {
	ctx := context.Background()
	c := newE2EClient(t)

	task := callE2ETool(ctx, t, c, "clickup_create_task", map[string]any{
		"list_id": e2eListID,
		"name":    "e2e checklist lifecycle",
	})
	taskID, _ := task["id"].(string)
	if taskID == "" {
		t.Fatalf("create task response has no id: %+v", task)
	}
	t.Cleanup(func() {
		req := mcp.CallToolRequest{}
		req.Params.Name = "clickup_delete_task"
		req.Params.Arguments = map[string]any{"task_id": taskID}
		_, _ = c.CallTool(context.Background(), req)
	})

	checklist := callE2ETool(ctx, t, c, "clickup_create_checklist", map[string]any{
		"task_id": taskID,
		"name":    "e2e checklist",
	})
	checklistMap, _ := checklist["checklist"].(map[string]any)
	checklistID, _ := checklistMap["id"].(string)
	if checklistID == "" {
		t.Fatalf("checklist response has no id: %+v", checklist)
	}

	item := callE2ETool(ctx, t, c, "clickup_create_checklist_item", map[string]any{
		"checklist_id": checklistID,
		"name":         "e2e checklist item",
	})
	itemChecklist, _ := item["checklist"].(map[string]any)
	items, _ := itemChecklist["items"].([]any)
	if len(items) == 0 {
		t.Fatalf("expected at least one checklist item, got none: %+v", item)
	}
	firstItem, _ := items[0].(map[string]any)
	itemID, _ := firstItem["id"].(string)
	if itemID == "" {
		t.Fatalf("checklist item has no id: %+v", firstItem)
	}

	callE2ETool(ctx, t, c, "clickup_update_checklist_item", map[string]any{
		"checklist_id":      checklistID,
		"checklist_item_id": itemID,
		"resolved":          true,
	})

	got := callE2ETool(ctx, t, c, "clickup_get_task", map[string]any{"task_id": taskID})
	checklists, _ := got["checklists"].([]any)
	found := false
	for _, cl := range checklists {
		clm, _ := cl.(map[string]any)
		if clm["id"] != checklistID {
			continue
		}
		found = true
		clItems, _ := clm["items"].([]any)
		for _, ci := range clItems {
			cim, _ := ci.(map[string]any)
			if cim["id"] == itemID && cim["resolved"] != true {
				t.Errorf("checklist item resolved = %v, want true", cim["resolved"])
			}
		}
	}
	if !found {
		t.Errorf("checklist %s not found on task: %+v", checklistID, got)
	}

	callE2ETool(ctx, t, c, "clickup_delete_checklist", map[string]any{"checklist_id": checklistID})
}

// TestE2ETimeTrackingLifecycle exercises clickup_create_time_entry ->
// clickup_list_time_entries -> clickup_get_list_time_report ->
// clickup_delete_time_entry against the real API. clickup_get_list_time_report
// is a cross-cutting report tool (internal/tools/report_tools.go) that
// aggregates GetTasksInList + ListTimeEntries itself, so this is the one
// real-API check that its aggregation/rollup math (not just each underlying
// endpoint call) matches ClickUp's actual data.
func TestE2ETimeTrackingLifecycle(t *testing.T) {
	ctx := context.Background()
	c := newE2EClient(t)

	task := callE2ETool(ctx, t, c, "clickup_create_task", map[string]any{
		"list_id": e2eListID,
		"name":    "e2e time tracking lifecycle",
	})
	taskID, _ := task["id"].(string)
	if taskID == "" {
		t.Fatalf("create task response has no id: %+v", task)
	}
	t.Cleanup(func() {
		req := mcp.CallToolRequest{}
		req.Params.Name = "clickup_delete_task"
		req.Params.Arguments = map[string]any{"task_id": taskID}
		_, _ = c.CallTool(context.Background(), req)
	})

	const durationMs = float64(3600000) // 1 hour
	now := time.Now()
	start := now.Add(-2 * time.Hour).UTC().Format("2006-01-02 15:04:05")
	// Report window wide enough to comfortably contain the entry regardless
	// of clock skew between this test and ClickUp's server.
	reportStart := now.Add(-24 * time.Hour).UTC().Format("2006-01-02 15:04:05")
	reportEnd := now.Add(24 * time.Hour).UTC().Format("2006-01-02 15:04:05")

	entry := callE2ETool(ctx, t, c, "clickup_create_time_entry", map[string]any{
		"task_id":     taskID,
		"start":       start,
		"duration":    durationMs,
		"description": "e2e time tracking entry",
	})
	entryData, _ := entry["data"].(map[string]any)
	entryID, _ := entryData["id"].(string)
	if entryID == "" {
		t.Fatalf("create time entry response has no data.id: %+v", entry)
	}
	t.Cleanup(func() {
		req := mcp.CallToolRequest{}
		req.Params.Name = "clickup_delete_time_entry"
		req.Params.Arguments = map[string]any{"timer_id": entryID}
		_, _ = c.CallTool(context.Background(), req)
	})

	listed := callE2ETool(ctx, t, c, "clickup_list_time_entries", map[string]any{
		"task_id":    taskID,
		"start_date": reportStart,
		"end_date":   reportEnd,
	})
	entries, _ := listed["data"].([]any)
	found := false
	for _, e := range entries {
		em, _ := e.(map[string]any)
		if em["id"] == entryID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created time entry %s not found in clickup_list_time_entries: %+v", entryID, listed)
	}

	report := callE2ETool(ctx, t, c, "clickup_get_list_time_report", map[string]any{
		"list_id":    e2eListID,
		"start_date": reportStart,
		"end_date":   reportEnd,
	})
	if report["total_duration_ms"] != durationMs {
		t.Errorf("report total_duration_ms = %v, want %v", report["total_duration_ms"], durationMs)
	}
	tasks, _ := report["tasks"].([]any)
	foundInReport := false
	for _, tr := range tasks {
		trm, _ := tr.(map[string]any)
		if trm["task_id"] != taskID {
			continue
		}
		foundInReport = true
		if trm["duration_ms"] != durationMs {
			t.Errorf("report task duration_ms = %v, want %v", trm["duration_ms"], durationMs)
		}
	}
	if !foundInReport {
		t.Errorf("task %s not found in clickup_get_list_time_report output: %+v", taskID, report)
	}

	// Some ClickUp plans reject DELETE on time entries ("Time Tracking is not
	// available on your plan") even though create/list/report all succeed on
	// the same workspace — a real account-tier restriction, not a bug in our
	// code. Tolerate it here (the sandbox List's cascading delete in TestMain
	// still removes the entry along with its task) rather than failing the
	// whole suite on accounts without that entitlement.
	deleteReq := mcp.CallToolRequest{}
	deleteReq.Params.Name = "clickup_delete_time_entry"
	deleteReq.Params.Arguments = map[string]any{"timer_id": entryID}
	deleteRes, err := c.CallTool(ctx, deleteReq)
	if err != nil {
		t.Fatalf("clickup_delete_time_entry: %v", err)
	}
	if deleteRes.IsError {
		errText := textContent(deleteRes)
		if strings.Contains(errText, "not available on your plan") {
			t.Logf("clickup_delete_time_entry blocked by account plan (expected on this workspace), skipping post-delete verification: %s", errText)
			return
		}
		t.Fatalf("clickup_delete_time_entry returned an unexpected error: %s", errText)
	}

	listedAfterDelete := callE2ETool(ctx, t, c, "clickup_list_time_entries", map[string]any{
		"task_id":    taskID,
		"start_date": reportStart,
		"end_date":   reportEnd,
	})
	remaining, _ := listedAfterDelete["data"].([]any)
	for _, e := range remaining {
		em, _ := e.(map[string]any)
		if em["id"] == entryID {
			t.Errorf("time entry %s still present after delete", entryID)
		}
	}
}
