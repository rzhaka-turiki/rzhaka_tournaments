package model

import (
	"time"

	"github.com/google/uuid"
)

type Match struct {
	ID          uuid.UUID
	ExternalMID *string

	MapID        uuid.UUID
	StatsTokenID int

	Status    string
	StartedAt time.Time

	Settings   MatchSettings
	Slots      []MatchSlot
	LegendBans []MatchLegendBan
}

type MatchSettings struct {
	PlaylistName string
	AdminChat    bool
	TeamRename   bool
	SelfAssign   bool
	AimAssist    bool
	AnonMode     bool
}

type MatchSlot struct {
	ID            uuid.UUID
	MatchID       uuid.UUID
	SlotNumber    int
	TeamID        *uuid.UUID
	MapLocationID *uuid.UUID

	Players []MatchSlotPlayer
}

type MatchSlotPlayer struct {
	ID              uuid.UUID
	SlotID          uuid.UUID
	UserID          *uuid.UUID
	ExpectedNIDHash *string
}
