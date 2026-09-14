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

type MatchSettingsHandler struct {
	matchSettingsSerivce service.MatchSettingsService
}

func NewMatchSettingsHandler(matchSettingsSerivce service.MatchSettingsService) *MatchSettingsHandler {
	return &MatchSettingsHandler{
		matchSettingsSerivce: matchSettingsSerivce,
	}
}

// POST api/v1/organisations/:organisation_id/matches/:match_id
func (h *MatchSettingsHandler) Create(c *gin.Context) {
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
	var req dto.CreateMatchSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	matchSettings := &model.MatchSettings{
		MatchID:          matchID,
		MapID:            *req.MapID,
		PlaylistName:     *req.PlaylistName,
		AdminChat:        *req.AdminChat,
		TeamRename:       *req.TeamRename,
		SelfAssign:       *req.SelfAssign,
		AimAssist:        *req.AimAssist,
		AnonMode:         *req.AnonMode,
		DropSpotsEnabled: *req.DropSpotsEnabled,
		FillBotsMode:     *req.FillBotsMode,
	}
	err = h.matchSettingsSerivce.Create(c.Request.Context(), actorID, organisationID, matchSettings)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, dto.FromMatchSettings(matchSettings))
}

// GET api/v1/organisations/:organisation_id/matches/:match_id/settings
// blya, так-то мне тут не нужен org_id, но чисто ради сохранения структуры я бы оставил его
func (h *MatchSettingsHandler) GetByID(c *gin.Context) {
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
	matchSettings, err := h.matchSettingsSerivce.GetByMatchID(c.Request.Context(), matchID, organisationID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromMatchSettings(matchSettings))
}

// PATCH api/v1/organisations/:organisation_id/matches/:match_id/settings
func (h *MatchSettingsHandler) Update(c *gin.Context) {
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
	var req dto.UpdateMatchSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	matchSettingsUPD := service.MatchSettingsUpdate{
		MapID:            req.MapID,
		PlaylistName:     req.PlaylistName,
		AdminChat:        req.AdminChat,
		TeamRename:       req.TeamRename,
		SelfAssign:       req.SelfAssign,
		AimAssist:        req.AimAssist,
		AnonMode:         req.AnonMode,
		DropSpotsEnabled: req.DropSpotsEnabled,
		FillBotsMode:     req.FillBotsMode,
	}
	err = h.matchSettingsSerivce.Update(c.Request.Context(), organisationID, actorID, matchID, matchSettingsUPD)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

// DELETE api/v1/organisations/:organisation_id/matches/:match_id/settings
func (h *MatchSettingsHandler) Delete(c *gin.Context) {
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
	err = h.matchSettingsSerivce.Delete(c.Request.Context(), organisationID, actorID, matchID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}
