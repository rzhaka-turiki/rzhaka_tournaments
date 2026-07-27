package model

import (
	"time"

	"github.com/google/uuid"
)

type TeamParticipant struct {
	User     uuid.UUID
	Role     string
	JoinedAt time.Time
}

type Team struct {
	ID        uuid.UUID
	Name      string
	ShortName string
	LogoURL   string
	Members   []TeamParticipant
}
