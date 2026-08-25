package model

import (
	"time"

	"github.com/google/uuid"
)

type MatchResult struct {
	ID        uuid.UUID
	MatchID   *uuid.UUID
	MapID     uuid.UUID
	StartedAt time.Time

	Slots []SlotMatchResult
}

type SlotMatchResult struct {
	ID       uuid.UUID
	ResultID uuid.UUID
	SlotID   uuid.UUID

	Placement int
	Points    int
	Kills     int

	SourceTeamName string

	Players []PlayerMatchResult
}

type PlayerMatchResult struct {
	ID           uuid.UUID
	SlotResultID uuid.UUID

	NIDHash  string
	Gamertag string
	Legend   string

	Kills      int
	Assists    int
	Damage     int
	Headshots  int
	Revives    int
	Respawns   int
	TimeAlive  int
	Hits       int
	Shots      int
	Knockdowns int
}
