package service

import (
	"context"

	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/google/uuid"
)

type TeamService interface {
	Create(ctx context.Context, team *model.Team) error
	GetByID(ctx context.Context, teamID uuid.UUID) (*model.Team, error)
	GetOwnedTeams(ctx context.Context, userID uuid.UUID, includeArchived bool) ([]model.Team, error)
	GetMemberTeams(ctx context.Context, userID uuid.UUID, includeArchived bool) ([]model.Team, error)
	Update(ctx context.Context, team *model.Team) error
	Archive(ctx context.Context, teamID, userID uuid.UUID) error
	Restore(ctx context.Context, teamID, userID uuid.UUID) error
}

type teamService struct {
}
