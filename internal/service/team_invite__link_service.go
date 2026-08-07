package service

import (
	"context"
	"time"

	"github.com/a1uka/rzhaka_tournaments/internal/database"
	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/a1uka/rzhaka_tournaments/internal/repository"
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TeamInviteLinkService interface {
	CreateByLink(ctx context.Context, actorID uuid.UUID, req dto.CreateInviteLinkRequest) (*model.TeamInviteLink, error)
	GetByToken(ctx context.Context, token string) (*model.TeamInviteLink, error)
	AcceptByLink(ctx context.Context, token string, userID uuid.UUID) error
	DeleteLink(ctx context.Context, actorID, linkID uuid.UUID) error
	GetTeamLinks(ctx context.Context, actorID, teamID uuid.UUID) ([]model.TeamInviteLink, error)
}

type teamInviteLinkService struct {
	txManager *database.TxManager
}

func NewTeamInviteLinkService(
	txManager *database.TxManager,
) TeamInviteLinkService {
	return &teamInviteLinkService{
		txManager: txManager,
	}
}

func (s *teamInviteLinkService) CreateByLink(ctx context.Context, actorID uuid.UUID, req dto.CreateInviteLinkRequest) (*model.TeamInviteLink, error) {
	if req.MaxUses <= 0 {
		return nil, ErrInvalidInviteLink
	}
	if req.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidInviteLink
	}

	var link *model.TeamInviteLink
	err := s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		teamRepo := repository.NewTeamRepository(tx)
		inviteRepo := repository.NewTeamInviteLinkRepository(tx)

		isOwner, err := teamRepo.IsOwner(ctx, req.TeamID, actorID)
		if err != nil {
			return err
		}
		if !isOwner {
			return ErrForbidden
		}
		token := uuid.NewString()
		link = &model.TeamInviteLink{
			TeamID:    req.TeamID,
			Token:     token,
			CreatedBy: actorID,
			MaxUses:   &req.MaxUses,
			ExpiresAt: req.ExpiresAt,
		}
		if err := inviteRepo.Create(ctx, link); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return link, nil
}
