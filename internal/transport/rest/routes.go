package rest

import (
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, health *handlers.HealthHandler) {
	router.GET("/health", health.Health)
}
