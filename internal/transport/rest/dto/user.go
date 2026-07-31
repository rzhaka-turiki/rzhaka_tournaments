package dto

import (
	"time"

	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func FromUser(user *model.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}
}
