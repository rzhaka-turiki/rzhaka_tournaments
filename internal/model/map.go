package model

import "github.com/google/uuid"

type Map struct {
	ID                uuid.UUID
	Name              string
	InGameName        string
	ImageURL          string
	SupportsDropSpots bool
	Locations         []MapLocation
}

type MapLocation struct {
	ID          uuid.UUID
	MapID       uuid.UUID
	Name        string
	ImageURL    string
	MapPosition int
}
