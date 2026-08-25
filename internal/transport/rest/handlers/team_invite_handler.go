package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/service"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/auth"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/dto"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/response"
)

type TeamRequestService interface {
	CreateInvite(ctx context.Context, actorID, teamID, userID uuid.UUID) (*model.TeamRequest, error)
	GetUserRequests(ctx context.Context, userID uuid.UUID) ([]model.TeamRequest, error)
	GetTeamRequests(ctx context.Context, actorID, teamID uuid.UUID) ([]model.TeamRequest, error)
	AcceptRequest(ctx context.Context, userID, requestID uuid.UUID) error
	RejectRequest(ctx context.Context, userID, requestID uuid.UUID) error
	DeleteRequest(ctx context.Context, actorID, requestID uuid.UUID) error
}

type TeamRequestHandler struct {
	requestService service.TeamRequestService
}

func NewTeamRequestHandler(requestService service.TeamRequestService) *TeamRequestHandler {
	return &TeamRequestHandler{
		requestService: requestService,
	}
}

func (h *TeamRequestHandler) CreateInvite(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_TEAM_ID", "invalid team id")
		return
	}

	var req dto.CreateTeamInvitationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}

	actorID := auth.UserID(c)

	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}

	err = h.requestService.CreateInvite(c.Request.Context(), teamID, actorID, req.UserID)

	if err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, nil)
}

func (h *TeamRequestHandler) GetMyInvitations(c *gin.Context) {
	userID := auth.UserID(c)

	if userID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}

	requests, err := h.requestService.GetUserInvitations(c.Request.Context(), userID)

	if err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, dto.FromTeamInvitations(requests))
}

func (h *TeamRequestHandler) GetTeamInvitations(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_TEAM_ID", "invalid team id")
		return
	}

	actorID := auth.UserID(c)

	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}

	requests, err := h.requestService.GetTeamInvitations(c.Request.Context(), teamID, actorID)

	if err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, dto.FromTeamInvitations(requests))
}

func (h *TeamRequestHandler) AcceptRequest(c *gin.Context) {
	requestID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_REQUEST_ID", "invalid request id")
		return
	}

	actorID := auth.UserID(c)

	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}

	err = h.requestService.AcceptRequest(c.Request.Context(), requestID, actorID)

	if err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, http.StatusNoContent, nil)
}

func (h *TeamRequestHandler) RejectRequest(c *gin.Context) {
	requestID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_REQUEST_ID", "invalid request id")
		return
	}

	actorID := auth.UserID(c)

	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}

	err = h.requestService.RejectRequest(c.Request.Context(), requestID, actorID)

	if err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, http.StatusNoContent, nil)
}
