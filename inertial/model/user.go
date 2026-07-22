package model

import (
	"time"

	"net/netip"

	"encoding/json"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID `db:"id"`
	DiscordIDEncrypted string    `db:"discord_id_encrypted"`
	DiscordIDHash      []byte    `db:"discord_id_hash"`
	Username           string    `db:"username"`
	AvatarURL          *string   `db:"avatar_url"`
	CreatedAt          time.Time `db:"created_at"`
}

type LoginEvent struct {
	ID         uuid.UUID
	User       User
	LoggenInAt time.Time
	IPAdress   netip.Addr
	UserAgent  string
}

type UserEvent struct {
	ID         uuid.UUID
	User       User
	HappenedAt time.Time
	EventType  string
	EventBody  json.RawMessage
}

type RefreshToken struct {
	ID        uuid.UUID
	User      User
	TokenHash string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Event struct {
	ID         uuid.UUID
	User       User
	HappenedAt time.Time
	EventType  string
	EventBody  json.RawMessage
}
