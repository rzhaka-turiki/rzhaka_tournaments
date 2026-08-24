package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"
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
