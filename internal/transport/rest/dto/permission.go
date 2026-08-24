package dto

import "github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"

type PermissionResponse struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

func FromPermission(permission model.Permission) PermissionResponse {
	return PermissionResponse{
		ID:          permission.ID,
		Code:        permission.Code,
		Description: permission.Description,
	}
}

func FromPermissions(permissions []model.Permission) []PermissionResponse {
	var PermissionsResponse []PermissionResponse
	for _, permission := range permissions {
		PermissionsResponse = append(PermissionsResponse, FromPermission(permission))
	}
	return PermissionsResponse
}

type AddPermissionRequest struct {
	PermissionID int `json:"permission_id" binding:"required"`
}
