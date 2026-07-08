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

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "clickup-mcp: "+err.Error())
		os.Exit(1)
	}

	client := clickup.NewClient(cfg)

	s := server.NewMCPServer("clickup-mcp", "0.1.0",
		server.WithInstructions(
			"All timestamp fields returned by ClickUp (start, end, duration, "+
				"at, date_created, due_date, etc.) are Unix epoch milliseconds "+
				"in UTC, not seconds and not local time — convert to the "+
				"user's timezone before reporting a calendar date. "+
				"clickup_list_time_entries silently defaults to the last 30 "+
				"days AND only entries created by the token's own user when "+
				"start_date/end_date/assignee are omitted; ClickUp returns "+
				"this truncated result as a normal 200 OK with no warning. "+
				"For historical or multi-user reports, always pass explicit "+
				"start_date and end_date (and assignee, if reporting on "+
				"someone else), or use the task_id filter to scope to one task.",
		),
	)
	tools.RegisterAll(s, client)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintln(os.Stderr, "clickup-mcp: server error: "+err.Error())
		os.Exit(1)
	}
}
