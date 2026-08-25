package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"
)

func registerRoles(router *gin.RouterGroup, roleHandler *handlers.RoleHandler) {
	roles := router.Group("/roles")
	{
		roles.POST("", roleHandler.CreateRole)
		roles.GET("", roleHandler.GetAllRoles)
		roles.DELETE("/:id", roleHandler.DeleteRole)
		roles.PATCH(":id/restore", roleHandler.RestoreRole)
	}
	users := router.Group("/users")
	{
		users.GET("/:id/roles", roleHandler.GetUserRoles)
		users.GET("/:id/permissions", roleHandler.GetUserPermissions)

		users.POST("/:id/role", roleHandler.AssignRole)
		users.DELETE("/:id/role", roleHandler.RemoveRole)
	}
}
