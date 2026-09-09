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

// Done
func (s *tokenService) Create(ctx context.Context, actorID, organisationID uuid.UUID, token *model.MatchAPIToken) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		tokenRepo := repository.NewTokensRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)

		roleCode, err := orgMemberRepo.GetRole(ctx, organisationID, actorID)
		if err != nil {
			return err
		}
		if (*roleCode != "owner") && (*roleCode != "admin") {
			return ErrForbidden
		}

		token.OrganisationID = &organisationID
		var tokensRequest matchapipb.ListTokensRequest
		tokensRequest.OnlyActive = true
		tokenResponse, err := s.matchAPIClient.ListTokens(ctx, &tokensRequest)
		if err != nil {
			return err
		}
		tokenIDs := tokenResponse.Tokens
		flag := false
		for _, tokenResp := range tokenIDs {
			if (tokenResp.StatsToken == token.StatsToken) && (tokenResp.PlayerToken == *token.PlayerToken) && (tokenResp.AdminToken == *token.AdminToken) {
				flag = true
				token.MatchAPITokenID = int(tokenResp.Id)
			}
		}
		if !flag {
			var createTokenRequest matchapipb.AddTokenRequest
			// sends time to API with ISO 8601 via UTC
			createTokenRequest.Activation = token.Activation.UTC().Format(time.RFC3339)
			createTokenRequest.Expiration = token.Expiration.UTC().Format(time.RFC3339)
			createTokenRequest.AdminToken = *token.AdminToken
			createTokenRequest.PlayerToken = *token.PlayerToken
			createTokenRequest.StatsToken = token.StatsToken
			response, err := s.matchAPIClient.AddToken(ctx, &createTokenRequest)
			if err != nil {
				return err
			}
			token.MatchAPITokenID = int(response.TokenId)
		}
		return tokenRepo.Create(ctx, token)
	})
}

// Done
func (s *tokenService) GetByID(ctx context.Context, tokenID, organisationID, actorID uuid.UUID) (*model.MatchAPIToken, error) {
	roleCode, err := s.organisationMemberRepository.GetRole(ctx, organisationID, actorID)
	if err != nil {
		return nil, err
	}
	if (*roleCode != "owner") && (*roleCode != "admin") {
		return nil, ErrForbidden
	}
	return s.tokenRepository.GetByID(ctx, tokenID, organisationID)
}

// Done
func (s *tokenService) GetByOrganisationID(ctx context.Context, organisationID, actorID uuid.UUID, includeInactive bool) ([]model.MatchAPIToken, error) {
	roleCode, err := s.organisationMemberRepository.GetRole(ctx, organisationID, actorID)
	if err != nil {
		return nil, err
	}
	if (*roleCode != "owner") && (*roleCode != "admin") {
		return nil, ErrForbidden
	}
	return s.tokenRepository.GetByOrganisation(ctx, organisationID, includeInactive)
}

// Done
func (s *tokenService) Update(ctx context.Context, organisationID, actorID, tokenID uuid.UUID, req TokenUpdate) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		tokenRepo := repository.NewTokensRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)

		roleCode, err := orgMemberRepo.GetRole(ctx, organisationID, actorID)
		if err != nil {
			return err
		}
		if (*roleCode != "owner") && (*roleCode != "admin") {
			return ErrForbidden
		}
		token, err := tokenRepo.GetByID(ctx, tokenID, organisationID)
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

// Done
func (s *tokenService) Delete(ctx context.Context, organisationID, actorID, tokenID uuid.UUID) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		tokenRepo := repository.NewTokensRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)
		roleCode, err := orgMemberRepo.GetRole(ctx, organisationID, actorID)
		if err != nil {
			return err
		}
		if (*roleCode != "owner") && (*roleCode != "admin") {
			return ErrForbidden
		}
		return tokenRepo.Delete(ctx, tokenID, organisationID)
	})
}
