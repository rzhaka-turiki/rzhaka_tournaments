package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type MatchResultRepository interface {
	Create(ctx context.Context, result *model.MatchResult) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.MatchResult, error)
	GetByExternalMID(ctx context.Context, mid string) (*model.MatchResult, error)
	ListByMatchID(ctx context.Context, matchID uuid.UUID) ([]model.MatchResult, error)
	Update(ctx context.Context, result *model.MatchResult) error

	AttachToMatch(ctx context.Context, resultID, matchID uuid.UUID) error
	DeattachFromMatch(ctx context.Context, resultID uuid.UUID) error
}
