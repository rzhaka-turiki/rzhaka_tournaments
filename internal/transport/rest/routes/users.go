package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"
)

func registerUsers(router *gin.RouterGroup, handler *handlers.UserHandler) {
	users := router.Group("/users")
	{
		users.GET("/me", handler.Me)
		users.GET("/:id", handler.GetByID)
	}
}
