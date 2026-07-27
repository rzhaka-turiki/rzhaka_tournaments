package model

import (
	"time"

	"github.com/google/uuid"
)

type TournamentRegistration struct {
	ID         uuid.UUID
	TeamID     uuid.UUID
	PlayerIDs  []uuid.UUID
	CheckIn    []bool
	CreatedAt  time.Time
	InWaitList bool
}
