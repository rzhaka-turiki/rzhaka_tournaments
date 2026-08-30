package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type MatchSettingsRepository interface {
	Create(ctx context.Context, settings *model.MatchSettings) error
	Update(ctx context.Context, settings *model.MatchSettings) error
	GetByID(ctx context.Context, matchgID uuid.UUID) (*model.MatchSettings, error)
}
