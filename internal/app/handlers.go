package app

import "github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"

type Handlers struct {
	Health *handlers.HealthHandler
	User   *handlers.UserHandler
}
