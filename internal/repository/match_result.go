package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type MatchResultRepository interface {
	Create(ctx context.Context, result *model.MatchResult) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.MatchResult, error)
	GetByExternalMID(ctx context.Context, mid string) (*model.MatchResult, error)
	Update(ctx context.Context, result *model.MatchResult) error
	Delete(ctx context.Context, id uuid.UUID) error

	AttachToMatch(ctx context.Context, resultID, matchID uuid.UUID) error
	DeattachFromMatch(ctx context.Context, resultID uuid.UUID) error
}

type matchResultRepository struct {
	db DBTX
}

func NewMatchResultRepository(db DBTX) MatchResultRepository {
	return &matchResultRepository{
		db: db,
	}
}

func (r *matchResultRepository) Create(ctx context.Context, result *model.MatchResult) error {
	query := `
	INSERT INTO match_results (
		match_id,
		external_mid,
		map_id,
		started_at,
		aim_assist_allowed,
	) VALUES ($1, $2, $3, $4, $5)
	RETURNING
		id,
		created_at,
		updated_at
	`
	return r.db.QueryRow(ctx, query,
		result.MatchID,
		result.ExternalMID,
		result.MapID,
		result.StartedAt,
		result.AimAssistAllowed,
	).Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt)
}

func (r *matchResultRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.MatchResult, error) {
	query := `
	SELECT
		id,
		match_id,
		external_mid,
		map_id,
		map_name,
		started_at,
		aim_assist_allowed,
		created_at,
		updated_at
	FROM match_results
	WHERE id = $1
	`
	var result model.MatchResult
	err := r.db.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&result.MatchID,
		&result.ExternalMID,
		&result.MapID,
		&result.MapName,
		&result.StartedAt,
		&result.AimAssistAllowed,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *matchResultRepository) GetByExternalMID(ctx context.Context, mid string) (*model.MatchResult, error) {
	query := `
	SELECT
		id,
		match_id,
		external_mid,
		map_id,
		map_name,
		started_at,
		aim_assist_allowed,
		created_at,
		updated_at
	FROM match_results
	WHERE external_mid = $1
	`
	var result model.MatchResult
	err := r.db.QueryRow(ctx, query, mid).Scan(
		&result.ID,
		&result.MatchID,
		&result.ExternalMID,
		&result.MapID,
		&result.MapName,
		&result.StartedAt,
		&result.AimAssistAllowed,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *matchResultRepository) Update(ctx context.Context, matchResult *model.MatchResult) error {
	query := `
	UPDATE match_results 
	SET
		match_id = $1,
		external_mid = $2,
		map_name = $3,
		map_id = $4,
		started_at = $5,
		aim_assist_allowed = $6,
		updated_at = NOW()
	WHERE id = $7
	`
	result, err := r.db.Exec(ctx, query,
		matchResult.MatchID,
		matchResult.ExternalMID,
		matchResult.MapName,
		matchResult.MapID,
		matchResult.StartedAt,
		matchResult.AimAssistAllowed,
		matchResult.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *matchResultRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE 
	FROM match_results 
	WHERE id = $7
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

func (r *matchResultRepository) AttachToMatch(ctx context.Context, resultID, matchID uuid.UUID) error {
	query := `
	UPDATE match_results
	SET
		match_id = $1,
		updated_at = NOW()
	WHERE id = $2
	`
	result, err := r.db.Exec(ctx, query, matchID, resultID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *matchResultRepository) DeattachFromMatch(ctx context.Context, resultID uuid.UUID) error {
	query := `
	UPDATE match_results
	SET
		match_id = NULL,
		updated_at = NOW()
	WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query, resultID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
