package model

import (
	"time"

	"github.com/google/uuid"
)

type ApexAccount struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	UID       string    `json:"uid" db:"uid"`
	NIDHash   string    `json:"nid_hash" db:"nid_hash"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
