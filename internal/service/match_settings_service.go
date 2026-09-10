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
	GetByMatchID(ctx context.Context, matchID uuid.UUID) (*model.MatchSettings, error)
	List(ctx context.Context, matchIDs []uuid.UUID) ([]model.MatchSettings, error)
	Update(ctx context.Context, organisationID, actorID, matchID uuid.UUID, req MatchSettingsUpdate) error
	Delete(ctx context.Context, organisationID, actorID, matchID uuid.UUID) error
}

type matchSettingsService struct {
	txManager               *database.TxManager
	matchSettingsRepository repository.MatchSettingsRepository
}

type MatchSettingsUpdate struct {
	MapID            *uuid.UUID
	PlaylistName     *string
	MapName          *string
	AdminChat        *bool
	TeamRename       *bool
	SelfAssign       *bool
	AimAssist        *bool
	AnonMode         *bool
	DropSpotsEnabled *bool
	FillBotsMode     *bool
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

// Да, это репчик
func (s *matchSettingsService) GetByMatchID(ctx context.Context, matchID uuid.UUID) (*model.MatchSettings, error) {
	return s.matchSettingsRepository.GetByID(ctx, matchID)
}

func (s *matchSettingsService) List(ctx context.Context, matchIDs []uuid.UUID) ([]model.MatchSettings, error) {
	return s.matchSettingsRepository.List(ctx, matchIDs)
}

func (s *matchSettingsService) Update(ctx context.Context, organisationID, actorID, matchID uuid.UUID, req MatchSettingsUpdate) error {
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
		matchSetting, err := matchSettingsRepo.GetByID(ctx, matchID)
		if err != nil {
			return err
		}
		if req.MapID != nil {
			matchSetting.MapID = *req.MapID
		}
		if req.PlaylistName != nil {
			matchSetting.MapName = *req.MapName
		}
		if req.AdminChat != nil {
			matchSetting.AdminChat = *req.AdminChat
		}
		if req.MapName != nil {
			matchSetting.MapName = *req.MapName
		}
		if req.TeamRename != nil {
			matchSetting.TeamRename = *req.TeamRename
		}
		if req.SelfAssign != nil {
			matchSetting.SelfAssign = *req.SelfAssign
		}
		if req.AimAssist != nil {
			matchSetting.AimAssist = *req.AimAssist
		}
		if req.AnonMode != nil {
			matchSetting.AnonMode = *req.AnonMode
		}
		if req.DropSpotsEnabled != nil {
			matchSetting.DropSpotsEnabled = *req.DropSpotsEnabled
		}
		if req.FillBotsMode != nil {
			matchSetting.FillBotsMode = *req.FillBotsMode
		}
		return matchSettingsRepo.Update(ctx, matchSetting)
	})
}

func (s *matchSettingsService) Delete(ctx context.Context, organisationID, actorID, matchID uuid.UUID) error {
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
		return matchSettingsRepo.Delete(ctx, matchID)
	})
}
