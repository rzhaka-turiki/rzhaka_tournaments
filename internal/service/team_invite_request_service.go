package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/database"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/repository"
)

type TeamRequestService interface {
	CreateInvite(ctx context.Context, teamID, actorID, userID uuid.UUID) error
	AcceptRequest(ctx context.Context, requestID, actorID uuid.UUID) error
	RejectRequest(ctx context.Context, requestID, actorID uuid.UUID) error
	GetTeamInvitations(ctx context.Context, teamID, actorID uuid.UUID) ([]model.TeamRequest, error)
	GetUserInvitations(ctx context.Context, userID uuid.UUID) ([]model.TeamRequest, error)
}

type teamRequestService struct {
	txManager *database.TxManager

	teamRepo    repository.TeamRepository
	memberRepo  repository.TeamMemberRepository
	requestRepo repository.TeamRequestRepository
}

func NewTeamRequestService(
	txManager *database.TxManager,
	teamRepo repository.TeamRepository,
	memberRepo repository.TeamMemberRepository,
	requestRepo repository.TeamRequestRepository,
) TeamRequestService {
	return &teamRequestService{
		txManager:   txManager,
		teamRepo:    teamRepo,
		memberRepo:  memberRepo,
		requestRepo: requestRepo,
	}
}

const (
	TeamInviteTTL = 7 * 24 * time.Hour
)

func (s *teamRequestService) CreateInvite(ctx context.Context, teamID, actorID, userID uuid.UUID) error {
	isOwner, err := s.teamRepo.IsOwner(ctx, teamID, actorID)
	if err != nil {
		return err
	}
	if !isOwner {
		return ErrForbidden

	}
	exists, err := s.memberRepo.Exists(ctx, teamID, userID)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyMember
	}
	exists, err = s.requestRepo.Exists(ctx, teamID, userID)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyExists
	}
	req := &model.TeamRequest{
		TeamID:    teamID,
		UserID:    userID,
		CreatedBy: actorID,
		Type:      model.TeamRequestInvite,
		ExpiresAt: time.Now().Add(TeamInviteTTL),
	}

	return s.create(ctx, req)
}

func (s *teamRequestService) create(ctx context.Context, req *model.TeamRequest) error {
	exists, err := s.requestRepo.Exists(ctx, req.TeamID, req.UserID)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyExists
	}

	return s.requestRepo.Create(ctx, req)
}

func (s *teamRequestService) AcceptRequest(ctx context.Context, requestID, actorID uuid.UUID) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		requestRepo := repository.NewTeamRequestRepository(tx)
		memberRepo := repository.NewTeamMemberRepository(tx)
		teamRepo := repository.NewTeamRepository(tx)

		req, err := requestRepo.GetByID(ctx, requestID)
		if err != nil {
			return err
		}
		if time.Now().After(req.ExpiresAt) {
			_ = requestRepo.Delete(ctx, req.ID)
			return ErrInvitationExpired
		}
		switch req.Type {
		case model.TeamRequestInvite:
			if req.UserID != actorID {
				return ErrForbidden
			}
		case model.TeamRequestInviteLink:
			isOwner, err := teamRepo.IsOwner(ctx, req.TeamID, actorID)
			if err != nil {
				return err
			}
			if !isOwner {
				return ErrForbidden
			}
		default:
			return repository.ErrInvalid
		}
		exists, err := memberRepo.Exists(ctx, req.TeamID, req.UserID)
		if err != nil {
			return err
		}
		if exists {
			_ = requestRepo.Delete(ctx, req.ID)
			return ErrAlreadyMember
		}
		err = memberRepo.Create(ctx, &model.TeamMember{
			TeamID: req.TeamID,
			UserID: req.UserID,
			Role:   model.TeamRolePlayer,
		})
		if err != nil {
			return err
		}
		return requestRepo.Delete(ctx, req.ID)
	})
}

func (s *teamRequestService) RejectRequest(ctx context.Context, requestID, actorID uuid.UUID) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		requestRepo := repository.NewTeamRequestRepository(tx)
		teamRepo := repository.NewTeamRepository(tx)
		req, err := requestRepo.GetByID(ctx, requestID)
		if err != nil {
			return err
		}
		if time.Now().After(req.ExpiresAt) {
			_ = requestRepo.Delete(ctx, req.ID)
			return ErrInvitationExpired
		}
		switch req.Type {
		case model.TeamRequestInvite:
			if req.UserID != actorID {
				return ErrForbidden
			}
		case model.TeamRequestInviteLink:
			isOwner, err := teamRepo.IsOwner(ctx, req.TeamID, actorID)
			if err != nil {
				return err
			}
			if !isOwner {
				return ErrForbidden
			}
		default:
			return repository.ErrInvalid
		}
		return requestRepo.Delete(ctx, req.ID)
	})
}

func (s *teamRequestService) GetUserInvitations(ctx context.Context, userID uuid.UUID) ([]model.TeamRequest, error) {
	var requests []model.TeamRequest

	errs := s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		requestRepo := repository.NewTeamRequestRepository(tx)
		var err error
		if err := requestRepo.DeleteExpiredByUser(ctx, userID); err != nil {
			return err
		}

		requests, err = requestRepo.GetByUser(ctx, userID)
		if err != nil {
			return err
		}

		return nil
	})

	if errs != nil {
		return nil, errs
	}

	return requests, nil
}

func (s *teamRequestService) GetTeamInvitations(ctx context.Context, teamID, actorID uuid.UUID) ([]model.TeamRequest, error) {
	var requests []model.TeamRequest

	errs := s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		requestRepo := repository.NewTeamRequestRepository(tx)
		var err error
		if err := requestRepo.DeleteExpiredByTeam(ctx, teamID); err != nil {
			return err
		}

		requests, err = requestRepo.GetByTeam(ctx, teamID)
		if err != nil {
			return err
		}

		return nil
	})

	if errs != nil {
		return nil, errs
	}

	return requests, nil
}
