package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type CreateTokenRequest struct {
	StatsToken  string    `json:"stats_token" binding:"required"`
	PlayerToken *string   `json:"player_token"`
	AdminToken  *string   `json:"admin_token"`
	Activation  time.Time `json:"activation" binding:"required"`
	Expiration  time.Time `json:"expiration" binding:"required"`
}

type TokenResponse struct {
	ID              uuid.UUID `json:"id"`
	MatchAPITokenID int       `json:"match_api_token_id"`
	Activation      time.Time `json:"activation"`
	Expiration      time.Time `json:"expiration"`
}

type UpdateTokenRequest struct {
	AdminToken  *string    `json:"admin_token"`
	PlayerToken *string    `json:"player_token"`
	StatsToken  *string    `json:"stats_token"`
	Activation  *time.Time `json:"activation"`
	Expiration  *time.Time `json:"expiration"`
}

func FromToken(token model.MatchAPIToken) TokenResponse {
	return TokenResponse{
		ID:              token.ID,
		MatchAPITokenID: token.MatchAPITokenID,
		Activation:      token.Activation,
		Expiration:      token.Expiration,
	}
}

func FromTokens(tokens []model.MatchAPIToken) []TokenResponse {
	TokensResponse := make([]TokenResponse, 0, len(tokens))
	for _, token := range tokens {
		TokensResponse = append(TokensResponse, FromToken(token))
	}
	return TokensResponse
}
