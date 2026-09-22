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

type MatchSlotHandler struct {
	matchSlotService service.MatchSlotService
}

func NewMatchSlotHandler(matchSlotService service.MatchSlotService) *MatchSlotHandler {
	return &MatchSlotHandler{
		matchSlotService: matchSlotService,
	}
}

// POST api/v1/organisations/:organisation_id/matches/:match_id/slot
func (h *MatchSlotHandler) Create(c *gin.Context) {
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
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	var req dto.CreateMatchSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	matchSlot := &model.MatchSlot{
		MatchID:    matchID,
		DropSpotID: req.DropSpotID,
		SlotNumber: req.SlotNumber,
	}
	err = h.matchSlotService.Create(c.Request.Context(), actorID, organisationID, matchSlot)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, dto.FromMatchSlot(matchSlot))
}

// GET api/v1/organisations/:organisation_id/matches/:match_id/slot/:slot_id
func (h *MatchSlotHandler) GetByID(c *gin.Context) {
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
	matchSlot, err := h.matchSlotService.GetByID(c.Request.Context(), slotID, organisationID, matchID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromMatchSlot(matchSlot))
}

// GET api/v1/organisations/:organisation_id/matches/:match_id/slot
func (h *MatchSlotHandler) GetByMatchID(c *gin.Context) {
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
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	matchSlots, err := h.matchSlotService.ListByMatchID(c.Request.Context(), matchID, organisationID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromMatchSlots(matchSlots))
}

// PATCH api/v1/organisations/:organisation_id/matches/:match_id/slot/:slot_id
func (h *MatchSlotHandler) Update(c *gin.Context) {
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
	var req dto.UpdateMatchSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	matchSlotUPD := service.MatchSlotUpdate{
		SlotNumber: req.SlotNumber,
		DropSpotID: req.DropSpotID,
	}
	err = h.matchSlotService.Update(c.Request.Context(), actorID, organisationID, matchID, slotID, matchSlotUPD)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

// DELETE api/v1/organisations/:organisation_id/matches/:match_id/slot/:slot_id
func (h *MatchSlotHandler) Delete(c *gin.Context) {
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
	err = h.matchSlotService.Delete(c.Request.Context(), actorID, slotID, matchID, organisationID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}
