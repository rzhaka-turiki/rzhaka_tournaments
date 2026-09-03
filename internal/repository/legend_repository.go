package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type LegendsRepository interface {
	Create(ctx context.Context, legend *model.Legend) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Legend, error)
	GetByInGameName(ctx context.Context, inGameName string) (*model.Legend, error)
	Update(ctx context.Context, legend *model.Legend) error
	Delete(ctx context.Context, id uuid.UUID) error
}
