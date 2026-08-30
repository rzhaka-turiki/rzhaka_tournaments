package model

import "github.com/google/uuid"

type Map struct {
	ID                uuid.UUID
	Name              string
	InGameName        string
	ImageURL          string
	MinimapImageURL   string
	SupportsDropSpots bool
}

type MapLocation struct {
	ID       uuid.UUID
	MapID    uuid.UUID
	Name     string
	ImageURL string
	Position int
}
