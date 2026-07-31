package rest

import (
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	healthHandler *handlers.HealthHandler,
	users *handlers.UserHandler,
) {
	// service ends
	router.GET("/health", healthHandler.Health)
	router.GET("health/db", healthHandler.DBHealth)
	v1 := router.Group("/api/v1")
	{
		usersGroup := v1.Group("/users")
		{
			usersGroup.GET("/:id", users.GetByID)
		}
	}
}
