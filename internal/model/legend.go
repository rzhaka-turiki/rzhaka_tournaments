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
}

type MatchLegendBan struct {
	ID        uuid.UUID
	MatchID   uuid.UUID
	LegendID  uuid.UUID
	CreatedAt time.Time
}

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
