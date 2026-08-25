package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type CreateTeamInvitationRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

type TeamInvitationResponse struct {
	ID        uuid.UUID `json:"id"`
	TeamID    uuid.UUID `json:"team_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedBy uuid.UUID `json:"created_by"`

	Type string `json:"type"`

	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func FromTeamInvitation(req model.TeamRequest) TeamInvitationResponse {
	return TeamInvitationResponse{
		ID:        req.ID,
		TeamID:    req.TeamID,
		UserID:    req.UserID,
		CreatedBy: req.CreatedBy,
		Type:      string(req.Type),
		ExpiresAt: req.ExpiresAt,
		CreatedAt: req.CreatedAt,
	}
}

func FromTeamInvitations(requests []model.TeamRequest) []TeamInvitationResponse {

	result := make([]TeamInvitationResponse, 0, len(requests))

	for _, req := range requests {
		result = append(result, FromTeamInvitation(req))
	}

	return result
}
