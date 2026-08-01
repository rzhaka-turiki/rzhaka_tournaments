package routes

import (
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/handlers"
	"github.com/gin-gonic/gin"
)

func registerRoles(router *gin.RouterGroup, roleHandler *handlers.RoleHandler) {
	roles := router.Group("/roles")
	{
		roles.POST("", roleHandler.CreateRole)
		roles.GET("", roleHandler.GetAllRoles)
	}
	permissions := router.Group("/permissions")
	{
		permissions.GET("", roleHandler.GetAllPermissions)
	}
	users := router.Group("/users")
	{
		users.GET("/:id/roles", roleHandler.GetUserRoles)
		users.GET("/:id/permissions", roleHandler.GetUserPermissions)

		users.POST("/:id/role", roleHandler.AssignRole)
		users.DELETE("/:id/role", roleHandler.RemoveRole)
	}
}
