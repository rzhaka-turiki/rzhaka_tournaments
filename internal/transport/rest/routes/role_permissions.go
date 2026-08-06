package routes

import (
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/handlers"
	"github.com/gin-gonic/gin"
)

func registerRolePermissions(rg *gin.RouterGroup, h *handlers.RolePermissionHandler) {
	rg.GET("/roles/:id/permissions", h.GetRolePermissions)
	rg.POST("/roles/:id/permissions", h.AddPermission)
	rg.DELETE("/roles/:id/permissions/:permissionID", h.RemovePermission)
}
