package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/database"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/repository"
)

type MatchSettingsService interface {
	// likely need just a matchSettings model
	Create(ctx context.Context, actorID, organisationID uuid.UUID, matchSettings *model.MatchSettings) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.MatchSettings, error)
}

type matchSettingsService struct {
	txManager                     *database.TxManager
	matchSettingsRepository       repository.MatchSettingsRepository
	organisationMembersRepository repository.OrganisationMembersRepository
}

type MatchSettingsUpdate struct {
}

func NewMatchSettingsService(txManager *database.TxManager, matchSettingsRepository repository.MatchSettingsRepository) MatchSettingsService {
	return &matchSettingsService{
		txManager:               txManager,
		matchSettingsRepository: matchSettingsRepository,
	}
}

func (s *matchSettingsService) Create(ctx context.Context, actorID, organisationID uuid.UUID, matchSettings *model.MatchSettings) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		matchSettingsRepo := repository.NewMatchSettingsRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)
		roleCode, err := orgMemberRepo.GetRole(ctx, organisationID, actorID)
		if err != nil {
			return err
		}
		if (*roleCode != "owner") && (*roleCode != "admin") {
			return ErrForbidden
		}
		return matchSettingsRepo.Create(ctx, matchSettings)
	})
}

func (s *matchSettingsService) GetByID(ctx context.Context, matchSettingsID uuid.UUID) (*model.MatchSettings, error) {
	return s.matchSettingsRepository.GetByID(ctx, matchSettingsID)
}
