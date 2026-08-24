package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type BindApexAccountRequest struct {
	Player   string `json:"player" binding:"required"`
	Platform string `json:"platform" binding:"required"`
	Level    int32  `json:"level" binding:"required,min=1"`
}

type ApexAccountResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	UID       string    `json:"uid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromApexAccount(apexAccount *model.ApexAccount) *ApexAccountResponse {
	return &ApexAccountResponse{
		ID:        apexAccount.ID,
		UserID:    apexAccount.UserID,
		UID:       apexAccount.UID,
		CreatedAt: apexAccount.CreatedAt,
		UpdatedAt: apexAccount.UpdatedAt,
	}
}
