package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

func RegisterChatTools(s *server.MCPServer, c *clickup.Client) {
	s.AddTool(
		mcp.NewTool("clickup_list_chat_channels",
			mcp.WithDescription("List the chat channels in a ClickUp workspace."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			out, err := c.ListChatChannels(ctx, teamIDOrDefault(req, c))
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_chat_channel",
			mcp.WithDescription("Create a chat channel in a ClickUp workspace. Channels are PUBLIC "+
				"— visible and joinable by the whole workspace — by default; this tool does not "+
				"support setting a private visibility or an initial member list."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Channel name")),
			mcp.WithString("description", mcp.Description("Channel description")),
			mcp.WithString("topic", mcp.Description("Channel topic")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name, err := req.RequireString("name")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{"name": name}
			setString(body, req, "description")
			setString(body, req, "topic")
			out, err := c.CreateChatChannel(ctx, teamIDOrDefault(req, c), body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_get_chat_channel",
			mcp.WithDescription("Get a single ClickUp chat channel."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			channelID, err := req.RequireString("channel_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.GetChatChannel(ctx, teamIDOrDefault(req, c), channelID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_chat_channel",
			mcp.WithDescription("Update a ClickUp chat channel."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
			mcp.WithString("name", mcp.Description("Channel name")),
			mcp.WithString("description", mcp.Description("Channel description")),
			mcp.WithString("topic", mcp.Description("Channel topic")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			channelID, err := req.RequireString("channel_id")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{}
			setString(body, req, "name")
			setString(body, req, "description")
			setString(body, req, "topic")
			out, err := c.UpdateChatChannel(ctx, teamIDOrDefault(req, c), channelID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_chat_messages",
			mcp.WithDescription("List the messages in a ClickUp chat channel."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			channelID, err := req.RequireString("channel_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListChatMessages(ctx, teamIDOrDefault(req, c), channelID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_chat_message",
			mcp.WithDescription("Post a message to a ClickUp chat channel."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
			mcp.WithString("content", mcp.Required(), mcp.Description("Message content")),
			mcp.WithString("type", mcp.Description(`Message type: "message" (default) or "post"`)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			channelID, err := req.RequireString("channel_id")
			if err != nil {
				return ErrorResult(err)
			}
			content, err := req.RequireString("content")
			if err != nil {
				return ErrorResult(err)
			}
			body := map[string]any{"content": content, "type": req.GetString("type", "message")}
			out, err := c.CreateChatMessage(ctx, teamIDOrDefault(req, c), channelID, body)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_update_chat_message",
			mcp.WithDescription("Edit a ClickUp chat message."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("message_id", mcp.Required(), mcp.Description("Message ID")),
			mcp.WithString("content", mcp.Required(), mcp.Description("New message content")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			messageID, err := req.RequireString("message_id")
			if err != nil {
				return ErrorResult(err)
			}
			content, err := req.RequireString("content")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.UpdateChatMessage(ctx, teamIDOrDefault(req, c), messageID, map[string]any{"content": content})
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_delete_chat_message",
			mcp.WithDescription("Delete a ClickUp chat message."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("message_id", mcp.Required(), mcp.Description("Message ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			messageID, err := req.RequireString("message_id")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.DeleteChatMessage(ctx, teamIDOrDefault(req, c), messageID); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"deleted": true, "message_id": messageID})
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_chat_reactions",
			mcp.WithDescription("List the reactions on a ClickUp chat message."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("message_id", mcp.Required(), mcp.Description("Message ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			messageID, err := req.RequireString("message_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListChatReactions(ctx, teamIDOrDefault(req, c), messageID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_create_chat_reaction",
			mcp.WithDescription("Add a reaction to a ClickUp chat message."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("message_id", mcp.Required(), mcp.Description("Message ID")),
			mcp.WithString("reaction", mcp.Required(), mcp.Description("Emoji name, e.g. thumbsup")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			messageID, err := req.RequireString("message_id")
			if err != nil {
				return ErrorResult(err)
			}
			reaction, err := req.RequireString("reaction")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.CreateChatReaction(ctx, teamIDOrDefault(req, c), messageID, map[string]any{"reaction": reaction})
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_delete_chat_reaction",
			mcp.WithDescription("Remove a reaction from a ClickUp chat message."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("message_id", mcp.Required(), mcp.Description("Message ID")),
			mcp.WithString("reaction", mcp.Required(), mcp.Description("Emoji name, e.g. thumbsup")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			messageID, err := req.RequireString("message_id")
			if err != nil {
				return ErrorResult(err)
			}
			reaction, err := req.RequireString("reaction")
			if err != nil {
				return ErrorResult(err)
			}
			if err := c.DeleteChatReaction(ctx, teamIDOrDefault(req, c), messageID, reaction); err != nil {
				return ErrorResult(err)
			}
			return JSONResult(map[string]any{"deleted": true, "message_id": messageID, "reaction": reaction})
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_chat_followers",
			mcp.WithDescription("List the followers of a ClickUp chat channel."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			channelID, err := req.RequireString("channel_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListChatFollowers(ctx, teamIDOrDefault(req, c), channelID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)

	s.AddTool(
		mcp.NewTool("clickup_list_chat_members",
			mcp.WithDescription("List the members of a ClickUp chat channel."),
			mcp.WithString("team_id", mcp.Description("Workspace ID; defaults to CLICKUP_TEAM_ID")),
			mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			channelID, err := req.RequireString("channel_id")
			if err != nil {
				return ErrorResult(err)
			}
			out, err := c.ListChatMembers(ctx, teamIDOrDefault(req, c), channelID)
			if err != nil {
				return ErrorResult(err)
			}
			return JSONResult(out)
		},
	)
}
