package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/database"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/repository"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/service"
)

type Container struct {
	// repos
	Repositories Repositories

	//services
	Services Services
	// transaction man
	TxManager *database.TxManager
}

func NewContainer(db *pgxpool.Pool) *Container {
	txManager := database.NewTxManager(db)
	userRepository := repository.NewUserRepository(db)
	roleRepository := repository.NewRoleRepository(db)
	eventRepository := repository.NewEventRepository(db)
	permissionRepository := repository.NewPermissionRepository(db)
	rolePermissionRepository := repository.NewRolePermissionRepository(db)

	userService := service.NewUserService(userRepository, roleRepository)
	roleService := service.NewRoleService(txManager, roleRepository, eventRepository)
	permissionService := service.NewPermissionService(permissionRepository)
	rolePermissionService := service.NewRolePermissionService(
		txManager,
		roleRepository,
		permissionRepository,
		rolePermissionRepository,
		eventRepository,
	)

	return &Container{
		Repositories: Repositories{
			User:       userRepository,
			Role:       roleRepository,
			Event:      eventRepository,
			Permission: permissionRepository,
		},
		Services: Services{
			User:           userService,
			Role:           roleService,
			Permission:     permissionService,
			RolePermission: rolePermissionService,
		},
	}
}
