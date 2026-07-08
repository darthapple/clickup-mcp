// Package clickup is the ClickUp REST API client used by the MCP tool layer.
//
// Every ClickUp identifier (task, list, folder, space, team/workspace, view,
// comment, etc.) is modeled as a Go string end to end — in client method
// signatures and in MCP tool schemas alike. ClickUp task IDs are alphanumeric
// (e.g. "868czmkqz"), and numeric-looking IDs (space/folder/list/team) can
// exceed the safe integer range for JSON numbers decoded as float64. Treating
// all IDs as strings avoids precision loss and keeps ID handling uniform.
package clickup
