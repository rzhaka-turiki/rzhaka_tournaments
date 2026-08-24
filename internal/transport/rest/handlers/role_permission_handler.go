package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/service"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/auth"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/dto"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/response"
)

type RolePermissionHandler struct {
	rolePermissionService service.RolePermissionService
}

func NewRolePermissionHandler(rolePermissionSevice service.RolePermissionService) *RolePermissionHandler {
	return &RolePermissionHandler{
		rolePermissionService: rolePermissionSevice,
	}
}

func (h *RolePermissionHandler) GetRolePermissions(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ROLE_ID", "invalid role id")
		return
	}
	permissions, err := h.rolePermissionService.GetRolePermissions(c.Request.Context(), roleID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromPermissions(permissions))
}

func (h *RolePermissionHandler) AddPermission(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ROLE_ID", "invalid role id")
		return
	}
	var req dto.AddPermissionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	// change after JWT add
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	err = h.rolePermissionService.AddPermission(c.Request.Context(), actorID, roleID, req.PermissionID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

func (h *RolePermissionHandler) RemovePermission(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ROLE_ID", "invalid role id")
		return
	}
	permissionID, err := strconv.Atoi(c.Param("permissionID"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_PERMISSION_ID", "invalid permission id")
		return
	}
	// change after JWT add
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	err = h.rolePermissionService.RemovePermission(c.Request.Context(), actorID, roleID, permissionID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}
