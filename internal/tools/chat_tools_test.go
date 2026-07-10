package tools

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupListChatChannels(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"channels":[{"id":"c1"}]}`))
	})
	s := server.NewMCPServer("test", "1.0.0")
	RegisterChatTools(s, c)
	res := callTool(t, s, "clickup_list_chat_channels", map[string]any{"team_id": "123"})
	if res.IsError {
		t.Fatalf("IsError = true, want false: %s", textOf(t, res))
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/workspaces/123/chat/channels" {
		t.Errorf("path = %q, want /workspaces/123/chat/channels", gotPath)
	}
}

func TestClickupCreateChatChannel(t *testing.T) {
	t.Run("required arg validation", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_create_chat_channel", map[string]any{"team_id": "123"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing name)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("argument wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"c1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_create_chat_channel", map[string]any{
			"team_id":     "123",
			"name":        "general",
			"description": "General discussion",
			"topic":       "everything",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/workspaces/123/chat/channels" {
			t.Errorf("path = %q, want /workspaces/123/chat/channels", gotPath)
		}
		if gotBody["name"] != "general" || gotBody["description"] != "General discussion" || gotBody["topic"] != "everything" {
			t.Errorf("body = %+v", gotBody)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_create_chat_channel", map[string]any{"team_id": "123", "name": "general"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupGetChatChannel(t *testing.T) {
	t.Run("required arg validation", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_get_chat_channel", map[string]any{"team_id": "123"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing channel_id)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("argument wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"id":"c1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_get_chat_channel", map[string]any{"team_id": "123", "channel_id": "c1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %q, want GET", gotMethod)
		}
		if gotPath != "/workspaces/123/chat/channels/c1" {
			t.Errorf("path = %q, want /workspaces/123/chat/channels/c1", gotPath)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_get_chat_channel", map[string]any{"team_id": "123", "channel_id": "c1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupUpdateChatChannel(t *testing.T) {
	t.Run("required arg validation", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_update_chat_channel", map[string]any{"team_id": "123"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing channel_id)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("partial update semantics", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"c1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_update_chat_channel", map[string]any{
			"team_id":    "123",
			"channel_id": "c1",
			"topic":      "new topic",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/workspaces/123/chat/channels/c1" {
			t.Errorf("path = %q, want /workspaces/123/chat/channels/c1", gotPath)
		}
		if len(gotBody) != 1 {
			t.Errorf("body = %+v, want exactly one field", gotBody)
		}
		if gotBody["topic"] != "new topic" {
			t.Errorf("body[topic] = %v, want %q", gotBody["topic"], "new topic")
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_update_chat_channel", map[string]any{"team_id": "123", "channel_id": "c1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupChatMessagesCRUD(t *testing.T) {
	t.Run("list requires channel_id", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { hit = true })
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_list_chat_messages", map[string]any{"team_id": "123"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing channel_id)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("list wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"messages":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_list_chat_messages", map[string]any{"team_id": "123", "channel_id": "c1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %q, want GET", gotMethod)
		}
		if gotPath != "/workspaces/123/chat/channels/c1/messages" {
			t.Errorf("path = %q, want /workspaces/123/chat/channels/c1/messages", gotPath)
		}
	})

	t.Run("create requires content", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { hit = true })
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_create_chat_message", map[string]any{"team_id": "123", "channel_id": "c1"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing text)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("create wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"m1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_create_chat_message", map[string]any{
			"team_id": "123", "channel_id": "c1", "content": "hello",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/workspaces/123/chat/channels/c1/messages" {
			t.Errorf("path = %q, want /workspaces/123/chat/channels/c1/messages", gotPath)
		}
		if gotBody["content"] != "hello" {
			t.Errorf("body = %+v", gotBody)
		}
		if gotBody["type"] != "message" {
			t.Errorf("body[type] = %v, want default \"message\"", gotBody["type"])
		}
	})

	t.Run("update requires message_id and content", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { hit = true })
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_update_chat_message", map[string]any{"team_id": "123"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing message_id/text)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("update wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"id":"m1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_update_chat_message", map[string]any{
			"team_id": "123", "message_id": "m1", "content": "edited",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/workspaces/123/chat/messages/m1" {
			t.Errorf("path = %q, want /workspaces/123/chat/messages/m1", gotPath)
		}
		if len(gotBody) != 1 || gotBody["content"] != "edited" {
			t.Errorf("body = %+v, want only content=edited", gotBody)
		}
	})

	t.Run("delete requires message_id", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { hit = true })
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_delete_chat_message", map[string]any{"team_id": "123"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing message_id)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("delete wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_delete_chat_message", map[string]any{"team_id": "123", "message_id": "m1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/workspaces/123/chat/messages/m1" {
			t.Errorf("path = %q, want /workspaces/123/chat/messages/m1", gotPath)
		}
	})

	t.Run("error passthrough on create", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_create_chat_message", map[string]any{"team_id": "123", "channel_id": "c1", "content": "hi"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupChatReactions(t *testing.T) {
	t.Run("list requires message_id", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { hit = true })
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_list_chat_reactions", map[string]any{"team_id": "123"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing message_id)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("list wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"reactions":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_list_chat_reactions", map[string]any{"team_id": "123", "message_id": "m1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %q, want GET", gotMethod)
		}
		if gotPath != "/workspaces/123/chat/messages/m1/reactions" {
			t.Errorf("path = %q, want /workspaces/123/chat/messages/m1/reactions", gotPath)
		}
	})

	t.Run("create requires message_id and reaction", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { hit = true })
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_create_chat_reaction", map[string]any{"team_id": "123", "message_id": "m1"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing reaction)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("create wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		var gotBody map[string]any
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			_, _ = w.Write([]byte(`{"reaction":"thumbsup"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_create_chat_reaction", map[string]any{
			"team_id": "123", "message_id": "m1", "reaction": "thumbsup",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/workspaces/123/chat/messages/m1/reactions" {
			t.Errorf("path = %q, want /workspaces/123/chat/messages/m1/reactions", gotPath)
		}
		if gotBody["reaction"] != "thumbsup" {
			t.Errorf("body = %+v", gotBody)
		}
	})

	t.Run("delete requires message_id and reaction", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { hit = true })
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_delete_chat_reaction", map[string]any{"team_id": "123", "message_id": "m1"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing reaction)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("delete wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_delete_chat_reaction", map[string]any{
			"team_id": "123", "message_id": "m1", "reaction": "thumbsup",
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/workspaces/123/chat/messages/m1/reactions/thumbsup" {
			t.Errorf("path = %q, want /workspaces/123/chat/messages/m1/reactions/thumbsup", gotPath)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_list_chat_reactions", map[string]any{"team_id": "123", "message_id": "m1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupChatFollowersAndMembers(t *testing.T) {
	t.Run("followers requires channel_id", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { hit = true })
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_list_chat_followers", map[string]any{"team_id": "123"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing channel_id)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("followers wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"followers":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_list_chat_followers", map[string]any{"team_id": "123", "channel_id": "c1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %q, want GET", gotMethod)
		}
		if gotPath != "/workspaces/123/chat/channels/c1/followers" {
			t.Errorf("path = %q, want /workspaces/123/chat/channels/c1/followers", gotPath)
		}
	})

	t.Run("members requires channel_id", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { hit = true })
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_list_chat_members", map[string]any{"team_id": "123"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing channel_id)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("members wiring", func(t *testing.T) {
		var gotMethod, gotPath string
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"members":[]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_list_chat_members", map[string]any{"team_id": "123", "channel_id": "c1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodGet {
			t.Errorf("method = %q, want GET", gotMethod)
		}
		if gotPath != "/workspaces/123/chat/channels/c1/members" {
			t.Errorf("path = %q, want /workspaces/123/chat/channels/c1/members", gotPath)
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterChatTools(s, c)
		res := callTool(t, s, "clickup_list_chat_members", map[string]any{"team_id": "123", "channel_id": "c1"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}
