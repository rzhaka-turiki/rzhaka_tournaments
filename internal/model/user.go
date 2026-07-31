package model

import (
	"time"

	"net/netip"

	"encoding/json"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	DiscordIDEncrypted string    `json:"-" db:"discord_id_encrypted"`
	DiscordIDHash      []byte    `json:"-" db:"discord_id_hash"`
	Username           string    `json:"username" db:"username"`
	AvatarURL          *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	Status             string    `json:"status" db:"status"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

type LoginEntry struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	LoggedInAt  time.Time  `json:"logged_in_at"`
	IPAddress   netip.Addr `json:"ip_address,omitempty"`
	UserAgent   string     `json:"user_agent,omitempty"`
	LoginMethod string     `json:"login_method"`
	Success     bool       `json:"success"`
}

type Event struct {
	ID             uuid.UUID       `json:"id"`
	ActorID        *uuid.UUID      `json:"actor_id,omitempty"` // кто совершил
	UserID         *uuid.UUID      `json:"user_id,omitempty"`
	TeamID         *uuid.UUID      `json:"team_id,omitempty"`
	TournamentID   *uuid.UUID      `json:"tournament_id,omitempty"`
	RegistrationID *uuid.UUID      `json:"registration_id,omitempty"`
	StageID        *uuid.UUID      `json:"stage_id,omitempty"`
	MatchID        *uuid.UUID      `json:"match_id,omitempty"`
	EventType      string          `json:"event_type"`
	Payload        json.RawMessage `json:"payload"`
	CreatedAt      time.Time       `json:"created_at"`
}

type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TokenHash []byte    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
