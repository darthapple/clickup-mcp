package tools

import (
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// hasArg reports whether the caller explicitly supplied key, as opposed to
// it being absent and defaulted by a Get* call. Used when building
// create/update request bodies so unset optional fields are omitted rather
// than sent as zero values (important for PUT/PATCH semantics).
func hasArg(req mcp.CallToolRequest, key string) bool {
	_, ok := req.GetArguments()[key]
	return ok
}

// setString copies key into body as a string if the caller supplied it.
func setString(body map[string]any, req mcp.CallToolRequest, key string) {
	if hasArg(req, key) {
		body[key] = req.GetString(key, "")
	}
}

// setBool copies key into body as a bool if the caller supplied it.
func setBool(body map[string]any, req mcp.CallToolRequest, key string) {
	if hasArg(req, key) {
		body[key] = req.GetBool(key, false)
	}
}

// setFloat copies key into body as a number if the caller supplied it.
func setFloat(body map[string]any, req mcp.CallToolRequest, key string) {
	if hasArg(req, key) {
		body[key] = req.GetFloat(key, 0)
	}
}

// setStringSlice copies key into body as a string array if the caller
// supplied it.
func setStringSlice(body map[string]any, req mcp.CallToolRequest, key string) {
	if hasArg(req, key) {
		body[key] = req.GetStringSlice(key, nil)
	}
}

// setRawJSON parses the caller-supplied JSON text under argKey and sets it
// on body[bodyKey]. Used for fields whose shape is fixed but complex enough
// (arrays of objects, variable-typed values) that flattening into scalar
// params isn't practical. Returns an error if argKey was supplied but isn't
// valid JSON.
func setRawJSON(body map[string]any, req mcp.CallToolRequest, argKey, bodyKey string) error {
	if !hasArg(req, argKey) {
		return nil
	}
	raw := req.GetString(argKey, "")
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return fmt.Errorf("%s must be valid JSON: %w", argKey, err)
	}
	body[bodyKey] = v
	return nil
}
