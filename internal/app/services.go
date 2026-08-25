package app

import "github.com/rzhaka-turiki/rzhaka_tournaments/internal/service"

type Services struct {
	User       service.UserService
	Role       service.RoleService
	Permission service.PermissionService

	RolePermission service.RolePermissionService
	ApexAccount    service.ApexAccountService
}
