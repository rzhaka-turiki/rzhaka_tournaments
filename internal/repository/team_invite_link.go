package repository

import (
	"context"

	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/google/uuid"
)

type TeamInviteLinkRepository interface {
	Create(ctx context.Context, req *model.TeamInviteLink) error
	GetByToken(ctx context.Context, token string) (*model.TeamInviteLink, error)
	IncrementUsage(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}
