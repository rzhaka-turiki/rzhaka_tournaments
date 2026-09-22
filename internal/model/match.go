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

// done
type Match struct {
	ID             uuid.UUID
	MapID          uuid.UUID
	OrganisationID uuid.UUID
	StatsTokenID   *uuid.UUID
	Status         MatchStatus
	StartAt        *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// done
type MatchSettings struct {
	MatchID          uuid.UUID
	MapID            uuid.UUID
	PlaylistName     string
	MapName          string
	AdminChat        bool
	TeamRename       bool
	SelfAssign       bool
	AimAssist        bool
	AnonMode         bool
	DropSpotsEnabled bool
	FillBotsMode     bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// done
type MatchSlot struct {
	ID         uuid.UUID
	MatchID    uuid.UUID
	SlotNumber int
	DropSpotID *uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Done
type MatchSlotPlayer struct {
	ID              uuid.UUID
	MatchSlotID     uuid.UUID
	UserID          *uuid.UUID
	ExpectedNIDHash *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// done
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
