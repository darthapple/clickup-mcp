package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

// Note: every method in chat.go sets APIVersion: apiV3; since testClient
// points BaseURLv2 and BaseURLv3 at the same fake server, that's exercised
// implicitly by every request below hitting the handler at all.

func TestListChatChannelsHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"channels":[{"id":"c1"}]}`))
	})

	out, err := c.ListChatChannels(context.Background(), "999")
	if err != nil {
		t.Fatalf("ListChatChannels: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/workspaces/999/chat/channels" {
		t.Errorf("path = %q, want /workspaces/999/chat/channels", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded type = %T, want map[string]any", out)
	}
	if channels, ok := m["channels"].([]any); !ok || len(channels) != 1 {
		t.Errorf("channels = %+v", m["channels"])
	}
}

func TestCreateChatChannelPostsExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"c1"}`))
	})

	out, err := c.CreateChatChannel(context.Background(), "999", map[string]any{"name": "general"})
	if err != nil {
		t.Fatalf("CreateChatChannel: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotPath != "/workspaces/999/chat/channels" {
		t.Errorf("path = %q, want /workspaces/999/chat/channels", gotPath)
	}
	if gotBody["name"] != "general" {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok || m["id"] != "c1" {
		t.Errorf("decoded = %+v", out)
	}
}

func TestGetChatChannelHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"id":"c1"}`))
	})

	out, err := c.GetChatChannel(context.Background(), "999", "c1")
	if err != nil {
		t.Fatalf("GetChatChannel: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/workspaces/999/chat/channels/c1" {
		t.Errorf("path = %q, want /workspaces/999/chat/channels/c1", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok || m["id"] != "c1" {
		t.Errorf("decoded = %+v", out)
	}
}

func TestUpdateChatChannelPatchesExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"c1","name":"renamed"}`))
	})

	out, err := c.UpdateChatChannel(context.Background(), "999", "c1", map[string]any{"name": "renamed"})
	if err != nil {
		t.Fatalf("UpdateChatChannel: %v", err)
	}
	if gotMethod != http.MethodPatch {
		t.Errorf("method = %q, want PATCH", gotMethod)
	}
	if gotPath != "/workspaces/999/chat/channels/c1" {
		t.Errorf("path = %q, want /workspaces/999/chat/channels/c1", gotPath)
	}
	if gotBody["name"] != "renamed" {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok || m["name"] != "renamed" {
		t.Errorf("decoded = %+v", out)
	}
}

func TestListChatMessagesHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"messages":[{"id":"m1"}]}`))
	})

	out, err := c.ListChatMessages(context.Background(), "999", "c1")
	if err != nil {
		t.Fatalf("ListChatMessages: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/workspaces/999/chat/channels/c1/messages" {
		t.Errorf("path = %q, want /workspaces/999/chat/channels/c1/messages", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded type = %T, want map[string]any", out)
	}
	if messages, ok := m["messages"].([]any); !ok || len(messages) != 1 {
		t.Errorf("messages = %+v", m["messages"])
	}
}

func TestCreateChatMessagePostsExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"m1"}`))
	})

	out, err := c.CreateChatMessage(context.Background(), "999", "c1", map[string]any{"content": "hello"})
	if err != nil {
		t.Fatalf("CreateChatMessage: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotPath != "/workspaces/999/chat/channels/c1/messages" {
		t.Errorf("path = %q, want /workspaces/999/chat/channels/c1/messages", gotPath)
	}
	if gotBody["content"] != "hello" {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok || m["id"] != "m1" {
		t.Errorf("decoded = %+v", out)
	}
}

func TestUpdateChatMessagePatchesExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"m1","content":"edited"}`))
	})

	out, err := c.UpdateChatMessage(context.Background(), "999", "m1", map[string]any{"content": "edited"})
	if err != nil {
		t.Fatalf("UpdateChatMessage: %v", err)
	}
	if gotMethod != http.MethodPatch {
		t.Errorf("method = %q, want PATCH", gotMethod)
	}
	if gotPath != "/workspaces/999/chat/messages/m1" {
		t.Errorf("path = %q, want /workspaces/999/chat/messages/m1", gotPath)
	}
	if gotBody["content"] != "edited" {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok || m["content"] != "edited" {
		t.Errorf("decoded = %+v", out)
	}
}

func TestDeleteChatMessageHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	if err := c.DeleteChatMessage(context.Background(), "999", "m1"); err != nil {
		t.Fatalf("DeleteChatMessage: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %q, want DELETE", gotMethod)
	}
	if gotPath != "/workspaces/999/chat/messages/m1" {
		t.Errorf("path = %q, want /workspaces/999/chat/messages/m1", gotPath)
	}
}

func TestListChatReactionsHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"reactions":[{"emoji":"👍"}]}`))
	})

	out, err := c.ListChatReactions(context.Background(), "999", "m1")
	if err != nil {
		t.Fatalf("ListChatReactions: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/workspaces/999/chat/messages/m1/reactions" {
		t.Errorf("path = %q, want /workspaces/999/chat/messages/m1/reactions", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded type = %T, want map[string]any", out)
	}
	if reactions, ok := m["reactions"].([]any); !ok || len(reactions) != 1 {
		t.Errorf("reactions = %+v", m["reactions"])
	}
}

func TestCreateChatReactionPostsExpectedRequest(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"emoji":"👍"}`))
	})

	out, err := c.CreateChatReaction(context.Background(), "999", "m1", map[string]any{"reaction": "👍"})
	if err != nil {
		t.Fatalf("CreateChatReaction: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotPath != "/workspaces/999/chat/messages/m1/reactions" {
		t.Errorf("path = %q, want /workspaces/999/chat/messages/m1/reactions", gotPath)
	}
	if gotBody["reaction"] != "👍" {
		t.Errorf("body = %+v", gotBody)
	}
	m, ok := out.(map[string]any)
	if !ok || m["emoji"] != "👍" {
		t.Errorf("decoded = %+v", out)
	}
}

func TestDeleteChatReactionHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	if err := c.DeleteChatReaction(context.Background(), "999", "m1", "thumbsup"); err != nil {
		t.Fatalf("DeleteChatReaction: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %q, want DELETE", gotMethod)
	}
	if gotPath != "/workspaces/999/chat/messages/m1/reactions/thumbsup" {
		t.Errorf("path = %q, want /workspaces/999/chat/messages/m1/reactions/thumbsup", gotPath)
	}
}

func TestListChatFollowersHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"followers":[{"id":1}]}`))
	})

	out, err := c.ListChatFollowers(context.Background(), "999", "c1")
	if err != nil {
		t.Fatalf("ListChatFollowers: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/workspaces/999/chat/channels/c1/followers" {
		t.Errorf("path = %q, want /workspaces/999/chat/channels/c1/followers", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded type = %T, want map[string]any", out)
	}
	if followers, ok := m["followers"].([]any); !ok || len(followers) != 1 {
		t.Errorf("followers = %+v", m["followers"])
	}
}

func TestListChatMembersHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"members":[{"id":1}]}`))
	})

	out, err := c.ListChatMembers(context.Background(), "999", "c1")
	if err != nil {
		t.Fatalf("ListChatMembers: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/workspaces/999/chat/channels/c1/members" {
		t.Errorf("path = %q, want /workspaces/999/chat/channels/c1/members", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("decoded type = %T, want map[string]any", out)
	}
	if members, ok := m["members"].([]any); !ok || len(members) != 1 {
		t.Errorf("members = %+v", m["members"])
	}
}

func TestListChatChannelsReturnsAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"err":"Chat not enabled","ECODE":"CHAT_001"}`))
	})

	_, err := c.ListChatChannels(context.Background(), "999")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusForbidden || apiErr.ECode != "CHAT_001" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}
