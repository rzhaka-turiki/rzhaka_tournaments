package app

import "github.com/a1uka/rzhaka_tournaments/internal/service"

type Services struct {
	User       service.UserService
	Role       service.RoleService
	Permission service.PermissionService

	RolePermission service.RolePermissionService
}
