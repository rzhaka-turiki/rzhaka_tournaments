package model

import (
	"time"

	"github.com/google/uuid"
)

type Ban struct {
	ID        uuid.UUID
	User      User
	Type      string
	Body      string
	Duration  time.Duration
	CreatedAt time.Time
}

type Suspend struct {
	ID        uuid.UUID
	User      User
	Type      string
	Body      string
	Duration  time.Duration
	CreatedAt time.Time
}
