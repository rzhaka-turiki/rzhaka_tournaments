package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type MatchRepository interface {
	Create(ctx context.Context, match *model.Match) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Match, error)
	List(ctx context.Context, limit, offset int) ([]model.Match, error)
	Update(ctx context.Context, match *model.Match) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type matchReposity struct {
	db DBTX
}

func NewMatchRepository(db DBTX) MatchRepository {
	return &matchReposity{
		db: db,
	}
}
