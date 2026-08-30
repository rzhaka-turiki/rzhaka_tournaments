package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type MatchSlotPlayerRepository interface {
	Create(ctx context.Context, player *model.MatchSlotPlayer) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.MatchSlotPlayer, error)
	ListBySlotID(ctx context.Context, slotID uuid.UUID) ([]model.MatchSlotPlayer, error)
	Update(ctx context.Context, player *model.MatchSlotPlayer) error
	Delete(ctx context.Context, id uuid.UUID) error
}
