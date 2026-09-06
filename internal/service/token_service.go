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
	GetByID(ctx context.Context, tokenID uuid.UUID) (*model.MatchAPIToken, error)
	GetByOrganisationID(ctx context.Context, organisationID uuid.UUID) ([]model.MatchAPIToken, error)
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
) tokenService {
	return &tokenService{
		txManager:       txManager,
		tokenRepository: tokenRepository,
	}
}

// What should it do?
// It should create a new db record with token
// Should I create an Org table???
func (s *tokenService) Create(ctx context.Context, actorID uuid.UUID, token *model.MatchAPIToken) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		tokenRepo := repository.NewTokensRepository(tx)

	})
}
