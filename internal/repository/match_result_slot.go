package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/shopspring/decimal"
)

type MatchResultSlotRepository interface {
	Create(ctx context.Context, slot *model.MatchResultSlot) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.MatchResultSlot, error)
	ListByResultID(ctx context.Context, resultID uuid.UUID) ([]model.MatchResultSlot, error)
	Update(ctx context.Context, slot *model.MatchResultSlot) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type matchResultSlotRepository struct {
	db DBTX
}

func NewMatchResultSlotRepository(db DBTX) MatchResultSlotRepository {
	return &matchResultSlotRepository{
		db: db,
	}
}

func (r *matchResultSlotRepository) Create(ctx context.Context, slot *model.MatchResultSlot) error {
	query := `
	INSERT INTO match_result_slots (
		match_result_id,
		slot_number,
		team_name,
		team_placement,
		points,
		kills,
		match_slot_id
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING
		id,
		created_at
	`
	err := r.db.QueryRow(ctx, query,
		slot.MatchResultID,
		slot.SlotNumber,
		slot.TeamName,
		slot.TeamPlacement,
		slot.Points.String(),
		slot.Kills,
		slot.MatchSlotID,
	).Scan(&slot.ID, &slot.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *matchResultSlotRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.MatchResultSlot, error) {
	query := `
	SELECT
		id,
		match_result_id,
		slot_number,
		team_name,
		team_placement,
		points,
		kills,
		match_slot_id,
		created_at
	FROM match_result_slots
	WHERE id = $1
	`
	var slot model.MatchResultSlot
	var pointsStr string
	err := r.db.QueryRow(ctx, query, id).Scan(
		&slot.ID,
		&slot.MatchResultID,
		&slot.TeamName,
		&slot.TeamPlacement,
		&pointsStr,
		&slot.Kills,
		&slot.MatchSlotID,
		&slot.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	slot.Points, err = decimal.NewFromString(pointsStr)
	if err != nil {
		return nil, fmt.Errorf("parse points: %w", err)
	}
	return &slot, nil
}

func (r *matchResultSlotRepository) ListByResultID(ctx context.Context, resultID uuid.UUID) ([]model.MatchResultSlot, error) {
	query := `
	SELECT
		id,
		match_result_id,
		slot_number,
		team_name,
		team_placement,
		points,
		kills,
		match_slot_id,
		created_at
	FROM match_result_slots
	WHERE match_result_id = $1
	`
	rows, err := r.db.Query(ctx, query, resultID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var slots []model.MatchResultSlot
	for rows.Next() {
		var slot model.MatchResultSlot
		var slotPts string
		err := rows.Scan(
			&slot.ID,
			&slot.MatchResultID,
			&slot.SlotNumber,
			&slot.TeamName,
			&slot.TeamPlacement,
			&slotPts,
			&slot.Kills,
			&slot.MatchSlotID,
			&slot.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		slot.Points, err = decimal.NewFromString(slotPts)
		if err != nil {
			return nil, fmt.Errorf("parse points: %w", err)
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

func (r *matchResultSlotRepository) Update(ctx context.Context, slot *model.MatchResultSlot) error {
	query := `
	UPDATE match_result_slots
	SET 
		match_result_id = $1,
		slot_number = $2,
		team_name = $3,
		team_placement = $4,
		points = $5,
		kills = $6,
		match_slot_id = $7
	WHERE id = $8
	`
	result, err := r.db.Exec(ctx, query,
		slot.MatchResultID,
		slot.SlotNumber,
		slot.TeamName,
		slot.TeamPlacement,
		slot.Points.String(),
		slot.Kills,
		slot.MatchSlotID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *matchResultSlotRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE
	FROM match_result_slots
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
