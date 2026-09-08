package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/service"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/dto"
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
	var req dto.CreateTokenRequest
}
