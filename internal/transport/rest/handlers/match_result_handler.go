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

type MatchResultHandler struct {
	matchResultService service.MatchResultService
}

func NewMatchResultHandler(matchResultService service.MatchResultService) *MatchResultHandler {
	return &MatchResultHandler{
		matchResultService: matchResultService,
	}
}

// POST api/v1/organisations/:organisation_id/matches/:match_id/result
func (h *MatchResultHandler) Create(c *gin.Context) {
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
	var req dto.CreateMatchResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	matchResult := &model.MatchResult{
		MatchID:          matchID,
		MapID:            req.MapID,
		MapName:          req.MapName,
		ExternalMID:      req.ExternalMID,
		AimAssistAllowed: req.AimAssistAllowed,
	}
	err = h.matchResultService.Create(c.Request.Context(), actorID, organisationID, matchID, matchResult)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, dto.FromMatchResult(matchResult))
}

// GET api/v1/organisations/:organisation_id/matches/:match_id/result
func (h *MatchResultHandler) GetByMatchID(c *gin.Context) {
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
	matchResult, err := h.matchResultService.GetByMatchID(c.Request.Context(), matchID, organisationID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromMatchResult(matchResult))
}

// GET api/v1/organisations/:organisation_id/matches/raw/:external_mid/result
func (h *MatchResultHandler) GetByExternalMID(c *gin.Context) {
	organisationID, err := uuid.Parse(c.Param("organisation_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ORGANISATION_ID", "invalid organisation id")
		return
	}
	externalMID := c.Param("external_mid")
	matchResult, err := h.matchResultService.GetByExternalMID(c.Request.Context(), externalMID, organisationID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromMatchResult(matchResult))
}

// PATCH api/v1/organisations/:organisation_id/matches/:match_id/result
func (h *MatchResultHandler) Update(c *gin.Context) {
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
	var req dto.UpdateMatchResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	resultUPD := service.MatchResultUpdate{
		ExternalMID:      req.ExternalMID,
		MapName:          req.MapName,
		StartedAt:        req.StartedAt,
		AimAssistAllowed: req.AimAssistAllowed,
		MapID:            req.MapID,
	}
	err = h.matchResultService.Update(c.Request.Context(), actorID, organisationID, matchID, resultUPD)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

// DELETE api/v1/organisations/:organisation_id/matches/:match_id/result
func (h *MatchResultHandler) Delete(c *gin.Context) {
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
	err = h.matchResultService.Delete(c.Request.Context(), actorID, organisationID, matchID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}
