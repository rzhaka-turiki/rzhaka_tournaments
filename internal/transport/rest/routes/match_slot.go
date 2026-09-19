package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"
)

func registerMatchSlots(rg *gin.RouterGroup, h *handlers.MatchSlotHandler) {
	match := rg.Group("organisations/:organisation_id/matches/:match_id/slot")
	{
		match.POST("", h.Create)
		match.GET("/:slot_id", h.GetByID)
		match.GET("", h.GetByMatchID)
		match.PATCH("/:slot_id", h.Update)
		match.DELETE(":slot_id", h.Delete)
	}
}

func registerMatchSlotPlayers(rg *gin.RouterGroup, h *handlers.MatchSlotPlayerHandler) {
	match := rg.Group("organisations/:organisation_id/matches/:match_id/slot/:slot_id")
	{
		match.POST("", h.Create)
		match.GET("/:slot_player_id", h.GetByID)
		match.GET("", h.GetBySlotID)
		match.PATCH("/:slot_player_id", h.Update)
		match.DELETE("/:slot_player_id", h.Delete)
	}
}
