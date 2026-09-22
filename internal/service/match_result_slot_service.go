package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/database"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/repository"
)

type MatchSlotResultService interface {
	Create(ctx context.Context, actorID, organisationID uuid.UUID, slotResult *model.MatchResultSlot) error
	GetByID(ctx context.Context, slotResultID, organisationID uuid.UUID) (*model.MatchResultSlot, error)
	GetByMatchID(ctx context.Context, matchID, organisationID uuid.UUID) ([]model.MatchResultSlot, error)
	Update(ctx context.Context, actorID, organisationID, slotResultID uuid.UUID, req MatchSlotResultUpdate) error
	Delete(ctx context.Context, actorID, organisationID, slotResultID uuid.UUID) error
}

type MatchSlotResultUpdate struct {
}

type matchSlotResultService struct {
	txManager                 *database.TxManager
	matchSlotResultRepository repository.MatchResultSlotRepository
}

func NewMatchSlotResultService(txManager *database.TxManager, matchSlotResultRepository repository.MatchResultSlotRepository) MatchSlotResultService {
	return &matchSlotResultService{
		txManager:                 txManager,
		matchSlotResultRepository: matchSlotResultRepository,
	}
}

func (s *matchSlotResultService) Create(ctx context.Context, actorID, organisationID uuid.UUID, slotResult *model.MatchResultSlot) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {

	})
}
