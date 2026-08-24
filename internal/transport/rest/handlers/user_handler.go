package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/service"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/auth"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/dto"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/response"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}
	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			HandleError(c, err)
			return
		}
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromUser(user))
}

func (h *UserHandler) Me(c *gin.Context) {
	actorID := auth.UserID(c)
	me, err := h.userService.GetMe(c.Request.Context(), actorID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, me)
}
