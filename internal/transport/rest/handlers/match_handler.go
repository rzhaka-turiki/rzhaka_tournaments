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

type MatchHandler struct {
	matchService service.MatchService
}

func NewMatchHandler(matchService service.MatchService) *MatchHandler {
	return &MatchHandler{
		matchService: matchService,
	}
}

func (h *MatchHandler) Create(c *gin.Context) {
	organisationID, err := uuid.Parse(c.Param("organisation_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ORGANISATION_ID", "invalid organisation id")
		return
	}
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	var req dto.CreateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	match := &model.Match{
		MapID:          *req.MapID,
		StatsTokenID:   req.StatsTokenID,
		OrganisationID: organisationID,
		GroupID:        req.GroupID,
		StartAt:        req.StartAt,
	}
	err = h.matchService.Create(c.Request.Context(), actorID, match)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, dto.FromMatch(match))
}

func (h *MatchHandler) GetByID(c *gin.Context) {
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	// Needs to add various return info. F.e., common user doesn't need to know about tokenID, time created e.t.c.
	matchID, err := uuid.Parse(c.Param("match_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_TOKEN_ID", "invalid token id")
		return
	}
	match, err := h.matchService.GetByID(c.Request.Context(), matchID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromMatch(match))
}

func (h *MatchHandler) Update(c *gin.Context) {
	organisationID, err := uuid.Parse(c.Param("organisation_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ORGANISATION_ID", "invalid organisation id")
		return
	}
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	matchID, err := uuid.Parse(c.Param("match_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_MATCH_ID", "invalid match id")
		return
	}

}
