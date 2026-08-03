package handlers

import (
	"net/http"
	"strconv"

	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/a1uka/rzhaka_tournaments/internal/service"
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/auth"
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
		HandleError(c, err)
		return
	}
	roles, err := h.roleService.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromRoles(roles))
}

func (h *RoleHandler) GetUserPermissions(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}
	permissions, err := h.roleService.GetUserPermissions(c.Request.Context(), userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromPermissions(permissions))
}

func (h *RoleHandler) AssignRole(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}
	var req dto.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}
	// change after JWT add
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	err = h.roleService.AssignRole(c.Request.Context(), actorID, userID, req.RoleID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

func (h *RoleHandler) RemoveRole(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}
	roleID, err := strconv.Atoi(c.Param("roleID"))
	if err != nil {
		HandleError(c, err)
		return
	}
	// change after JWT add
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	err = h.roleService.RemoveRole(c.Request.Context(), actorID, userID, roleID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.roleService.GetAllRoles(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromRoles(roles))
}

func (h *RoleHandler) GetAllPermissions(c *gin.Context) {
	permissions, err := h.roleService.GetAllPermissions(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromPermissions(permissions))
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req dto.CreateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}
	// change after JWT add
	actorID := uuid.MustParse("test_val")
	role := &model.Role{
		Name:      req.Name,
		Position:  req.Position,
		RoleColor: req.RoleColor,
	}

	err := h.roleService.CreateRole(c.Request.Context(), actorID, role)

	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, dto.FromRole(*role))
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ID", "invalid role id")
		return
	}
	// change after JWT add
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	err = h.roleService.DeleteRole(c.Request.Context(), actorID, roleID)
	if err != nil {
		HandleError(c, err)
		return
	}
}

func (h *RoleHandler) RestoreRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ID", "invalid role id")
		return
	}
	// change after JWT add
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	err = h.roleService.RestoreRole(c.Request.Context(), actorID, roleID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
