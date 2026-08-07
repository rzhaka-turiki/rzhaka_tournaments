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
	CreateLink(ctx context.Context, actorID, teamID uuid.UUID, req dto.CreateInviteLinkRequest) (*model.TeamInviteLink, error)
	GetByToken(ctx context.Context, token string) (*model.TeamInviteLink, error)
	AcceptByLink(ctx context.Context, token string, userID uuid.UUID) error
	DeleteLink(ctx context.Context, actorID, linkID uuid.UUID) error
	PreviewLink(ctx context.Context, token string) (*model.InviteLinkPreview, error)
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

func (s *teamInviteLinkService) CreateLink(ctx context.Context, actorID, teamID uuid.UUID, req dto.CreateInviteLinkRequest) (*model.TeamInviteLink, error) {
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

		isOwner, err := teamRepo.IsOwner(ctx, teamID, actorID)
		if err != nil {
			return err
		}
		if !isOwner {
			return ErrForbidden
		}
		token := uuid.NewString()
		link = &model.TeamInviteLink{
			TeamID:    teamID,
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

func (s *teamInviteLinkService) GetByToken(ctx context.Context, token string) (*model.TeamInviteLink, error) {
	var link *model.TeamInviteLink
	errs := s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		repo := repository.NewTeamInviteLinkRepository(tx)
		var err error
		link, err = repo.GetByToken(ctx, token)
		if err != nil {
			return err
		}
		if link.ExpiresAt.Before(time.Now()) {
			return ErrInvitationExpired
		}
		if link.UsedCount >= *link.MaxUses {
			return ErrInvitationExpired
		}
		return nil
	})
	if errs != nil {
		return nil, errs
	}
	return link, nil
}

func (s *teamInviteLinkService) AcceptByLink(ctx context.Context, token string, userID uuid.UUID) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		inviteLinkRepo := repository.NewTeamInviteLinkRepository(tx)
		requestRepo := repository.NewTeamRequestRepository(tx)
		memberRepo := repository.NewTeamMemberRepository(tx)

		link, err := inviteLinkRepo.GetByToken(ctx, token)
		if err != nil {
			return err
		}
		if time.Now().After(link.ExpiresAt) {
			return ErrInvitationExpired
		}
		if link.UsedCount >= *link.MaxUses {
			return ErrInvitationExpired
		}
		exists, err := memberRepo.Exists(ctx, link.TeamID, userID)
		if err != nil {
			return err
		}
		if exists {
			return ErrAlreadyMember
		}
		exists, err = requestRepo.Exists(ctx, link.TeamID, userID)
		if err != nil {
			return err
		}
		if exists {
			return ErrAlreadyInvited
		}
		request := &model.TeamRequest{
			TeamID:    link.TeamID,
			UserID:    userID,
			CreatedBy: link.CreatedBy,
			Type:      model.TeamRequestInviteLink,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		}
		err = requestRepo.Create(ctx, request)
		if err != nil {
			return err
		}
		err = inviteLinkRepo.IncrementUsage(ctx, link.ID)
		if err != nil {
			return err
		}
		if link.UsedCount+1 >= *link.MaxUses {
			err = inviteLinkRepo.Delete(ctx, link.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *teamInviteLinkService) DeleteLink(ctx context.Context, actorID, linkID uuid.UUID) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		inviteRepo := repository.NewTeamInviteLinkRepository(tx)
		teamRepo := repository.NewTeamRepository(tx)
		link, err := inviteRepo.GetByID(ctx, linkID)
		if err != nil {
			return err
		}
		isOwner, err := teamRepo.IsOwner(ctx, link.TeamID, actorID)
		if err != nil {
			return err
		}
		if !isOwner {
			return ErrForbidden
		}
		return inviteRepo.Delete(ctx, link.ID)
	})
}

func (s *teamInviteLinkService) GetTeamLinks(ctx context.Context, actorID, teamID uuid.UUID) ([]model.TeamInviteLink, error) {
	var links []model.TeamInviteLink
	err := s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		inviteRepo := repository.NewTeamInviteLinkRepository(tx)
		teamRepo := repository.NewTeamRepository(tx)
		isOwner, err := teamRepo.IsOwner(ctx, teamID, actorID)
		if err != nil {
			return err
		}
		if !isOwner {
			return ErrForbidden
		}
		err = inviteRepo.DeleteExpiredByTeam(ctx, teamID)
		if err != nil {
			return err
		}
		links, err = inviteRepo.GetByTeam(ctx, teamID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return links, nil
}

func (s *teamInviteLinkService) PreviewLink(ctx context.Context, token string) (*model.InviteLinkPreview, error) {
	var preview *model.InviteLinkPreview
	err := s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		inviteRepo := repository.NewTeamInviteLinkRepository(tx)
		teamRepo := repository.NewTeamRepository(tx)

		link, err := inviteRepo.GetByToken(ctx, token)
		if err != nil {
			return err
		}

		if time.Now().After(link.ExpiresAt) {
			return ErrInvitationExpired
		}

		if link.MaxUses != nil && link.UsedCount >= *link.MaxUses {
			return ErrInvitationExpired
		}

		team, err := teamRepo.GetByID(ctx, link.TeamID)
		if err != nil {
			return err
		}

		preview = &model.InviteLinkPreview{
			TeamID:    team.ID,
			TeamName:  team.Name,
			ExpiresAt: link.ExpiresAt,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return preview, nil
}
