package model

import (
	"time"

	"github.com/google/uuid"
)

type PlayerResult struct {
	ID         uuid.UUID
	User       User
	Character  []Legend
	Kills      int
	Assists    int
	Damages    int
	Headshots  int
	Revives    int
	Respawns   int
	TimeAlive  int
	Hits       int
	Shots      int
	Knockdowns int
	// Можно использовать как в одиночном матче, так и во всех сразу
}

type TeamResult struct {
	ID        uuid.UUID
	Team      Team
	Players   []PlayerResult
	Placement int
	Points    int
	Kills     int
	// Можно использовать как водиночном матче, так и во всех сразу
}

type MatchResult struct {
	ID        uuid.UUID
	Match     Match
	Teams     []TeamResult
	Map       string
	StartedAt time.Time
}

type GroupResult struct {
	ID        uuid.UUID
	Group     Group
	Teams     []TeamResult
	Matches   []MatchResult
	StartedAt time.Time
}

type StageResult struct {
	ID        uuid.UUID
	Stage     Stage
	Groups    []GroupResult
	StartedAt time.Time
}

type TournamentResult struct {
	ID         uuid.UUID
	Tournament Tournament
	Stages     []StageResult
	StartedAt  time.Time
}

type LeagueResult struct {
	ID          uuid.UUID
	League      League
	Tournaments []TournamentResult
	StartedAt   time.Time
}
