package dto

import "github.com/a1uka/rzhaka_tournaments/internal/model"

type RoleResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	Position int    `json:"position"`
}

type AssignRoleRequest struct {
	RoleID int `json:"role_id" binding:"required"`
}

func FromRole(role model.Role) RoleResponse {
	return RoleResponse{
		ID:       role.ID,
		Name:     role.Name,
		Color:    role.RoleColor,
		Position: role.Position,
	}
}

func FromRoles(roles []model.Role) []RoleResponse {
	var RolesResponse []RoleResponse
	for _, role := range roles {
		RolesResponse = append(RolesResponse, FromRole(role))
	}
	return RolesResponse
}
