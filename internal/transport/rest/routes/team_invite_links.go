package routes

import (
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/handlers"
	"github.com/gin-gonic/gin"
)

func registerTeamInviteLinks(rg *gin.RouterGroup, h *handlers.TeamInviteLinkHandler) {

	teams := rg.Group("/teams")
	{
		teams.POST("/:id/invite-links", h.Create)
		teams.GET("/:id/invite-links", h.GetTeamLinks)
		teams.DELETE("/:id/invite-links/:linkID", h.DeleteLink)
	}

	inviteLinks := rg.Group("/invite-links")
	{
		inviteLinks.GET("/:token", h.PreviewLink)
		inviteLinks.POST("/:token/accept", h.AcceptByLink)
	}
}
