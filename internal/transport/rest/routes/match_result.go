package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"
)

func registerMatchResults(rg *gin.RouterGroup, h *handlers.MatchResultHandler) {
	matches := rg.Group("organisations/:organisation_id/matches")
	{
		matches.POST("/:match_id/result", h.Create)
		matches.GET("/:match_id/result", h.GetByMatchID)
		matches.GET("/raw/:external_mid/result", h.GetByExternalMID)
		matches.PATCH("/:match_id/result", h.Update)
		matches.DELETE("/:match_id/result", h.Delete)
	}
}
