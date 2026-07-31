package handlers

import (
	"net/http"
	"strconv"

	"github.com/a1uka/rzhaka_tournaments/internal/service"
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/dto"
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoleHandler struct {
	roleService service.RoleService
}

func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_UUID", "invalid uuid")
		return
	}
	roles, err := h.roleService.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}
	response.Success(c, http.StatusOK, dto.FromRoles(roles))
}

func (h *RoleHandler) GetUserPermissions(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_UUID", "invalid uuid")
		return
	}
	permissions, err := h.roleService.GetUserPermissions(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}
	response.Success(c, http.StatusOK, dto.FromPermissions(permissions))
}

func (h *RoleHandler) AssignRole(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_UUID", "invalid uuid")
		return
	}
	var req dto.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	actorID := userID
	err = h.roleService.AssignRole(c.Request.Context(), actorID, userID, req.RoleID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

func (h *RoleHandler) RemoveRole(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_UUID", "invalid uuid")
		return
	}
	roleID, err := strconv.Atoi(c.Param("roleID"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ROLE_ID", "invalid role id")
		return
	}
	actorID := userID
	err = h.roleService.RemoveRole(c.Request.Context(), actorID, userID, roleID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}
