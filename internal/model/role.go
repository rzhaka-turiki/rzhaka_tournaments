package model

type Role struct {
	ID        int
	Name      string
	RoleColor string
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
