package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TeamParticipant struct {
	User     User
	Role     string
	JoinedAt time.Time
}

type Team struct {
	ID       uuid.UUID
	Name     string
	ImageURL string
	Members  []TeamParticipant
}

type TeamEvent struct {
	ID        uuid.UUID
	Team      Team
	Actioner  User
	EventType string
	Body      json.RawMessage
	CreatedAt time.Time
}
