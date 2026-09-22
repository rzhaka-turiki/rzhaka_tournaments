package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/database"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/repository"
)

type MatchSlotService interface {
	Create(ctx context.Context, actorID, organisationID uuid.UUID, matchSlot *model.MatchSlot) error
	GetByID(ctx context.Context, slotID, organisationID, matchID uuid.UUID) (*model.MatchSlot, error)
	ListByMatchID(ctx context.Context, matchID, organisationID uuid.UUID) ([]model.MatchSlot, error)
	Update(ctx context.Context, actorID, organisationID, matchID, slotID uuid.UUID, req MatchSlotUpdate) error
	Delete(ctx context.Context, actorID, slotID, matchID, organisationID uuid.UUID) error
}

type matchSlotService struct {
	txManager           *database.TxManager
	matchSlotRepository repository.MatchSlotRepository
}

type MatchSlotUpdate struct {
	SlotNumber *int
	DropSpotID *uuid.UUID
}

func NewMatchSlotService(txManager *database.TxManager, matchSlotRepository repository.MatchSlotRepository) MatchSlotService {
	return &matchSlotService{
		txManager:           txManager,
		matchSlotRepository: matchSlotRepository,
	}
}

// maybe later I should add it to matchCreate tx
func (s *matchSlotService) Create(ctx context.Context, actorID, organisationID uuid.UUID, matchSlot *model.MatchSlot) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		matchSlotRepo := repository.NewMatchSlotRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)

		roleCode, err := orgMemberRepo.GetRole(ctx, organisationID, actorID)
		if err != nil {
			return err
		}
		if (*roleCode != "owner") && (*roleCode != "admin") {
			return ErrForbidden
		}
		return matchSlotRepo.Create(ctx, matchSlot)
	})
}

func (s *matchSlotService) GetByID(ctx context.Context, slotID, organisationID, matchID uuid.UUID) (*model.MatchSlot, error) {
	return s.matchSlotRepository.GetByID(ctx, slotID)
}

func (s *matchSlotService) ListByMatchID(ctx context.Context, matchID, organisationID uuid.UUID) ([]model.MatchSlot, error) {
	return s.matchSlotRepository.ListByMatchID(ctx, matchID)
}

func (s *matchSlotService) Update(ctx context.Context, actorID, organisationID, matchID, slotID uuid.UUID, req MatchSlotUpdate) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		matchSlotRepo := repository.NewMatchSlotRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)

		roleCode, err := orgMemberRepo.GetRole(ctx, organisationID, actorID)
		if err != nil {
			return err
		}
		if (*roleCode != "owner") && (*roleCode != "admin") {
			return ErrForbidden
		}
		matchSlot, err := matchSlotRepo.GetByID(ctx, slotID)
		if err != nil {
			return err
		}
		if req.DropSpotID != nil {
			matchSlot.DropSpotID = req.DropSpotID
		}
		if req.SlotNumber != nil {
			matchSlot.SlotNumber = *req.SlotNumber
		}
		return matchSlotRepo.Update(ctx, matchSlot)
	})
}

func (s *matchSlotService) Delete(ctx context.Context, actorID, matchID, slotID, organisationID uuid.UUID) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		orgMembreRepo := repository.NewOrganisationMembersRepository(tx)
		slotRepo := repository.NewMatchSlotRepository(tx)

		roleCode, err := orgMembreRepo.GetRole(ctx, organisationID, actorID)
		if err != nil {
			return err
		}
		if (*roleCode != "owner") && (*roleCode != "admin") {
			return ErrForbidden
		}
		return slotRepo.Delete(ctx, slotID)
	})
}
