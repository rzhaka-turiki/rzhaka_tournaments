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

type MatchResultService interface {
	Create(ctx context.Context, actorID, organisationID, matchID uuid.UUID, matchResult *model.MatchResult) error
	GetByMatchID(ctx context.Context, matchID, organisationID uuid.UUID) (*model.MatchResult, error)
	GetByExternalMID(ctx context.Context, mid string, organisationID uuid.UUID) (*model.MatchResult, error)
	Update(ctx context.Context, actorID, organisationID, resultID uuid.UUID, req MatchResultUpdate) error
	Delete(ctx context.Context, actorID, organisationID, resultID uuid.UUID) error
	/*
		AttachToMatch(ctx context.Context, actorID, organisationID, resultID, matchID uuid.UUID) error
		DeattachFromMatch(ctx context.Context, actorID, organisationID, resultID uuid.UUID) error
	*/
}

type MatchResultUpdate struct {
	ExternalMID      *string
	MapName          *string
	MapID            *uuid.UUID
	StartedAt        *time.Time
	AimAssistAllowed *bool
}

type matchResultService struct {
	txManager             *database.TxManager
	matchResultRepository repository.MatchResultRepository
	matchRepository       repository.MatchRepository
}

func NewMatchResultService(txManager *database.TxManager, matchResultRepository repository.MatchResultRepository, matchRepository repository.MatchRepository) MatchResultService {
	return &matchResultService{
		txManager:             txManager,
		matchResultRepository: matchResultRepository,
		matchRepository:       matchRepository,
	}
}

func (s *matchResultService) Create(ctx context.Context, actorID, organisationID, matchID uuid.UUID, matchResult *model.MatchResult) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		matchResultRepo := repository.NewMatchResultRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)
		matchRepo := repository.NewMatchRepository(tx)

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
		if match == nil {
			return ErrForbidden
		}

		return matchResultRepo.Create(ctx, matchResult)
	})
}

func (s *matchResultService) GetByExternalMID(ctx context.Context, mid string, organisationID uuid.UUID) (*model.MatchResult, error) {
	matchResult, err := s.matchResultRepository.GetByExternalMID(ctx, mid)
	if err != nil {
		return nil, err
	}
	match, err := s.matchRepository.GetByID(ctx, matchResult.MatchID, organisationID)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, ErrForbidden
	}
	return matchResult, err
}

func (s *matchResultService) GetByMatchID(ctx context.Context, matchID, organisationID uuid.UUID) (*model.MatchResult, error) {
	matchResult, err := s.matchResultRepository.GetByMatchID(ctx, matchID)
	if err != nil {
		return nil, err
	}
	match, err := s.matchRepository.GetByID(ctx, matchResult.MatchID, organisationID)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, ErrForbidden
	}
	return matchResult, err
}

func (s *matchResultService) Update(ctx context.Context, actorID, organisationID, resultID uuid.UUID, req MatchResultUpdate) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)
		matchResultRepo := repository.NewMatchResultRepository(tx)
		matchRepo := repository.NewMatchRepository(tx)

		roleCode, err := orgMemberRepo.GetRole(ctx, organisationID, actorID)
		if err != nil {
			return err
		}
		if (*roleCode != "admin") && (*roleCode != "owner") {
			return ErrForbidden
		}
		match, err := matchRepo.GetByID(ctx, resultID, organisationID)
		if err != nil {
			return err
		}
		if match == nil {
			return ErrForbidden
		}
		matchResult, err := matchResultRepo.GetByMatchID(ctx, match.ID)
		if err != nil {
			return err
		}
		if req.AimAssistAllowed != nil {
			matchResult.AimAssistAllowed = *req.AimAssistAllowed
		}
		if req.MapName != nil {
			matchResult.MapName = *req.MapName
		}
		if req.ExternalMID != nil {
			matchResult.ExternalMID = *req.ExternalMID
		}
		if req.StartedAt != nil {
			matchResult.StartedAt = *req.StartedAt
		}
		if req.MapID != nil {
			matchResult.MapID = req.MapID
		}
		return matchResultRepo.Update(ctx, matchResult)
	})
}

func (s *matchResultService) Delete(ctx context.Context, actorID, organisationID, matchID uuid.UUID) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		matchResultRepo := repository.NewMatchResultRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)
		matchRepo := repository.NewMatchRepository(tx)

		roleCode, err := orgMemberRepo.GetRole(ctx, organisationID, actorID)
		if err != nil {
			return err
		}
		if (*roleCode != "admin") && (*roleCode != "owner") {
			return ErrForbidden
		}
		match, err := matchRepo.GetByID(ctx, matchID, organisationID)
		if err != nil {
			return err
		}
		if match == nil {
			return ErrForbidden
		}
		return matchResultRepo.Delete(ctx, matchID)
	})
}
