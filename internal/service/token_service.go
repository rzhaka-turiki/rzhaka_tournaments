package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/client/matchapi"
	matchapipb "github.com/rzhaka-turiki/rzhaka_tournaments/internal/client/matchapi/proto"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/database"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/repository"
)

type TokenService interface {
	Create(ctx context.Context, organisationID, actorID uuid.UUID, token *model.MatchAPIToken) error
	GetByID(ctx context.Context, organisationID, actorID, tokenID uuid.UUID) (*model.MatchAPIToken, error)
	GetByOrganisationID(ctx context.Context, organisationID, actorID uuid.UUID, includeInactive bool) ([]model.MatchAPIToken, error)
	Update(ctx context.Context, organisationID, actorID, tokenID uuid.UUID, req TokenUpdate) error
	Delete(ctx context.Context, organisationID, actorID, tokenID uuid.UUID) error
}

type TokenUpdate struct {
	Activation  *time.Time
	Expiration  *time.Time
	AdminToken  *string
	PlayerToken *string
	StatsToken  *string
}

type tokenService struct {
	txManager                    *database.TxManager
	matchAPIClient               matchapi.Client
	tokenRepository              repository.TokenRepository
	organisationRepository       repository.OrganisationsRepository
	organisationMemberRepository repository.OrganisationMembersRepository
}

func NewTokenService(
	txManager *database.TxManager,
	matchAPIClient matchapi.Client,
	tokenRepository repository.TokenRepository,
	organisationRepository repository.OrganisationsRepository,
	organisationMemberRepository repository.OrganisationMembersRepository,
) TokenService {
	return &tokenService{
		txManager:                    txManager,
		matchAPIClient:               matchAPIClient,
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
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)

		// Поменять потом isOwner на проверку разрешений роли
		isOwner, err := orgMemberRepo.HasRole(ctx, organisationID, actorID, "owner")
		if err != nil {
			return err
		}
		isAdmin, err := orgMemberRepo.HasRole(ctx, organisationID, actorID, "admin")
		if err != nil {
			return err
		}
		if !isOwner && !isAdmin {
			return ErrForbidden
		}

		token.OrganisationID = &organisationID
		// needs to add API request to Match API 4 check token_id and check if it exists, if dont - add it
		var tokensRequest matchapipb.ListTokensRequest
		tokensRequest.OnlyActive = true
		tokenResponse, err := s.matchAPIClient.ListTokens(ctx, &tokensRequest)
		if err != nil {
			return err
		}
		tokenIDs := tokenResponse.Tokens
		for _, token := range tokenIDs {
			if token.
		}
		return tokenRepo.Create(ctx, token)
	})
}

func (s *tokenService) GetByID(ctx context.Context, tokenID, organisationID, actorID uuid.UUID) (*model.MatchAPIToken, error) {
	isOwner, err := s.organisationMemberRepository.HasRole(ctx, organisationID, actorID, "owner")
	if err != nil {
		return nil, err
	}
	isAdmin, err := s.organisationMemberRepository.HasRole(ctx, organisationID, actorID, "admin")
	if err != nil {
		return nil, err
	}
	if !isOwner && !isAdmin {
		return nil, ErrForbidden
	}
	return s.tokenRepository.GetByID(ctx, tokenID)
}

func (s *tokenService) GetByOrganisationID(ctx context.Context, organisationID, actorID uuid.UUID, includeInactive bool) ([]model.MatchAPIToken, error) {
	isOwner, err := s.organisationMemberRepository.HasRole(ctx, organisationID, actorID, "owner")
	if err != nil {
		return nil, err
	}
	isAdmin, err := s.organisationMemberRepository.HasRole(ctx, organisationID, actorID, "admin")
	if err != nil {
		return nil, err
	}
	if !isOwner && !isAdmin {
		return nil, ErrForbidden
	}
	return s.tokenRepository.GetByOrganisation(ctx, organisationID, includeInactive)
}

func (s *tokenService) Update(ctx context.Context, organisationID, actorID, tokenID uuid.UUID, req TokenUpdate) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		tokenRepo := repository.NewTokensRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)

		isOwner, err := orgMemberRepo.HasRole(ctx, organisationID, actorID, "owner")
		if err != nil {
			return err
		}
		isAdmin, err := orgMemberRepo.HasRole(ctx, organisationID, actorID, "admin")
		if err != nil {
			return err
		}
		if !isOwner && !isAdmin {
			return ErrForbidden
		}
		token, err := tokenRepo.GetByID(ctx, tokenID)
		if err != nil {
			return err
		}
		if req.Activation != nil {
			token.Activation = *req.Activation
		}
		if req.Expiration != nil {
			token.Expiration = *req.Expiration
		}
		if req.AdminToken != nil {
			token.AdminToken = req.AdminToken
		}
		if req.PlayerToken != nil {
			token.PlayerToken = req.PlayerToken
		}
		if req.StatsToken != nil {
			token.StatsToken = *req.StatsToken
		}
		return tokenRepo.Update(ctx, token)
	})
}

func (s *tokenService) Delete(ctx context.Context, organisationID, actorID, tokenID uuid.UUID) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		tokenRepo := repository.NewTokensRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)
		isOwner, err := orgMemberRepo.HasRole(ctx, organisationID, actorID, "owner")
		if err != nil {
			return err
		}
		isAdmin, err := orgMemberRepo.HasRole(ctx, organisationID, actorID, "admin")
		if err != nil {
			return err
		}
		if !isAdmin && !isOwner {
			return ErrForbidden
		}
		return tokenRepo.Delete(ctx, tokenID)
	})
}
