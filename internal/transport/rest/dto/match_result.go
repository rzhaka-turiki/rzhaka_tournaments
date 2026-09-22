package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type CreateMatchResultRequest struct {
	ExternalMID      string     `json:"external_mid" binding:"required"`
	MapID            *uuid.UUID `json:"map_id"`
	MapName          string     `json:"map_name" binding:"required"`
	AimAssistAllowed bool       `json:"aim_assist_allowed"`
}

type MatchResultResponse struct {
	MatchID          uuid.UUID  `json:"match_id"`
	ExternalMID      string     `json:"external_mid"`
	MapID            *uuid.UUID `json:"map_id"`
	MapName          string     `json:"map_name"`
	AimAssistAllowed bool       `json:"aim_assist_allowed"`
	StartedAt        time.Time  `json:"started_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type UpdateMatchResultRequest struct {
	ExternalMID      *string    `json:"external_mid"`
	MapID            *uuid.UUID `json:"map_id"`
	MapName          *string    `json:"map_name"`
	AimAssistAllowed *bool      `json:"aim_assist_allowed"`
	StartedAt        *time.Time `json:"started_at"`
}

func FromMatchResult(matchResult *model.MatchResult) MatchResultResponse {
	return MatchResultResponse{
		MatchID:          matchResult.MatchID,
		ExternalMID:      matchResult.ExternalMID,
		MapID:            matchResult.MapID,
		MapName:          matchResult.MapName,
		AimAssistAllowed: matchResult.AimAssistAllowed,
		CreatedAt:        matchResult.CreatedAt,
		UpdatedAt:        matchResult.UpdatedAt,
		StartedAt:        matchResult.StartedAt,
	}
}
