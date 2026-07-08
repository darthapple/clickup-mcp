package clickup

import (
	"context"
	"net/http"
)

func chatBase(workspaceID string) string {
	return "/workspaces/" + workspaceID + "/chat"
}

// ListChatChannels returns the chat channels in a workspace.
// GET /v3/workspaces/{workspace_id}/chat/channels
func (c *Client) ListChatChannels(ctx context.Context, workspaceID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, APIVersion: apiV3, Path: chatBase(workspaceID) + "/channels"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateChatChannel creates a chat channel in a workspace.
// POST /v3/workspaces/{workspace_id}/chat/channels
func (c *Client) CreateChatChannel(ctx context.Context, workspaceID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, APIVersion: apiV3, Path: chatBase(workspaceID) + "/channels", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetChatChannel returns a single chat channel.
// GET /v3/workspaces/{workspace_id}/chat/channels/{channel_id}
func (c *Client) GetChatChannel(ctx context.Context, workspaceID, channelID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, APIVersion: apiV3, Path: chatBase(workspaceID) + "/channels/" + channelID}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateChatChannel updates a chat channel.
// PATCH /v3/workspaces/{workspace_id}/chat/channels/{channel_id}
func (c *Client) UpdateChatChannel(ctx context.Context, workspaceID, channelID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPatch, APIVersion: apiV3, Path: chatBase(workspaceID) + "/channels/" + channelID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListChatMessages returns the messages in a chat channel.
// GET /v3/workspaces/{workspace_id}/chat/channels/{channel_id}/messages
func (c *Client) ListChatMessages(ctx context.Context, workspaceID, channelID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, APIVersion: apiV3, Path: chatBase(workspaceID) + "/channels/" + channelID + "/messages"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateChatMessage posts a message to a chat channel.
// POST /v3/workspaces/{workspace_id}/chat/channels/{channel_id}/messages
func (c *Client) CreateChatMessage(ctx context.Context, workspaceID, channelID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, APIVersion: apiV3, Path: chatBase(workspaceID) + "/channels/" + channelID + "/messages", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateChatMessage edits a chat message.
// PATCH /v3/workspaces/{workspace_id}/chat/messages/{message_id}
func (c *Client) UpdateChatMessage(ctx context.Context, workspaceID, messageID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPatch, APIVersion: apiV3, Path: chatBase(workspaceID) + "/messages/" + messageID, Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteChatMessage deletes a chat message.
// DELETE /v3/workspaces/{workspace_id}/chat/messages/{message_id}
func (c *Client) DeleteChatMessage(ctx context.Context, workspaceID, messageID string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, APIVersion: apiV3, Path: chatBase(workspaceID) + "/messages/" + messageID}, nil)
}

// ListChatReactions returns the reactions on a chat message.
// GET /v3/workspaces/{workspace_id}/chat/messages/{message_id}/reactions
func (c *Client) ListChatReactions(ctx context.Context, workspaceID, messageID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, APIVersion: apiV3, Path: chatBase(workspaceID) + "/messages/" + messageID + "/reactions"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateChatReaction adds a reaction to a chat message.
// POST /v3/workspaces/{workspace_id}/chat/messages/{message_id}/reactions
func (c *Client) CreateChatReaction(ctx context.Context, workspaceID, messageID string, body map[string]any) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodPost, APIVersion: apiV3, Path: chatBase(workspaceID) + "/messages/" + messageID + "/reactions", Body: body}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteChatReaction removes a reaction from a chat message.
// DELETE /v3/workspaces/{workspace_id}/chat/messages/{message_id}/reactions/{emoji}
func (c *Client) DeleteChatReaction(ctx context.Context, workspaceID, messageID, emoji string) error {
	return c.do(ctx, requestParams{Method: http.MethodDelete, APIVersion: apiV3, Path: chatBase(workspaceID) + "/messages/" + messageID + "/reactions/" + emoji}, nil)
}

// ListChatFollowers returns the followers of a chat channel.
// GET /v3/workspaces/{workspace_id}/chat/channels/{channel_id}/followers
func (c *Client) ListChatFollowers(ctx context.Context, workspaceID, channelID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, APIVersion: apiV3, Path: chatBase(workspaceID) + "/channels/" + channelID + "/followers"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListChatMembers returns the members of a chat channel.
// GET /v3/workspaces/{workspace_id}/chat/channels/{channel_id}/members
func (c *Client) ListChatMembers(ctx context.Context, workspaceID, channelID string) (any, error) {
	var out any
	if err := c.do(ctx, requestParams{Method: http.MethodGet, APIVersion: apiV3, Path: chatBase(workspaceID) + "/channels/" + channelID + "/members"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
