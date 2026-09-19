package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/service"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/auth"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/dto"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/response"
)

type MatchSlotPlayerHandler struct {
	matchSlotPlayer service.MatchSlotPlayerService
}

func NewMatchSlotPlayerHandler(matchSlotPlayerService service.MatchSlotPlayerService) *MatchSlotPlayerHandler {
	return &MatchSlotPlayerHandler{
		matchSlotPlayer: matchSlotPlayerService,
	}
}

// POST api/v1/organisations/:organisation_id/matches/:match_id/slot/:slot_id
func (h *MatchSlotPlayerHandler) Create(c *gin.Context) {
	organisationID, err := uuid.Parse(c.Param("organisation_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ORGANISATION_ID", "invalid organisation id")
		return
	}
	matchID, err := uuid.Parse(c.Param("match_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_MATCH_ID", "invalid match id")
		return
	}
	slotID, err := uuid.Parse(c.Param("slot_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_SLOT_ID", "invalid slot id")
		return
	}
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	var req dto.CreateMatchSlotPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	slotPlayer := &model.MatchSlotPlayer{
		UserID:          req.UserID,
		ExpectedNIDHash: req.ExpectedNIDHash,
	}
	err = h.matchSlotPlayer.Create(c.Request.Context(), actorID, organisationID, matchID, slotID, slotPlayer)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, dto.FromMatchSlotPlayer(slotPlayer))
}

// GET api/v1/organisations/:organisation_id/matches/:match_id/slot/:slot_id/player/:slot_player_id
func (h *MatchSlotPlayerHandler) GetByID(c *gin.Context) {
	organisationID, err := uuid.Parse(c.Param("organisation_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ORGANISATION_ID", "invalid organisation id")
		return
	}
	matchID, err := uuid.Parse(c.Param("match_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_MATCH_ID", "invalid match id")
		return
	}
	slotID, err := uuid.Parse(c.Param("slot_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_SLOT_ID", "invalid slot id")
		return
	}
	playerID, err := uuid.Parse(c.Param("slot_player_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_SLOT_PLAYER_ID", "invalid slot player id")
	}
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	slotPlayer, err := h.matchSlotPlayer.GetByID(c.Request.Context(), actorID, organisationID, matchID, slotID, playerID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromMatchSlotPlayer(slotPlayer))
}

// GET api/v1/organisations/:organisation_id/matches/:match_id/slot/:slot_id/player
func (h *MatchSlotPlayerHandler) GetBySlotID(c *gin.Context) {
	organisationID, err := uuid.Parse(c.Param("organisation_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ORGANISATION_ID", "invalid organisation id")
		return
	}
	matchID, err := uuid.Parse(c.Param("match_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_MATCH_ID", "invalid match id")
		return
	}
	slotID, err := uuid.Parse(c.Param("slot_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_SLOT_ID", "invalid slot id")
		return
	}
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	slotPlayers, err := h.matchSlotPlayer.ListBySlotID(c.Request.Context(), actorID, organisationID, matchID, slotID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromMatchSlotPlayers(slotPlayers))
}

// PATCH api/v1/organisations/:organisation_id/matches/:match_id/slot/:slot_id/player/:slot_player_id
func (h *MatchSlotPlayerHandler) Update(c *gin.Context) {
	organisationID, err := uuid.Parse(c.Param("organisation_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ORGANISATION_ID", "invalid organisation id")
		return
	}
	matchID, err := uuid.Parse(c.Param("match_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_MATCH_ID", "invalid match id")
		return
	}
	slotID, err := uuid.Parse(c.Param("slot_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_SLOT_ID", "invalid slot id")
		return
	}
	playerID, err := uuid.Parse(c.Param("slot_player_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_SLOT_PLAYER_ID", "invalid slot player id")
	}
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	var req dto.UpdateMatchSlotPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	slotPlayerUPD := service.MatchSlotPlayerUpdate{
		UserID:          req.UserID,
		ExpectedNIDHash: req.ExpectedNIDHash,
	}
	err = h.matchSlotPlayer.Update(c.Request.Context(), actorID, organisationID, matchID, slotID, playerID, slotPlayerUPD)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

// DELETE api/v1/organisations/:organisation_id/matches/:match_id/slot/:slot_id/player/:slot_player_id
func (h *MatchSlotPlayerHandler) Delete(c *gin.Context) {
	organisationID, err := uuid.Parse(c.Param("organisation_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ORGANISATION_ID", "invalid organisation id")
		return
	}
	matchID, err := uuid.Parse(c.Param("match_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_MATCH_ID", "invalid match id")
		return
	}
	slotID, err := uuid.Parse(c.Param("slot_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_SLOT_ID", "invalid slot id")
		return
	}
	playerID, err := uuid.Parse(c.Param("slot_player_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_SLOT_PLAYER_ID", "invalid slot player id")
	}
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	err = h.matchSlotPlayer.Delete(c.Request.Context(), actorID, organisationID, matchID, slotID, playerID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}
