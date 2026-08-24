package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/service"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/auth"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/dto"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/response"
)

type ApexAccountHandler struct {
	apexAccountService service.ApexAccountService
}

func NewApexAccountHandler(apexAccountService service.ApexAccountService) *ApexAccountHandler {
	return &ApexAccountHandler{
		apexAccountService: apexAccountService,
	}
}

func (h *ApexAccountHandler) GetMyAccount(c *gin.Context) {
	userID := auth.UserID(c)
	account, err := h.apexAccountService.GetByUserID(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		if errors.Is(err, service.ErrApexAccountNotFound) {
			HandleError(c, err)
			return
		}
	}
	response.Success(c, http.StatusOK, dto.FromApexAccount(account))
}

func (h *ApexAccountHandler) Bind(c *gin.Context) {
	var req dto.BindApexAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}
	userID := auth.UserID(c)
	account, err := h.apexAccountService.Bind(
		c.Request.Context(),
		userID,
		req.Player,
		req.Platform,
		req.Level,
	)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, dto.FromApexAccount(account))
}
