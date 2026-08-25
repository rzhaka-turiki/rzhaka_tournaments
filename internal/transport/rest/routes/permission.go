package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"
)

func registerPermissions(rg *gin.RouterGroup, h *handlers.PermissionHandler) {
	permissions := rg.Group("/permissions")
	{
		permissions.GET("", h.GetAll)
		permissions.GET("/:id", h.GetByID)
		permissions.GET("/code/:code", h.GetByCode)
	}
}
