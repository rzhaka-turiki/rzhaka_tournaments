package handlers

import (
	"net/http"
	"strconv"

	"github.com/a1uka/rzhaka_tournaments/internal/service"
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/dto"
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/response"
	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permissionService service.PermissionService
}

func NewPermissionHandler(permissionService service.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
	}
}

func (h *PermissionHandler) GetAll(c *gin.Context) {
	permissions, err := h.permissionService.GetAll(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromPermissions(permissions))
}

func (h *PermissionHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}
	permission, err := h.permissionService.GetByID(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromPermission(*permission))
}

func (h *PermissionHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")
	permission, err := h.permissionService.GetByCode(c.Request.Context(), code)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromPermission(*permission))
}
