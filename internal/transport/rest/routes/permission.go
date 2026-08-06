package routes

import (
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/handlers"
	"github.com/gin-gonic/gin"
)

func registerPermissions(rg *gin.RouterGroup, h *handlers.PermissionHandler) {
	permissions := rg.Group("/permissions")
	{
		permissions.GET("", h.GetAll)
		permissions.GET("/:id", h.GetByID)
		permissions.GET("/code/:code", h.GetByCode)
	}
}
