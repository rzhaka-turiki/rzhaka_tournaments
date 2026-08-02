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

type CreateRoleRequest struct {
	Name      string `json:"name" binding:"required"`
	Position  int    `json:"position" binding:"required"`
	RoleColor string `json:"role_color"`
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
	RolesResponse := make([]RoleResponse, 0, len(roles))
	for _, role := range roles {
		RolesResponse = append(RolesResponse, FromRole(role))
	}
	return RolesResponse
}
