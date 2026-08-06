package routes

import (
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/handlers"
	"github.com/gin-gonic/gin"
)

func registerHealth(router *gin.Engine, handler *handlers.HealthHandler) {
	router.GET("/health", handler.Health)
	router.GET("/health/db", handler.DBHealth)
}
