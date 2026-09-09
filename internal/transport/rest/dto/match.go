package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type CreateMatchRequest struct {
	MapID          *uuid.UUID `json:"map_id"`
	GroupID        uuid.UUID  `json:"group_id" binding:"required"`
	OrganisationID uuid.UUID  `json:"organisation_id" binding:"required"`
	StatsTokenID   *uuid.UUID `json:"stats_token_id"`
	StartAt        *time.Time `json:"start_at"`
}

type MatchResponse struct {
	MatchID        uuid.UUID         `json:"id" binding:"required"`
	MapID          uuid.UUID         `json:"map_id" binding:"required"`
	GroupID        uuid.UUID         `json:"group_id" binding:"required"`
	OrganisationID uuid.UUID         `json:"organisation_id" binding:"required"`
	StatsTokenID   *uuid.UUID        `json:"stats_token_id"`
	Status         model.MatchStatus `json:"status" binding:"required"`
	StartAt        time.Time         `json:"start_at" binding:"required"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

func FromMatch(match *model.Match) MatchResponse {
	return MatchResponse{
		MatchID:        match.ID,
		MapID:          match.MapID,
		GroupID:        match.GroupID,
		OrganisationID: match.OrganisationID,
		StatsTokenID:   match.StatsTokenID,
		Status:         match.Status,
		StartAt:        *match.StartAt,
		CreatedAt:      match.CreatedAt,
		UpdatedAt:      match.UpdatedAt,
	}
}

func FromMatches(matches []model.Match) []MatchResponse {
	MatchesResponse := make([]MatchResponse, 0, len(matches))
	for _, match := range matches {
		MatchesResponse = append(MatchesResponse, FromMatch(&match))
	}
	return MatchesResponse
}
