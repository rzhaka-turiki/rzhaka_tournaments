package service

import (
	"context"

	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/dto"
	"github.com/google/uuid"
)

type TeamInviteLinkService interface {
	Create(ctx context.Context, actorID uuid.UUID, req dto.CreateInviteLinkRequest) (*model.TeamInviteLink, error)
}
