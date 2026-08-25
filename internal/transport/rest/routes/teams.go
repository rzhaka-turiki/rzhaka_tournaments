package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"
)

func registerTeams(rg *gin.RouterGroup, h *handlers.TeamHandler) {
	teams := rg.Group("/teams")
	{
		teams.POST("", h.Create)
		teams.GET("/:id", h.GetByID)

		teams.PATCH("/:id", h.Update)

		teams.POST("/:id/archive", h.Archive)
		teams.POST("/:id/restore", h.Restore)

		teams.GET("/:id/members", h.GetMembers)

		teams.DELETE("/:id/members/:userID", h.RemoveMember)

		teams.POST("/:id/transfer-ownership", h.TransferOwnership)
	}
	me := rg.Group("/users/me")
	{
		me.GET("/teams/owned", h.GetOwnedTeams)
		me.GET("/teams/member", h.GetMemberTeams)
	}
}
