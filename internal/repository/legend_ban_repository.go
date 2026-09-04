package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

// So, this table represents Legends ban prior each match. So, if we're talking
// about ALGS ban system - it won't be any records for the first match
// If we're talking about custom system - it may be
// It is also may be two or more records for each match, as there might be custom ban rules
type MatchLegendsBanRepository interface {
	Create(ctx context.Context, bans *model.MatchLegendBan) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.MatchLegendBan, error)
	GetByMatchID(ctx context.Context, matchID uuid.UUID) ([]model.MatchLegendBan, error)
	GetByLegendID(ctx context.Context, legendID uuid.UUID) ([]model.MatchLegendBan, error)
	Update(ctx context.Context, ban *model.MatchLegendBan) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type matchLegendsBanRepository struct {
	db DBTX
}

func NewMatchLegendsBanRepository(db DBTX) MatchLegendsBanRepository {
	return &matchLegendsBanRepository{
		db: db,
	}
}

func (r *matchLegendsBanRepository) Create(ctx context.Context, bans *model.MatchLegendBan) error {
	query := `
	INSERT INTO match_legend_bans (
		match_id,
		legend_id
	) VALUES ($1, $2)
	RETURNING
		id,
		created_at
	`
	return r.db.QueryRow(ctx, query, bans.MatchID, bans.LegendID).Scan(
		&bans.ID,
		&bans.CreatedAt,
	)
}

func (r *matchLegendsBanRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.MatchLegendBan, error) {
	query := `
	SELECT
		id,
		match_id,
		legend_id,
		created_at
	FROM match_legend_bans
	WHERE id = $1
	`
	var ban model.MatchLegendBan
	err := r.db.QueryRow(ctx, query, id).Scan(
		&ban.ID,
		&ban.MatchID,
		&ban.LegendID,
		&ban.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ban, nil
}

func (r *matchLegendsBanRepository) GetByMatchID(ctx context.Context, matchID uuid.UUID) ([]model.MatchLegendBan, error) {
	query := `
	SELECT
		id,
		match_id,
		legend_id,
		created_at
	FROM match_legend_bans
	WHERE match_id = $1
	`
	rows, err := r.db.Query(ctx, query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bans []model.MatchLegendBan
	for rows.Next() {
		var ban model.MatchLegendBan
		err = rows.Scan(
			&ban.ID,
			&ban.MatchID,
			&ban.LegendID,
			&ban.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		bans = append(bans, ban)
	}
	return bans, nil
}

func (r *matchLegendsBanRepository) GetByLegendID(ctx context.Context, legendID uuid.UUID) ([]model.MatchLegendBan, error) {
	query := `
	SELECT
		id,
		match_id,
		legend_id,
		created_at
	FROM match_legend_bans
	WHERE legend_id = $1
	`
	rows, err := r.db.Query(ctx, query, legendID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bans []model.MatchLegendBan
	for rows.Next() {
		var ban model.MatchLegendBan
		err = rows.Scan(
			&ban.ID,
			&ban.MatchID,
			&ban.LegendID,
			&ban.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		bans = append(bans, ban)
	}
	return bans, nil
}

func (r *matchLegendsBanRepository) Update(ctx context.Context, ban *model.MatchLegendBan) error {
	query := `
	UPDATE match_legend_bans
	SET
		match_id = $1,
		legend_id = $2,
		updated_at = NOW()
	WHERE id = $3
	`

	result, err := r.db.Exec(ctx, query, ban.MatchID, ban.LegendID, ban.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *matchLegendsBanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE 
	FROM match_legend_bans
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
