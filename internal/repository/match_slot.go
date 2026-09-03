package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type MatchSlotRepository interface {
	Create(ctx context.Context, slot *model.MatchSlot) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.MatchSlot, error)
	ListByMatchID(ctx context.Context, matchID uuid.UUID) ([]model.MatchSlot, error)
	Update(ctx context.Context, slot *model.MatchSlot) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type matchSlotRepository struct {
	db DBTX
}

func NewMatchSlotRepository(db DBTX) MatchSlotRepository {
	return &matchSlotRepository{
		db: db,
	}
}

func (r *matchSlotRepository) Create(ctx context.Context, slot *model.MatchSlot) error {
	query := `
	INSERT INTO match_slots (
		match_id,
		slot_number,
		drop_spot_id,
	) VALUES ($1, $2, $3)
	RETURNING
		id,
		created_at,
		updated_at
	`

	return r.db.QueryRow(ctx, query, slot.MatchID, slot.SlotNumber, slot.DropSpotID).Scan(
		&slot.ID,
		&slot.CreatedAt,
		&slot.UpdatedAt,
	)
}

func (r *matchSlotRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.MatchSlot, error) {
	query := `
	SELECT
		id,
		match_id,
		slot_number,
		drop_spot_id,
		created_at,
		updated_at
	FROM match_slots
	WHERE id = $1
	`
	var slot model.MatchSlot
	err := r.db.QueryRow(ctx, query, id).Scan(
		&slot.ID,
		&slot.MatchID,
		&slot.SlotNumber,
		&slot.DropSpotID,
		&slot.CreatedAt,
		&slot.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *matchSlotRepository) ListByMatchID(ctx context.Context, matchID uuid.UUID) ([]model.MatchSlot, error) {
	query := `
	SELECT
		id,
		match_id,
		slot_number,
		drop_spot_id,
		created_at,
		updated_at
	FROM match_slots
	WHERE match_id = $1
	`

	rows, err := r.db.Query(ctx, query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var slots []model.MatchSlot
	for rows.Next() {
		var slot model.MatchSlot
		err = rows.Scan(
			&slot.ID,
			&slot.MatchID,
			&slot.SlotNumber,
			&slot.DropSpotID,
			&slot.CreatedAt,
			&slot.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

func (r *matchSlotRepository) Update(ctx context.Context, slot *model.MatchSlot) error {
	query := `
	UPDATE match_slots
	SET
		match_id = $1,
		slot_number = $2,
		drop_spot_id = $3,
		updated_at = NOW()
	WHERE id = $4
	`
	result, err := r.db.Exec(ctx, query,
		slot.MatchID,
		slot.SlotNumber,
		slot.DropSpotID,
		slot.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *matchSlotRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE
	FROM match_slots
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
