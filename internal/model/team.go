package model

import (
	"time"

	"github.com/google/uuid"
)

type TeamMember struct {
	UserID   uuid.UUID
	TeamID   uuid.UUID
	Role     string
	JoinedAt time.Time
}

type Team struct {
	ID           uuid.UUID
	Name         string
	ShortName    string
	LogoPath     string
	LogoDarkPath string
	OwnerID      uuid.UUID
	CreatedAt    time.Time
	ArchivedAt   *time.Time
}
