package app

import (
	"github.com/a1uka/rzhaka_tournaments/internal/repository"
	"github.com/a1uka/rzhaka_tournaments/internal/service"
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	// repos
	UserRepository repository.UserRepository

	//services
	UserService service.UserService
	UserHandler *handlers.UserHandler
}

func NewContainer(db *pgxpool.Pool) *Container {
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)
	return &Container{
		UserRepository: userRepository,
		UserService:    userService,
		UserHandler:    userHandler,
	}
}
