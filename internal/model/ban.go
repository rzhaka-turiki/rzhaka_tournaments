package model

import (
	"time"

	"github.com/google/uuid"
)

type Ban struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Type      string // "ban", "suspend"
	Body      string
	Duration  time.Duration
	CreatedAt time.Time
}
