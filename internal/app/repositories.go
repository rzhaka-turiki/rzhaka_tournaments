package app

import "github.com/rzhaka-turiki/rzhaka_tournaments/internal/repository"

type Repositories struct {
	User       repository.UserRepository
	Role       repository.RoleRepository
	Event      repository.EventRepository
	Permission repository.PermissionRepository
}
