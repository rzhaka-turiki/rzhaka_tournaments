package model

import (
	"time"

	"github.com/google/uuid"
)

type Legend struct {
	ID              uuid.UUID
	Name            string
	InGameName      string
	ImageURL        string
	ProfileImageURL string
	Class           string
	Ability         string
	Ultimate        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type MatchLegendBan struct {
	ID        uuid.UUID
	MatchID   uuid.UUID
	LegendID  uuid.UUID
	CreatedAt time.Time
}

// I'll leave it as it is, but wont use it for now cause its not so important
type LegendSnapshot struct {
	ID              uuid.UUID
	LegendID        uuid.UUID
	Name            string
	InGameName      string
	ImageURL        string
	ProfileImageURL string
	Class           string
	Ability         string
	Ultimate        string
}
