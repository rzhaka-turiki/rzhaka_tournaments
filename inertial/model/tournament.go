package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type DropSpot struct {
	ID      int
	Name    string
	DraftID int
	Spot    int
}

type Legend struct {
	ID       uuid.UUID
	Name     string
	ImageURL string
	Class    string
	Ability  string
	Ultimate string
}

type Match struct {
	ID          uuid.UUID
	Map         string
	Status      string
	StartAt     time.Time
	Drafts      []DropSpot
	StatsCodeID int
	// Idk should we store teams here or not
}

type Group struct {
	ID            uuid.UUID
	Format        string
	Status        string
	StartAt       time.Time
	FinishAt      time.Time
	Matches       []Match
	BannedLegends []Legend
	NextGroupInfo json.RawMessage
	PrevGroupInfo json.RawMessage
	// Idk should we store teams here or not
}

type Stage struct {
	ID            uuid.UUID
	Format        string
	Status        string
	StartAt       time.Time
	FinishAt      time.Time
	Groups        []Group
	NextStageInfo json.RawMessage
	PrevStageInfo json.RawMessage
	// Idk should we store teams here or not
}

type Tournament struct {
	ID                   uuid.UUID
	Format               string
	StartAt              time.Time
	FinishAt             time.Time
	Status               string
	RegistrationStartAt  time.Time
	RegistrationFinishAt time.Time
	Info                 json.RawMessage
	CheckInStartAt       time.Time
	CheckInFinishAt      time.Time
	Stages               []Stage
	NextTournamentInfo   json.RawMessage
	PrevTournamentInfo   json.RawMessage
	// Idk should we store teams here or not
}

type League struct {
	ID                   uuid.UUID
	Format               string
	StartAt              time.Time
	FinishAt             time.Time
	Status               string
	RegistrationStartAt  time.Time
	RegistrationFinishAt time.Time
	Info                 json.RawMessage
	CheckInStartAt       time.Time
	CheckInFinishAt      time.Time
	Tournaments          []Tournament
	NextLeagueInfo       json.RawMessage
	PrevLeagueInfo       json.RawMessage
	// Idk should we store teams here or not
}
