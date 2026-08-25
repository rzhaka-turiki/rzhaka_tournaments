package model

import (
	"time"

	"github.com/google/uuid"
)

type TeamInviteLink struct {
	ID        uuid.UUID
	TeamID    uuid.UUID
	Token     string
	CreatedBy uuid.UUID
	MaxUses   *int
	UsedCount int
	ExpiresAt time.Time
	CreatedAt time.Time
}

type InviteLinkPreview struct {
	TeamID    uuid.UUID
	TeamName  string
	ExpiresAt time.Time
}
