package rest

import (
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	healthHandler *handlers.HealthHandler,
) {
	api := router.Group("/")

	{
		api.GET("/health", healthHandler.Health)
		api.GET("/health/db", healthHandler.DBHealth)
	}
}
