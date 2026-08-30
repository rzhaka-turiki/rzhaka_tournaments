package model

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/type/decimal"
)

type MatchResult struct {
	ID               uuid.UUID
	MatchID          *uuid.UUID
	ExternalMID      string
	MapID            *uuid.UUID
	StartedAt        time.Time
	AimAssistAllowed bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type MatchResultSlot struct {
	ID            uuid.UUID
	MatchResultID uuid.UUID
	SlotNumber    int

	TeamName      *string
	TeamPlacement *int
	Points        decimal.Decimal
	Kills         int

	MatchSlotID *uuid.UUID

	CreatedAt time.Time
}

type MatchResultPlayer struct {
	ID                uuid.UUID
	MatchResultSlotID uuid.UUID

	NIDHash    string
	PlayerName string

	LegendID      *uuid.UUID
	CharacterName string

	Kills        int
	Assists      int
	Knockdowns   int
	DamageDealt  int
	SurvivalTime int

	CreatedAt time.Time
}
