package rest

import "github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"

type Handlers struct {
	Health          *handlers.HealthHandler
	User            *handlers.UserHandler
	Role            *handlers.RoleHandler
	Permission      *handlers.PermissionHandler
	RolePermission  *handlers.RolePermissionHandler
	Team            *handlers.TeamHandler
	TeamInviteLinks *handlers.TeamInviteLinkHandler
	TeamInvites     *handlers.TeamRequestHandler
	ApexAccounts    *handlers.ApexAccountHandler
}
