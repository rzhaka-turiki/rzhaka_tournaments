package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"
)

func registerMatchSettings(rg *gin.RouterGroup, h *handlers.MatchSettingsHandler) {
	matches := rg.Group("organisations/:organisation_id/matches/:match_id/settings")
	{
		matches.POST("", h.Create)
		matches.GET("", h.GetByID)
		matches.PATCH("", h.Update)
		matches.DELETE("", h.Delete)
	}
}
