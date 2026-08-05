package model

import (
	"time"

	"github.com/google/uuid"
)

type TeamRequestType string

const (
	TeamRequestInvite     TeamRequestType = "invite"
	TeamRequestJoin       TeamRequestType = "join_request"
	TeamRequestInviteLink TeamRequestType = "invite_link"
)

type TeamRequest struct {
	ID        uuid.UUID
	TeamID    uuid.UUID
	UserID    uuid.UUID
	CreatedBy uuid.UUID
	Type      TeamRequestType
	ExpiresAt time.Time
	CreatedAt time.Time
}
