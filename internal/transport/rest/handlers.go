package rest

import "github.com/a1uka/rzhaka_tournaments/internal/transport/rest/handlers"

type Handlers struct {
	Health          *handlers.HealthHandler
	User            *handlers.UserHandler
	Role            *handlers.RoleHandler
	Permission      *handlers.PermissionHandler
	RolePermission  *handlers.RolePermissionHandler
	Team            *handlers.TeamHandler
	TeamInviteLinks *handlers.TeamInviteLinkHandler
	TeamInvites     *handlers.TeamRequestHandler
}
