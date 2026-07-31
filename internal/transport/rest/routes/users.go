package routes

import (
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/handlers"
	"github.com/gin-gonic/gin"
)

func registerUsers(router *gin.RouterGroup, handler *handlers.UserHandler) {
	users := router.Group("/users")
	{
		users.GET(":id", handler.GetByID)
	}
}
