package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"
)

func registerTokens(rg *gin.RouterGroup, h *handlers.TokenHandler) {
	tokens := rg.Group("/organisations/:organisation_id/tokens")
	{
		tokens.POST("", h.Create)
		tokens.GET("", h.GetByOrganisationID)
		tokens.GET("/:token_id", h.GetByID)
		tokens.PATCH("/:token_id", h.Update)
		tokens.DELETE("/:token_id", h.Delete)
	}
}
