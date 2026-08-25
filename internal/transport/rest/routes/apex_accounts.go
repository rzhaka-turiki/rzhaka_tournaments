package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"
)

func registerApexAccounts(router *gin.RouterGroup, handler *handlers.ApexAccountHandler) {
	users := router.Group("/users")
	{
		users.GET("/me/apex-account", handler.GetMyAccount)
		users.POST("/me/apex-account", handler.Bind)
	}
}
