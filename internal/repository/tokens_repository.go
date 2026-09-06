package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type TokenRepository interface {
	Create(ctx context.Context, token *model.MatchAPIToken) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.MatchAPIToken, error)
	GetByOrganisation(ctx context.Context, organisationID uuid.UUID, includeInactive bool) ([]model.MatchAPIToken, error)
	CountActiveTokens(ctx context.Context, organisationID uuid.UUID) (int, error)
	UpdateOrganisation(ctx context.Context, tokenID, organisationID uuid.UUID) error
	Update(ctx context.Context, token *model.MatchAPIToken) error

	Delete(ctx context.Context, id uuid.UUID) error
}

type tokenRepository struct {
	db DBTX
}

func NewTokensRepository(db DBTX) TokenRepository {
	return &tokenRepository{
		db: db,
	}
}

func (r *tokenRepository) Create(ctx context.Context, token *model.MatchAPIToken) error {
	query := `
	INSERT INTO match_api_tokens (
		match_api_token_id,
		added_by,
		organisation_id,
		activation,
		expiration,
		stats_token,
		admin_token,
		player_token
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING
		id,
		created_at,
		updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		token.MatchAPITokenID,
		token.AddedBy,
		token.OrganisationID,
		token.Activation,
		token.Expiration,
		token.StatsToken,
		token.AdminToken,
		token.PlayerToken,
	).Scan(
		&token.ID,
		&token.CreatedAt,
		&token.UpdatedAt,
	)
}

func (r *tokenRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.MatchAPIToken, error) {
	query := `
	SELECT
		id,
		match_api_token_id,
		added_by,
		organisation_id,
		activation,
		expiration,
		stats_token,
		admin_token,
		player_token,
		created_at,
		updated_at
	FROM match_api_tokens
	WHERE id = $1
	`

	var token model.MatchAPIToken

	err := r.db.QueryRow(ctx, query, id).Scan(
		&token.ID,
		&token.MatchAPITokenID,
		&token.AddedBy,
		&token.OrganisationID,
		&token.Activation,
		&token.Expiration,
		&token.StatsToken,
		&token.AdminToken,
		&token.PlayerToken,
		&token.CreatedAt,
		&token.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *tokenRepository) GetByOrganisation(ctx context.Context, organisationID uuid.UUID, includeInactive bool) ([]model.MatchAPIToken, error) {
	query := `
	SELECT
		id,
		match_api_token_id,
		added_by,
		organisation_id,
		activation,
		expiration,
		stats_token,
		admin_token,
		player_token,
		created_at,
		updated_at
	FROM match_api_tokens
	WHERE organisation_id = $1
	`
	if !includeInactive {
		query += ` AND expiration > NOW()`
	}
	query += ` ORDER BY expiration`
	rows, err := r.db.Query(ctx, query, organisationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tokens []model.MatchAPIToken
	for rows.Next() {
		var token model.MatchAPIToken
		err := rows.Scan(
			&token.ID,
			&token.MatchAPITokenID,
			&token.AddedBy,
			&token.OrganisationID,
			&token.Activation,
			&token.Expiration,
			&token.StatsToken,
			&token.AdminToken,
			&token.PlayerToken,
			&token.CreatedAt,
			&token.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, rows.Err()
}

func (r *tokenRepository) CountActiveTokens(ctx context.Context, organisationID uuid.UUID) (int, error) {
	query := `
	SELECT COUNT(*)
	FROM match_api_tokens
	WHERE organisation_id = $1
		AND expiration > NOW()
	`
	var count int
	err := r.db.QueryRow(ctx, query, organisationID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *tokenRepository) UpdateOrganisation(ctx context.Context, tokenID, organisationID uuid.UUID) error {
	query := `
	UPDATE match_api_tokens
	SET 
		updated_at = NOW(),
		organisation_id = $1
	WHERE id = $2
	`
	result, err := r.db.Exec(
		ctx,
		query,
		organisationID,
		tokenID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *tokenRepository) Update(ctx context.Context, token *model.MatchAPIToken) error {
	query := `
	UPDATE match_api_tokens
	SET
		match_api_token_id = $1,
		added_by = $2,
		organisation_id = $3,
		activation = $4,
		expiration = $5,
		stats_token = $6,
		admin_token = $7,
		player_token = $8,
		updated_at = NOW()
	WHERE id = $9
	`

	result, err := r.db.Exec(ctx, query,
		token.MatchAPITokenID,
		token.AddedBy,
		token.OrganisationID,
		token.Activation,
		token.Expiration,
		token.StatsToken,
		token.AdminToken,
		token.PlayerToken,
		token.ID,
	)
	if err != nil {
		return nil
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *tokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE 
	FROM match_api_tokens
	WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
