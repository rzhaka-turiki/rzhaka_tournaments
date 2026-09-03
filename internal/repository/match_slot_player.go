package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type MatchSlotPlayerRepository interface {
	Create(ctx context.Context, player *model.MatchSlotPlayer) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.MatchSlotPlayer, error)
	ListBySlotID(ctx context.Context, slotID uuid.UUID) ([]model.MatchSlotPlayer, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.MatchSlotPlayer, error)
	Update(ctx context.Context, player *model.MatchSlotPlayer) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type matchSlotPlayerRepository struct {
	db DBTX
}

func NewMatchSlotPlayerRepository(db DBTX) MatchSlotPlayerRepository {
	return &matchSlotPlayerRepository{
		db: db,
	}
}

func (r *matchSlotPlayerRepository) Create(ctx context.Context, player *model.MatchSlotPlayer) error {
	query := `
	INSERT INTO match_slot_players (
		match_slot_id,
		expected_nid_hash,
		user_id,
	) VALUES ($1, $2, $3)
	RETURNING
		id,
		created_at,
		updated_at
	`
	return r.db.QueryRow(ctx, query, player.MatchSlotID, player.ExpectedNIDHash, player.UserID).Scan(
		&player.ID,
		&player.CreatedAt,
		&player.UpdatedAt,
	)
}

func (r *matchSlotPlayerRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.MatchSlotPlayer, error) {
	query := `
	SELECT
		id,
		match_slot_id,
		expected_nid_hash,
		user_id,
		created_at,
		updated_at
	FROM match_slot_players
	WHERE id = $1
	`
	var player model.MatchSlotPlayer
	err := r.db.QueryRow(ctx, query, id).Scan(
		&player.ID,
		&player.MatchSlotID,
		&player.ExpectedNIDHash,
		&player.UserID,
		&player.CreatedAt,
		&player.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *matchSlotPlayerRepository) ListBySlotID(ctx context.Context, slotID uuid.UUID) ([]model.MatchSlotPlayer, error) {
	query := `
	SELECT
		id,
		match_slot_id,
		expected_nid_hash,
		user_id,
		created_at,
		updated_at
	FROM match_slot_players
	WHERE match_slot_id = $1
	`
	rows, err := r.db.Query(ctx, query, slotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var players []model.MatchSlotPlayer
	for rows.Next() {
		var player model.MatchSlotPlayer
		err = rows.Scan(
			&player.ID,
			&player.MatchSlotID,
			&player.ExpectedNIDHash,
			&player.UserID,
			&player.CreatedAt,
			&player.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil
}

func (r *matchSlotPlayerRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.MatchSlotPlayer, error) {
	query := `
	SELECT
		id,
		match_slot_id,
		expected_nid_hash,
		user_id,
		created_at,
		updated_at
	FROM match_slot_players
	WHERE user_id = $1
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var players []model.MatchSlotPlayer
	for rows.Next() {
		var player model.MatchSlotPlayer
		err = rows.Scan(
			&player.ID,
			&player.MatchSlotID,
			&player.ExpectedNIDHash,
			&player.UserID,
			&player.CreatedAt,
			&player.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil
}

func (r *matchSlotPlayerRepository) Update(ctx context.Context, player *model.MatchSlotPlayer) error {
	query := `
	UPDATE match_slot_players
	SET
		match_slot_id = $1,
		expected_nid_hash = $2,
		user_id = $3,
		updated_at = NOW()
	WHERE id = $4
	`
	result, err := r.db.Exec(ctx, query,
		player.MatchSlotID,
		player.ExpectedNIDHash,
		player.UserID,
		player.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *matchSlotPlayerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM match_slot_players
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
