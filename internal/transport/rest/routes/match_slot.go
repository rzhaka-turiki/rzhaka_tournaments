package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"
)

func RegisterMatchSlots(rg *gin.RouterGroup, h *handlers.MatchSlotHandler) {
	match := rg.Group("organisations/:organisation_id/matches/:match_id/slot")
	{
		match.POST("", h.Create)
		match.GET("/:slot_id", h.GetByID)
		match.GET("", h.GetByMatchID)
		match.PATCH("/:slot_id", h.Update)
		match.DELETE(":slot_id", h.Delete)
	}
}
