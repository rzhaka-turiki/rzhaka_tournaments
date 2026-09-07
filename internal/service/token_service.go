package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/database"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/repository"
)

type TokenService interface {
	Create(ctx context.Context, actorID uuid.UUID, token *model.MatchAPIToken) error
	GetByID(ctx context.Context, actorID, tokenID uuid.UUID) (*model.MatchAPIToken, error)
	GetByOrganisationID(ctx context.Context, organisationID uuid.UUID, includeInactive bool) ([]model.MatchAPIToken, error)
	Update(ctx context.Context, actorID uuid.UUID, token *model.MatchAPIToken) error
	Delete(ctx context.Context, actorID uuid.UUID, token *model.MatchAPIToken) error
}

type tokenService struct {
	txManager       *database.TxManager
	tokenRepository repository.TokenRepository
}

func NewTokenService(
	txManager *database.TxManager,
	tokenRepository repository.TokenRepository,
	organisationRepository repository.OrganisationsRepository,
	organisationMemberRepository repository.OrganisationMembersRepository,
) tokenService {
	return &tokenService{
		txManager:                    txManager,
		tokenRepository:              tokenRepository,
		organisationRepository:       organisationRepository,
		organisationMemberRepository: organisationMemberRepository,
	}
}

// What should it do?
// It should create a new db record with token
// Should I create an Org table???
func (s *tokenService) Create(ctx context.Context, actorID, organisationID uuid.UUID, token *model.MatchAPIToken) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		tokenRepo := repository.NewTokensRepository(tx)
		orgRepo := repository.NewOrganisationRepository(tx)

		// Поменять потом isOwner на проверку разрешений роли
		isOwner, err := orgRepo.IsOwner(ctx, organisationID, actorID)
		if err != nil {
			return err
		}
		if !isOwner {
			return ErrForbidden
		}
		token.OrganisationID = &organisationID
		return tokenRepo.Create(ctx, token)
	})
}

func (s *tokenService) GetByID(ctx context.Context, tokenID uuid.UUID) (*model.MatchAPIToken, error) {
	return s.tokenRepository.GetByID(ctx, tokenID)
}

func (s *tokenService) GetByOrganisationID(ctx context.Context, organisationID uuid.UUID, includeInactive bool) ([]model.MatchAPIToken, error) {
	return s.tokenRepository.GetByOrganisation(ctx, organisationID, includeInactive)
}

func (s *tokenService) Update(ctx context.Context, actorID uuid.UUID, token *model.MatchAPIToken) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		tokenRepo := repository.NewTokensRepository(tx)
		orgRepo := repository.NewOrganisationRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)

	})
}
