package tools

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupInviteGuest(t *testing.T) {
	t.Run("requires email", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGuestTools(s, c)

		res := callTool(t, s, "clickup_invite_guest", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires body and method/path", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"guest":{"id":"g1"}}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGuestTools(s, c)

		res := callTool(t, s, "clickup_invite_guest", map[string]any{
			"email":              "guest@example.com",
			"can_edit_tags":      true,
			"can_see_time_spent": false,
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %s, want POST", gotMethod)
		}
		if gotPath != "/team/999/guest" {
			t.Errorf("path = %q, want /team/999/guest", gotPath)
		}
		if gotBody["email"] != "guest@example.com" {
			t.Errorf("body[email] = %v", gotBody["email"])
		}
		if gotBody["can_edit_tags"] != true {
			t.Errorf("body[can_edit_tags] = %v, want true", gotBody["can_edit_tags"])
		}
		if gotBody["can_see_time_spent"] != false {
			t.Errorf("body[can_see_time_spent] = %v, want false", gotBody["can_see_time_spent"])
		}
		if _, present := gotBody["can_create_views"]; present {
			t.Errorf("body[can_create_views] present = true, want absent (not supplied)")
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"err":"bad","ECODE":"X_002"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGuestTools(s, c)

		res := callTool(t, s, "clickup_invite_guest", map[string]any{"email": "guest@example.com"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})
}

func TestClickupGetGuest(t *testing.T) {
	t.Run("requires guest_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGuestTools(s, c)

		res := callTool(t, s, "clickup_get_guest", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires guest_id into path", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"guest":{"id":"g1"}}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGuestTools(s, c)

		res := callTool(t, s, "clickup_get_guest", map[string]any{"guest_id": "g1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %s, want GET", gotMethod)
		}
		if gotPath != "/team/999/guest/g1" {
			t.Errorf("path = %q, want /team/999/guest/g1", gotPath)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGuestTools(s, c)

		res := callTool(t, s, "clickup_get_guest", map[string]any{"guest_id": "missing"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupUpdateGuest(t *testing.T) {
	t.Run("requires guest_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGuestTools(s, c)

		res := callTool(t, s, "clickup_update_guest", map[string]any{"can_edit_tags": true})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("partial update sends only supplied field", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"guest":{"id":"g1"}}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGuestTools(s, c)

		res := callTool(t, s, "clickup_update_guest", map[string]any{
			"guest_id":         "g1",
			"can_create_views": true,
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodPut {
			t.Errorf("method = %s, want PUT", gotMethod)
		}
		if gotPath != "/team/999/guest/g1" {
			t.Errorf("path = %q, want /team/999/guest/g1", gotPath)
		}
		if len(gotBody) != 1 || gotBody["can_create_views"] != true {
			t.Errorf("body = %+v, want only can_create_views=true", gotBody)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGuestTools(s, c)

		res := callTool(t, s, "clickup_update_guest", map[string]any{"guest_id": "g1", "can_edit_tags": true})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if got := textOf(t, res); got != want {
			t.Errorf("error text = %q, want %q", got, want)
		}
	})
}

func TestClickupRemoveGuestFromWorkspace(t *testing.T) {
	t.Run("requires guest_id", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGuestTools(s, c)

		res := callTool(t, s, "clickup_remove_guest_from_workspace", map[string]any{})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
	})

	t.Run("wires guest_id and reports removed", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGuestTools(s, c)

		res := callTool(t, s, "clickup_remove_guest_from_workspace", map[string]any{"guest_id": "g1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", gotMethod)
		}
		if gotPath != "/team/999/guest/g1" {
			t.Errorf("path = %q, want /team/999/guest/g1", gotPath)
		}
		var out map[string]any
		if err := json.Unmarshal([]byte(textOf(t, res)), &out); err != nil {
			t.Fatalf("decoding result: %v", err)
		}
		if out["removed"] != true || out["guest_id"] != "g1" {
			t.Errorf("result = %+v", out)
		}
	})
}

// TestClickupGuestScopes guards against a copy-paste/loop bug in
// registerGuestScope silently colliding two scopes' tool names or request
// paths: it iterates all 4 scopes (space, folder, list, task) and verifies
// both the add and remove tool are registered under distinct, scope-correct
// names, and that calling each one hits the fake server at the
// scope-correct path with the scope-correct HTTP method.
func TestClickupGuestScopes(t *testing.T) {
	scopes := []string{"space", "folder", "list", "task"}

	for _, scope := range scopes {
		scope := scope
		t.Run(scope, func(t *testing.T) {
			idParam := scope + "_id"
			addTool := "clickup_add_guest_to_" + scope
			removeTool := "clickup_remove_guest_from_" + scope
			resourceID := scope + "1"

			t.Run("add: registered and requires "+idParam, func(t *testing.T) {
				c, _ := newTestClient(t, unreachableHandler(t))
				s := server.NewMCPServer("test", "1.0.0")
				RegisterGuestTools(s, c)

				if s.GetTool(addTool) == nil {
					t.Fatalf("tool %q is not registered", addTool)
				}

				res := callTool(t, s, addTool, map[string]any{"guest_id": "g1"})
				if !res.IsError {
					t.Fatalf("IsError = false, want true (missing %s)", idParam)
				}
			})

			t.Run("add: wires scope-correct path", func(t *testing.T) {
				var gotMethod, gotPath string
				var gotBody map[string]any
				c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
					gotMethod = r.Method
					gotPath = r.URL.Path
					_ = json.NewDecoder(r.Body).Decode(&gotBody)
					_, _ = w.Write([]byte(`{"id":"g1"}`))
				})
				s := server.NewMCPServer("test", "1.0.0")
				RegisterGuestTools(s, c)

				res := callTool(t, s, addTool, map[string]any{
					idParam:            resourceID,
					"guest_id":         "g1",
					"permission_level": float64(2),
				})
				if res.IsError {
					t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
				}
				if gotMethod != http.MethodPost {
					t.Errorf("method = %s, want POST", gotMethod)
				}
				wantPath := "/" + scope + "/" + resourceID + "/guest/g1"
				if gotPath != wantPath {
					t.Errorf("path = %q, want %q", gotPath, wantPath)
				}
				if gotBody["permission_level"] != float64(2) {
					t.Errorf("body[permission_level] = %v, want 2", gotBody["permission_level"])
				}
			})

			t.Run("remove: registered and requires "+idParam, func(t *testing.T) {
				c, _ := newTestClient(t, unreachableHandler(t))
				s := server.NewMCPServer("test", "1.0.0")
				RegisterGuestTools(s, c)

				if s.GetTool(removeTool) == nil {
					t.Fatalf("tool %q is not registered", removeTool)
				}

				res := callTool(t, s, removeTool, map[string]any{"guest_id": "g1"})
				if !res.IsError {
					t.Fatalf("IsError = false, want true (missing %s)", idParam)
				}
			})

			t.Run("remove: wires scope-correct path", func(t *testing.T) {
				var gotMethod, gotPath string
				c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
					gotMethod = r.Method
					gotPath = r.URL.Path
					w.WriteHeader(http.StatusNoContent)
				})
				s := server.NewMCPServer("test", "1.0.0")
				RegisterGuestTools(s, c)

				res := callTool(t, s, removeTool, map[string]any{
					idParam:    resourceID,
					"guest_id": "g1",
				})
				if res.IsError {
					t.Fatalf("IsError = true, want false; text = %q", textOf(t, res))
				}
				if gotMethod != http.MethodDelete {
					t.Errorf("method = %s, want DELETE", gotMethod)
				}
				wantPath := "/" + scope + "/" + resourceID + "/guest/g1"
				if gotPath != wantPath {
					t.Errorf("path = %q, want %q", gotPath, wantPath)
				}
				var out map[string]any
				if err := json.Unmarshal([]byte(textOf(t, res)), &out); err != nil {
					t.Fatalf("decoding result: %v", err)
				}
				if out["removed"] != true || out[idParam] != resourceID || out["guest_id"] != "g1" {
					t.Errorf("result = %+v", out)
				}
			})
		})
	}

	// Cross-scope sanity check: all 8 generated tool names must be distinct
	// (guards against e.g. two scopes both registering "clickup_add_guest_to_list").
	t.Run("all scope tool names are distinct and registered", func(t *testing.T) {
		c, _ := newTestClient(t, unreachableHandler(t))
		s := server.NewMCPServer("test", "1.0.0")
		RegisterGuestTools(s, c)

		seen := map[string]bool{}
		for _, scope := range scopes {
			for _, name := range []string{"clickup_add_guest_to_" + scope, "clickup_remove_guest_from_" + scope} {
				if seen[name] {
					t.Fatalf("tool name %q registered more than once across scopes", name)
				}
				seen[name] = true
				if s.GetTool(name) == nil {
					t.Fatalf("tool %q is not registered", name)
				}
			}
		}
		if len(seen) != 8 {
			t.Fatalf("got %d distinct scope tool names, want 8", len(seen))
		}
	})
}
