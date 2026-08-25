package model

import "github.com/google/uuid"

type Legend struct {
	ID       uuid.UUID
	Name     string
	ImageURL string
	Class    string
	Ability  string
	Ultimate string
}

type MatchLegendBan struct {
	MatchID  uuid.UUID
	LegendID uuid.UUID
}
