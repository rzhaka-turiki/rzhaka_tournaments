package model

import (
	"time"

	"github.com/google/uuid"
)

type MatchStatus string

const (
	MatchStatusPending   MatchStatus = "pending"
	MatchStatusReady     MatchStatus = "ready"
	MatchStatusRunning   MatchStatus = "running"
	MatchStatusFinished  MatchStatus = "finished"
	MatchStatusCancelled MatchStatus = "cancelled"
)

type Match struct {
	ID           uuid.UUID
	MapID        uuid.UUID
	GroupID      uuid.UUID
	StatsTokenID *uuid.UUID
	Status       MatchStatus
	StartAt      *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type MatchSettings struct {
	MatchID          uuid.UUID
	PlaylistName     string
	AdminChat        bool
	TeamRename       bool
	SelfAssign       bool
	AimAssist        bool
	AnonMode         bool
	DropSpotsEnabled bool
}

type MatchSlot struct {
	ID         uuid.UUID
	MatchID    uuid.UUID
	SlotNumber int
	DropSpotID *uuid.UUID
}

type MatchSlotPlayer struct {
	ID              uuid.UUID
	MatchSlotID     uuid.UUID
	UserID          *uuid.UUID
	ExpectedNIDHash *string
}

type MatchAPIToken struct {
	ID              uuid.UUID
	MatchAPITokenID int

	AddedBy        uuid.UUID
	OrganisationID *uuid.UUID
	Activation     time.Time
	Expiration     time.Time

	StatsToken  string
	AdminToken  *string
	PlayerToken *string

	CreatedAt time.Time
	UpdatedAt time.Time
}
