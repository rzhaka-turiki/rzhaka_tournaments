package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/database"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/repository"
)

type MatchSlotPlayerService interface {
	Create(ctx context.Context, actorID, organisationID, matchID, slotID uuid.UUID, matchSlotPlayer *model.MatchSlotPlayer) error
	GetByID(ctx context.Context, actorID, organisationID, matchID, slotID, slotPlayerID uuid.UUID) (*model.MatchSlotPlayer, error)
	ListBySlotID(ctx context.Context, actorID, organisationID, matchID, slotID uuid.UUID) ([]model.MatchSlotPlayer, error)
	Update(ctx context.Context, actorID, organisationID, matchID, slotID, slotPlayerID uuid.UUID, req MatchSlotPlayerUpdate) error
	Delete(ctx context.Context, actorID, organisationID, matchID, slotID, slotPlayerID uuid.UUID) error
}

type matchSlotPlayerService struct {
	txManager                 *database.TxManager
	matchSlotPlayerRepository repository.MatchSlotPlayerRepository
}

type MatchSlotPlayerUpdate struct {
	UserID          *uuid.UUID
	ExpectedNIDHash *string
}

func NewMatchSlotPlayerService(txManager *database.TxManager, matchSlotPlayerRepository repository.MatchSlotPlayerRepository) MatchSlotPlayerService {
	return &matchSlotPlayerService{
		txManager:                 txManager,
		matchSlotPlayerRepository: matchSlotPlayerRepository,
	}
}

func (s *matchSlotPlayerService) Create(ctx context.Context, actorID, organisationID, matchID, slotID uuid.UUID, matchSlotPlayer *model.MatchSlotPlayer) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		matchSlotPlayerRepo := repository.NewMatchSlotPlayerRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)

		roleCode, err := orgMemberRepo.GetRole(ctx, organisationID, actorID)
		if err != nil {
			return err
		}
		if (*roleCode != "admin") && (*roleCode != "owner") {
			return ErrUnauthorized
		}
		return matchSlotPlayerRepo.Create(ctx, matchSlotPlayer)
	})
}

func (s *matchSlotPlayerService) GetByID(ctx context.Context, actorID, organisationID, matchID, slotID, slotPlayerID uuid.UUID) (*model.MatchSlotPlayer, error) {
	return s.matchSlotPlayerRepository.GetByID(ctx, slotPlayerID)
}

func (s *matchSlotPlayerService) ListBySlotID(ctx context.Context, actorID, organisationID, matchID, slotID uuid.UUID) ([]model.MatchSlotPlayer, error) {
	return s.matchSlotPlayerRepository.ListBySlotID(ctx, slotID)
}

func (s *matchSlotPlayerService) Update(ctx context.Context, actorID, organisationID, matchID, slotID, slotPlayerID uuid.UUID, req MatchSlotPlayerUpdate) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		matchSlotPlayerRepo := repository.NewMatchSlotPlayerRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)

		roleCode, err := orgMemberRepo.GetRole(ctx, organisationID, actorID)
		if err != nil {
			return err
		}
		if (*roleCode != "owner") && (*roleCode != "admin") {
			return ErrForbidden
		}
		matchSlotPlayer, err := matchSlotPlayerRepo.GetByID(ctx, slotPlayerID)
		if err != nil {
			return err
		}
		if req.ExpectedNIDHash != nil {
			matchSlotPlayer.ExpectedNIDHash = req.ExpectedNIDHash
		}
		if req.UserID != nil {
			matchSlotPlayer.UserID = req.UserID
		}
		return matchSlotPlayerRepo.Update(ctx, matchSlotPlayer)
	})
}

func (s *matchSlotPlayerService) Delete(ctx context.Context, actorID, organisationID, matchID, slotID, slotPlayerID uuid.UUID) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		matchSlotPlayerRepo := repository.NewMatchSlotPlayerRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)

		roleCode, err := orgMemberRepo.GetRole(ctx, organisationID, actorID)
		if err != nil {
			return err
		}
		if (*roleCode != "owner") && (*roleCode != "admin") {
			return ErrForbidden
		}
		return matchSlotPlayerRepo.Delete(ctx, slotPlayerID)
	})
}
