package app

import "github.com/rzhaka-turiki/rzhaka_tournaments/internal/service"

type Services struct {
	User       service.UserService
	Role       service.RoleService
	Permission service.PermissionService
	Team       service.TeamService

	Organisation service.OrganisationService

	Token service.TokenService

	RolePermission service.RolePermissionService
	ApexAccount    service.ApexAccountService

	Match service.MatchService
}
