package dto

import (
	"time"

	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/google/uuid"
)

type TeamInviteLink struct {
	ID        uuid.UUID
	TeamID    uuid.UUID
	Token     string
	CreatedBy uuid.UUID
	MaxUses   int
	UsedCount int
	ExpiresAt time.Time
	CreatedAt time.Time
}

type CreateInviteLinkRequest struct {
	MaxUses   int       `json:"max_uses" binding:"required,min=1"`
	ExpiresAt time.Time `json:"expires_at" binding:"required"`
}

type InviteLinkPreviewResponse struct {
	TeamID    uuid.UUID `json:"team_id"`
	TeamName  string    `json:"team_name"`
	ExpiresAt time.Time `json:"expires_at"`
}

type InviteLinkResponse struct {
	ID     uuid.UUID `json:"id"`
	TeamID uuid.UUID `json:"team_id"`
	Token  string    `json:"token"`

	MaxUses   int `json:"max_uses"`
	UsedCount int `json:"used_count"`

	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func FromInviteLink(link model.TeamInviteLink) InviteLinkResponse {
	return InviteLinkResponse{
		ID:        link.ID,
		TeamID:    link.TeamID,
		Token:     link.Token,
		MaxUses:   *link.MaxUses,
		UsedCount: link.UsedCount,
		ExpiresAt: link.ExpiresAt,
		CreatedAt: link.CreatedAt,
	}
}

func FromInviteLinks(links []model.TeamInviteLink) []InviteLinkResponse {
	result := make([]InviteLinkResponse, 0, len(links))

	for _, link := range links {
		result = append(result, FromInviteLink(link))
	}

	return result
}

func FromInviteLinkPreview(preview model.InviteLinkPreview) InviteLinkPreviewResponse {
	return InviteLinkPreviewResponse{
		TeamID:    preview.TeamID,
		TeamName:  preview.TeamName,
		ExpiresAt: preview.ExpiresAt,
	}
}

type CreateInviteRequest struct {
	UserID uuid.UUID `json:"user_id"`
}
