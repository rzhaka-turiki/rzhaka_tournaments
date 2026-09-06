package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type MatchResult struct {
	ID               uuid.UUID
	MatchID          *uuid.UUID
	ExternalMID      string
	MapID            *uuid.UUID
	MapName          string
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

	Shots     int
	Headshots int
	Hardware  string
	Hits      int
	Revives   int
	Respawns  int

	CreatedAt time.Time
}
