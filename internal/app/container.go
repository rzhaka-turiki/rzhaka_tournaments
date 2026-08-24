package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/client/apexverifier"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/config"
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
	Clients   Clients
}

func NewContainer(db *pgxpool.Pool, cfg *config.Config) (*Container, error) {
	//tx man
	txManager := database.NewTxManager(db)
	// repos
	userRepository := repository.NewUserRepository(db)
	roleRepository := repository.NewRoleRepository(db)
	eventRepository := repository.NewEventRepository(db)
	permissionRepository := repository.NewPermissionRepository(db)
	rolePermissionRepository := repository.NewRolePermissionRepository(db)
	apexAccountRepository := repository.NewApexAccountRepository(db)
	// cleints
	apexVerifierClient, err := apexverifier.NewClient(cfg.ApexVerifier.GRPCAddr)
	// services
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
	apexAccountService := service.NewApexAccountService(apexAccountRepository, apexVerifierClient)

	if err != nil {
		return nil, err
	}
	return &Container{
		Repositories: Repositories{
			User:        userRepository,
			Role:        roleRepository,
			Event:       eventRepository,
			Permission:  permissionRepository,
			ApexAccount: apexAccountRepository,
		},
		Services: Services{
			User:           userService,
			Role:           roleService,
			Permission:     permissionService,
			RolePermission: rolePermissionService,
			ApexAccount:    apexAccountService,
		},
		Clients: Clients{
			ApexVerifier: apexVerifierClient,
		},
		TxManager: txManager,
	}, nil
}
