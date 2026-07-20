package model

import (
	"encoding/json"
	"time"
)

type Role struct {
	ID   int
	Name string
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

type RoleEvent struct {
	ID        int
	Role      Role
	User      User
	EventType string
	EventBody json.RawMessage
	CreatedAt time.Time
}
