package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type MatchSlotRepository interface {
	Create(ctx context.Context, slot *model.MatchSlot) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.MatchSlot, error)
	ListByMatchID(ctx context.Context, matchID uuid.UUID) ([]model.MatchSlot, error)
	Update(ctx context.Context, slot *model.MatchSlot) error
	Delete(ctx context.Context, id uuid.UUID) error

	CreatePlayer(ctx context.Context, player *model.MatchSlotPlayer) error
	ListPlayers(ctx context.Context, slotID uuid.UUID) ([]model.MatchSlotPlayer, error)
	UpdatePlayer(ctx context.Context, player *model.MatchSlotPlayer) error
	DeletePlayer(ctx context.Context, playerID uuid.UUID) error
}
