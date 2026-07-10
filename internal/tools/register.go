package tools

import (
	"github.com/mark3labs/mcp-go/server"

	"clickup-mcp/internal/clickup"
)

// RegisterAll registers every implemented phase's tools on s, plus the
// standalone server-metadata tool (clickup_get_server_version).
func RegisterAll(s *server.MCPServer, c *clickup.Client) {
	RegisterPhase1(s, c)
	RegisterPhase2(s, c)
	RegisterPhase3(s, c)
	RegisterPhase4(s, c)
	RegisterVersionTools(s, c)
}

// RegisterPhase1 registers Auth/Workspace tools: user, workspaces, spaces,
// folders, lists.
func RegisterPhase1(s *server.MCPServer, c *clickup.Client) {
	RegisterAuthTools(s, c)
	RegisterSpaceTools(s, c)
	RegisterFolderTools(s, c)
	RegisterListTools(s, c)
}

// RegisterPhase2 registers Task tools: task CRUD/search, comments,
// checklists, custom fields, dependencies, tags.
func RegisterPhase2(s *server.MCPServer, c *clickup.Client) {
	RegisterTaskTools(s, c)
	RegisterCommentTools(s, c)
	RegisterChecklistTools(s, c)
	RegisterCustomFieldTools(s, c)
	RegisterDependencyTools(s, c)
	RegisterTagTools(s, c)
}

// RegisterPhase3 registers Goals, Time Tracking, Views, and cross-cutting
// time-tracking report tools.
func RegisterPhase3(s *server.MCPServer, c *clickup.Client) {
	RegisterGoalTools(s, c)
	RegisterTimeTrackingTools(s, c)
	RegisterViewTools(s, c)
	RegisterReportTools(s, c)
}

// RegisterPhase4 registers everything else: webhooks, guests, members,
// users, roles, templates, attachments, audit logs, chat, custom task types,
// shared hierarchy, and docs.
func RegisterPhase4(s *server.MCPServer, c *clickup.Client) {
	RegisterWebhookTools(s, c)
	RegisterGuestTools(s, c)
	RegisterMemberTools(s, c)
	RegisterUserTools(s, c)
	RegisterRoleTools(s, c)
	RegisterTemplateTools(s, c)
	RegisterAttachmentTools(s, c)
	RegisterAuditLogTools(s, c)
	RegisterChatTools(s, c)
	RegisterCustomTaskTypeTools(s, c)
	RegisterSharedHierarchyTools(s, c)
	RegisterDocTools(s, c)
}
