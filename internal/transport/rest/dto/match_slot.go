package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type CreateMatchSlotRequest struct {
	SlotNumber int        `json:"slot_number" binding:"required"`
	DropSpotID *uuid.UUID `json:"drop_spot_id"`
}

type MatchSlotResponse struct {
	ID         uuid.UUID `json:"id"`
	MatchID    uuid.UUID `json:"match_id"`
	SlotNumber int       `json:"slot_number"`
	DropSpotID uuid.UUID `json:"drop_spot_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UpdateMatchSlotRequest struct {
	SlotNumber *int       `json:"slot_number"`
	DropSpotID *uuid.UUID `json:"drop_spot_id"`
}

func FromMatchSlot(slot *model.MatchSlot) MatchSlotResponse {
	return MatchSlotResponse{
		ID:         slot.ID,
		MatchID:    slot.MatchID,
		SlotNumber: slot.SlotNumber,
		DropSpotID: *slot.DropSpotID,
		CreatedAt:  slot.CreatedAt,
		UpdatedAt:  slot.UpdatedAt,
	}
}

func FromMatchSlots(slots []model.MatchSlot) []MatchSlotResponse {
	MatchSlotsResponse := make([]MatchSlotResponse, 0, len(slots))
	for _, slot := range slots {
		MatchSlotsResponse = append(MatchSlotsResponse, FromMatchSlot(&slot))
	}
	return MatchSlotsResponse
}
