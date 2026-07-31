package model

import "time"

type Role struct {
	ID        int
	Name      string
	RoleColor string
	// higher pos - higher role
	Position int
	// cant be deleted
	IsSystem  bool
	CreatedAt time.Time
}

type Permission struct {
	ID          int
	Code        string
	Description string
}

type RolePermission struct {
	Role       Role
	Permission Permission
}

type UserRole struct {
	User User
	Role Role
}
