package routes

import (
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/handlers"
	"github.com/gin-gonic/gin"
)

func registerTeamRequests(
	rg *gin.RouterGroup,
	h *handlers.TeamRequestHandler,
) {

	teams := rg.Group("/teams")
	{
		teams.POST("/:id/invitations", h.CreateInvite)

		teams.GET("/:id/invitations", h.GetTeamInvitations)
	}

	users := rg.Group("/users/me")
	{
		users.GET("/invitations", h.GetMyInvitations)
	}

	invitations := rg.Group("/invitations")
	{
		invitations.POST("/:id/accept", h.AcceptRequest)

		invitations.POST("/:id/reject", h.RejectRequest)
	}
}
