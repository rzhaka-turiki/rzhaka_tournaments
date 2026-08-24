package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest"
)

func RegisterRoutes(
	router *gin.Engine,
	h *rest.Handlers,
) {
	registerHealth(router, h.Health)
	v1 := router.Group("/api/v1")
	{
		registerUsers(v1, h.User)
		registerRoles(v1, h.Role)
		registerPermissions(v1, h.Permission)
		registerRolePermissions(v1, h.RolePermission)
		registerTeams(v1, h.Team)
		registerTeamInviteLinks(v1, h.TeamInviteLinks)
		registerTeamRequests(v1, h.TeamInvites)
		registerApexAccounts(v1, h.ApexAccounts)
	}
}
