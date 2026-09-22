package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"
)

func registerMatches(rg *gin.RouterGroup, h *handlers.MatchHandler) {
	matches := rg.Group("/organisations/:organisation_id/matches")
	{
		matches.POST("", h.Create)
		matches.GET("", h.GetByOrganisationID)
		matches.GET("/:match_id", h.GetByID)
		matches.PATCH("/:match_id", h.Update)
		matches.DELETE("/:match_id", h.Delete)
	}
}
