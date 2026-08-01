package model

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	RoleColor string `json:"role_color"`
	// higher pos - higher role
	Position int `json:"position"`
	// cant be deleted
	IsSystem  bool       `json:"system"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	DeletedBy *uuid.UUID `json:"deleted_by,omitempty"`
}

type Permission struct {
	ID          int       `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type RolePermission struct {
	Role       Role
	Permission Permission
}

type UserRole struct {
	User User
	Role Role
}
