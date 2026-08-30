package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type MatchRepository interface {
	Create(ctx context.Context, match *model.Match) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Match, error)
	List(ctx context.Context, limit, offset int) ([]model.Match, error)
	Update(ctx context.Context, match *model.Match) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type matchRepository struct {
	db DBTX
}

func NewMatchRepository(db DBTX) MatchRepository {
	return &matchRepository{
		db: db,
	}
}

func (r *matchRepository) Create(ctx context.Context, match *model.Match) error {
	query := `
	INSERT INTO matches (
		map_id,
		stats_token_id,
		group_id,
		status,
		start_at
	)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING
		id,
		created_at,
		updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		match.MapID,
		match.StatsTokenID,
		match.Status,
		match.StartAt,
	).Scan(
		&match.ID,
		&match.CreatedAt,
		&match.UpdatedAt,
	)
}

func (r *matchRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Match, error) {
	query := `
	SELECT
		id,
		map_id,
		stats_token_id,
		group_id,
		status,
		start_at,
		created_at,
		updated_at
	FROM matches
	WHERE id = $1
	`

	var match model.Match
	err := r.db.QueryRow(ctx, query, id).Scan(
		&match.ID,
		&match.MapID,
		&match.StatsTokenID,
		&match.GroupID,
		&match.Status,
		&match.StartAt,
		&match.CreatedAt,
		&match.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *matchRepository) List(ctx context.Context, limit, offset int) ([]model.Match, error) {
	query := `
	SELECT
		id,
		map_id,
		stats_token_id,
		group_id,
		status,
		start_at,
		created_at,
		updated_at
	FROM matches
	ORDER BY start_at ASC
	LIMIT $1
	OFFSET $2
	`
	// Make a limit 4 that
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var matches []model.Match
	for rows.Next() {
		var match model.Match
		err := rows.Scan(
			&match.ID,
			&match.MapID,
			&match.StatsTokenID,
			&match.GroupID,
			&match.Status,
			&match.StartAt,
			&match.CreatedAt,
			&match.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return matches, nil
}

func (r *matchRepository) Update(ctx context.Context, match *model.Match) error {
	query := `
	UPDATE matches
	SET 
		map_id = $1,
		stats_token_id = $2,
		group_id = $3,
		status = $3,
		start_at = $4,
		updated_at = NOW()
	WHERE id = $5
	`
	cmd, err := r.db.Exec(
		ctx,
		query,
		match.MapID,
		match.StatsTokenID,
		match.GroupID,
		match.Status,
		match.StartAt,
		match.ID,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

func (r *matchRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM matches
	WHERE id = $1
	`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}
