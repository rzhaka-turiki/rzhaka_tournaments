package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TournamentRegistration struct {
	ID         uuid.UUID
	Team       Team
	Players    []User
	CheckIn    []bool
	CreatedAt  time.Time
	InWaitList bool
}

type RegistrationEvent struct {
	ID           uuid.UUID
	Actioner     User
	Registration TournamentRegistration
	Type         string
	Body         json.RawMessage
}
