package app

import (
	"github.com/a1uka/rzhaka_tournaments/internal/database"
	"github.com/a1uka/rzhaka_tournaments/internal/repository"
	"github.com/a1uka/rzhaka_tournaments/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
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
	userService := service.NewUserService(userRepository)
	roleRepository := repository.NewRoleRepository(db)
	eventRepository := repository.NewEventRepository(db)
	roleService := service.NewRoleService(txManager, roleRepository, eventRepository)

	return &Container{
		Repositories: Repositories{
			User:  userRepository,
			Role:  roleRepository,
			Event: eventRepository,
		},
		Services: Services{
			User: userService,
			Role: roleService,
		},
	}
}
