package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"
)

func registerRolePermissions(rg *gin.RouterGroup, h *handlers.RolePermissionHandler) {
	rg.GET("/roles/:id/permissions", h.GetRolePermissions)
	rg.POST("/roles/:id/permissions", h.AddPermission)
	rg.DELETE("/roles/:id/permissions/:permissionID", h.RemovePermission)
}
