package model

import (
	"time"

	"github.com/google/uuid"
)

type Map struct {
	ID                uuid.UUID
	Name              string
	InGameName        string
	ImageURL          string
	MinimapImageURL   string
	SupportsDropSpots bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type MapLocation struct {
	ID       uuid.UUID
	MapID    uuid.UUID
	Name     string
	ImageURL string
	Position int
}
