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

type MatchService interface {
	Create(ctx context.Context, actorID uuid.UUID, match *model.Match) error
	GetByID(ctx context.Context, matchID, organisationID uuid.UUID) (*model.Match, error)
	List(ctx context.Context, organisationID uuid.UUID, limit, offset *int) ([]model.Match, error)
	Update(ctx context.Context, actorID, matchID, organisationID uuid.UUID, req MatchUpdate) error
	Delete(ctx context.Context, actorID, matchID, organisationID uuid.UUID) error
}

type matchService struct {
	txManager       *database.TxManager
	matchRepository repository.MatchRepository
}

type MatchUpdate struct {
	Status    *string
	StartedAt *time.Time
	TokenID   *uuid.UUID
	MapID     *uuid.UUID
}

func NewMatchService(txManager *database.TxManager, matchRepository repository.MatchRepository) MatchService {
	return &matchService{
		txManager:       txManager,
		matchRepository: matchRepository,
	}
}

func (s *matchService) Create(ctx context.Context, actorID uuid.UUID, match *model.Match) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		matchRepo := repository.NewMatchRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)

		roleCode, err := orgMemberRepo.GetRole(ctx, match.OrganisationID, actorID)
		if err != nil {
			return err
		}
		if (*roleCode != "owner") && (*roleCode != "admin") {
			return ErrForbidden
		}
		return matchRepo.Create(ctx, match)
	})
}

func (s *matchService) GetByID(ctx context.Context, matchID, organisationID uuid.UUID) (*model.Match, error) {
	return s.matchRepository.GetByID(ctx, matchID, organisationID)
}

func (s *matchService) List(ctx context.Context, organisationID uuid.UUID, limit, offset *int) ([]model.Match, error) {
	return s.matchRepository.List(ctx, organisationID, limit, offset)
}

func (s *matchService) Update(ctx context.Context, actorID, matchID, organisationID uuid.UUID, req MatchUpdate) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		matchRepo := repository.NewMatchRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)

		roleCode, err := orgMemberRepo.GetRole(ctx, organisationID, actorID)
		if err != nil {
			return err
		}
		if (*roleCode != "owner") && (*roleCode != "admin") {
			return ErrForbidden
		}
		match, err := matchRepo.GetByID(ctx, matchID, organisationID)
		if err != nil {
			return err
		}
		if req.MapID != nil {
			match.MapID = *req.MapID
		}
		if req.StartedAt != nil {
			match.StartAt = req.StartedAt
		}
		if req.Status != nil {
			match.Status = model.MatchStatus(*req.Status)
		}
		if req.TokenID != nil {
			match.StatsTokenID = req.TokenID
		}
		return matchRepo.Update(ctx, match)
	})
}

func (s *matchService) Delete(ctx context.Context, actorID, matchID, organisationID uuid.UUID) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		matchRepo := repository.NewMatchRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)

		roleCode, err := orgMemberRepo.GetRole(ctx, organisationID, actorID)
		if err != nil {
			return err
		}
		if (*roleCode != "owner") && (*roleCode != "admin") {
			return ErrForbidden
		}
		return matchRepo.Delete(ctx, matchID, organisationID)
	})
}
