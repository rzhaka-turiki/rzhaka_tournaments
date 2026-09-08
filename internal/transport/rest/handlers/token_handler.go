package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/service"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/auth"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/dto"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/response"
)

type TokenHandler struct {
	tokenService service.TokenService
}

func NewTokenHandler(tokenService service.TokenService) *TokenHandler {
	return &TokenHandler{
		tokenService: tokenService,
	}
}

func (h *TokenHandler) Create(c *gin.Context) {
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
	var req dto.CreateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	token := &model.MatchAPIToken{
		AdminToken:  req.AdminToken,
		StatsToken:  req.StatsToken,
		PlayerToken: req.PlayerToken,
		Activation:  req.Activation,
		Expiration:  req.Expiration,
	}
	err = h.tokenService.Create(c.Request.Context(), organisationID, actorID, token)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, dto.FromToken(*token))
}

func (h *TokenHandler) GetByOrganisationID(c *gin.Context) {
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
	includeInactive := false
	if value := c.Query("include_inactive"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			response.Fail(
				c,
				http.StatusBadRequest,
				"INVALID_INCLUDE_INACTIVE",
				"invalid include_inactive",
			)
			return
		}

		includeInactive = parsed
	}
	tokens, err := h.tokenService.GetByOrganisationID(c.Request.Context(), organisationID, actorID, includeInactive)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromTokens(tokens))
}

func (h *TokenHandler) GetByID(c *gin.Context) {
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
	tokenID, err := uuid.Parse(c.Param("token_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_TOKEN_ID", "invalid token id")
		return
	}
	token, err := h.tokenService.GetByID(c.Request.Context(), organisationID, actorID, tokenID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromToken(*token))
}

func (h *TokenHandler) Update(c *gin.Context) {
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
	tokenID, err := uuid.Parse(c.Param("token_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_TOKEN_ID", "invalid token id")
		return
	}
	var req dto.UpdateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	tokenUPD := service.TokenUpdate{
		Activation:  req.Activation,
		Expiration:  req.Expiration,
		AdminToken:  req.AdminToken,
		PlayerToken: req.PlayerToken,
		StatsToken:  req.StatsToken,
	}
	err = h.tokenService.Update(c.Request.Context(), organisationID, actorID, tokenID, tokenUPD)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

func (h *TokenHandler) Delete(c *gin.Context) {
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
	tokenID, err := uuid.Parse(c.Param("token_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_TOKEN_ID", "invalid token id")
		return
	}
	err = h.tokenService.Delete(c.Request.Context(), organisationID, actorID, tokenID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}
