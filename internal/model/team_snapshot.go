package model

import (
	"time"

	"github.com/google/uuid"
)

type TeamSnapshot struct {
	ID           uuid.UUID
	TeamID       uuid.UUID
	Name         string
	ShortName    string
	LogoPath     string
	LogoDarkPath string
	CreatedAt    time.Time
}
