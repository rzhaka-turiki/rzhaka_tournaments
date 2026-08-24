package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"
)

func registerHealth(router *gin.Engine, handler *handlers.HealthHandler) {
	router.GET("/health", handler.Health)
	router.GET("/health/db", handler.DBHealth)
}
