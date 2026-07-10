// Command clickup-mcp is an MCP server that exposes the ClickUp REST API as
// tools over stdio.
package main

import (
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
	"clickup-mcp/internal/config"
	"clickup-mcp/internal/tools"
)

// version is overridden at build time via -ldflags "-X main.version=vX.Y.Z"
// (see .github/workflows/release.yml); "dev" for local, unreleased builds.
var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("clickup-mcp " + version)
		return
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "clickup-mcp: "+err.Error())
		os.Exit(1)
	}

	client := clickup.NewClient(cfg)

	s := server.NewMCPServer("clickup-mcp", version,
		server.WithInstructions(
			"All date/time fields, both in tool call arguments and in "+
				"responses (start, end, date_created, due_date, etc.), are "+
				"plain UTC strings — \"YYYY-MM-DD HH:MM:SS\" for a precise "+
				"moment, or bare \"YYYY-MM-DD\" where only a calendar date "+
				"applies (e.g. a task's due_date/start_date) — never raw Unix "+
				"epoch milliseconds. duration/duration_ms fields are the one "+
				"exception: they're an elapsed length, not a point in time, "+
				"and stay in milliseconds. Convert UTC to the user's local "+
				"timezone before reporting a calendar date, since a "+
				"UTC-vs-local mismatch can shift entries near midnight onto "+
				"the wrong day. "+
				"clickup_list_time_entries silently defaults to the last 30 "+
				"days AND only entries created by the token's own user when "+
				"start_date/end_date/assignee are omitted; ClickUp returns "+
				"this truncated result as a normal 200 OK with no warning. "+
				"For historical or multi-user reports, always pass explicit "+
				"start_date and end_date (and assignee, if reporting on "+
				"someone else — it accepts a comma-separated list of user "+
				"IDs to cover multiple people in one call), or use the "+
				"task_id filter to scope to one task, or use "+
				"clickup_get_list_time_report/clickup_get_user_time_report for "+
				"a ready-made aggregated report instead of assembling one by hand. "+
				"Only clickup_get_task supports ClickUp's custom task IDs (e.g. "+
				"\"CT-123\", via its custom_task_ids/team_id params) — every other "+
				"task-scoped tool (update, delete, comments, checklists, custom "+
				"fields, dependencies, links, list membership) only accepts the "+
				"internal task ID and will 404 on a custom one with no hint why; "+
				"resolve a custom ID via clickup_get_task first. "+
				"clickup_update_doc_page defaults content_edit_mode to \"replace\" "+
				"when omitted, silently overwriting the entire page instead of "+
				"appending — always pass content_edit_mode explicitly when adding "+
				"to existing content. Guest management and workspace-user-admin "+
				"tools (invite/update/remove) are Enterprise-plan-only at the API "+
				"level and return an expected 4xx on other plans.",
		),
	)
	tools.RegisterAll(s, client)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintln(os.Stderr, "clickup-mcp: server error: "+err.Error())
		os.Exit(1)
	}
}
