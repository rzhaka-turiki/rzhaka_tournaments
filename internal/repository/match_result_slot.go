package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type MatchResultSlotRepository interface {
	Create(ctx context.Context, slot *model.MatchResultSlot) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.MatchResultSlot, error)
	ListByResultID(ctx context.Context, resultID uuid.UUID) ([]model.MatchResultSlot, error)
	Update(ctx context.Context, slot *model.MatchResultSlot) error
	Delete(ctx context.Context, id uuid.UUID) error
}
