package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type CreateMatchSlotPlayerRequest struct {
	UserID          *uuid.UUID `json:"user_id"`
	ExpectedNIDHash *string    `json:"expected_nid_hash"`
}

type MatchSlotPlayerResponse struct {
	ID              uuid.UUID  `json:"id"`
	MatchSlotID     uuid.UUID  `json:"match_slot_id"`
	UserID          *uuid.UUID `json:"user_id"`
	ExpectedNIDHash *string    `json:"expected_nid_hash"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type UpdateMatchSlotPlayerRequest struct {
	UserID          *uuid.UUID `json:"user_id"`
	ExpectedNIDHash *string    `json:"expected_nid_hash"`
}

func FromMatchSlotPlayer(slotPlayer *model.MatchSlotPlayer) MatchSlotPlayerResponse {
	return MatchSlotPlayerResponse{
		ID:              slotPlayer.ID,
		MatchSlotID:     slotPlayer.MatchSlotID,
		UserID:          slotPlayer.UserID,
		ExpectedNIDHash: slotPlayer.ExpectedNIDHash,
		CreatedAt:       slotPlayer.CreatedAt,
		UpdatedAt:       slotPlayer.UpdatedAt,
	}
}

func FromMatchSlotPlayers(slotPlayers []model.MatchSlotPlayer) []MatchSlotPlayerResponse {
	slotPlayersResponse := make([]MatchSlotPlayerResponse, 0, len(slotPlayers))
	for _, player := range slotPlayers {
		slotPlayersResponse = append(slotPlayersResponse, FromMatchSlotPlayer(&player))
	}
	return slotPlayersResponse
}
