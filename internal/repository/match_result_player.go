package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type MatchResultPlayerRepository interface {
	Create(ctx context.Context, player *model.MatchResultPlayer) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.MatchResultPlayer, error)
	ListBySlotID(ctx context.Context, slotID uuid.UUID) ([]model.MatchResultPlayer, error)
	Update(ctx context.Context, player *model.MatchResultPlayer) error
	Delete(ctx context.Context, id uuid.UUID) error
}
