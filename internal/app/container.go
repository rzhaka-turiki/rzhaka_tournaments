package app

import (
	"github.com/a1uka/rzhaka_tournaments/internal/repository"
	"github.com/a1uka/rzhaka_tournaments/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	// repos
	Repositories Repositories

	//services
	Services Services
}

func NewContainer(db *pgxpool.Pool) *Container {
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	return &Container{
		Repositories: Repositories{
			User: userRepository,
		},
		Services: Services{
			User: userService,
		},
	}
}
