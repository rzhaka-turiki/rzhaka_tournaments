package handlers

import (
	"net/http"

	"github.com/a1uka/rzhaka_tournaments/internal/service"
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
		RespondError(c, http.StatusBadRequest, "invalid uuid")
		return
	}
	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		RespondError(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}
