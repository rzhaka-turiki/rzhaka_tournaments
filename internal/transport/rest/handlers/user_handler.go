package handlers

import (
	"errors"
	"net/http"

	"github.com/a1uka/rzhaka_tournaments/internal/service"
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/dto"
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		response.Fail(c, http.StatusBadRequest, "INVALID_UUID", "invalid uuid")
		return
	}
	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Fail(c, http.StatusNotFound, "USER_NOT_FOUND", "user not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}
	response.Success(c, http.StatusOK, dto.FromUser(user))
}
