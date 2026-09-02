package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type MatchResultPlayerRepository interface {
	Create(ctx context.Context, player *model.MatchResultPlayer) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.MatchResultPlayer, error)
	ListBySlotID(ctx context.Context, slotID uuid.UUID) ([]model.MatchResultPlayer, error)
	Update(ctx context.Context, player *model.MatchResultPlayer) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type matchResultPlayerRepository struct {
	db DBTX
}

func NewMatchResultPlayerRepository(db DBTX) MatchResultPlayerRepository {
	return &matchResultPlayerRepository{
		db: db,
	}
}

func (r *matchResultPlayerRepository) Create(ctx context.Context, player *model.MatchResultPlayer) error {
	query := `
	INSERT INTO match_result_players (
		match_result_slot_id,
		nid_hash,
		player_name,
		legend_id,
		character_name,
		kills,
		assists,
		knockdowns,
		damage_dealt,
		survival_time,
		hardware,
		headshots,
		shots,
		hits,
		respawns,
		revives
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
		$11, $12, $13, $14, $15, $16)
	RETURNING
		id,
		created_at
	`

	return r.db.QueryRow(ctx, query,
		player.MatchResultSlotID,
		player.NIDHash,
		player.PlayerName,
		player.LegendID,
		player.CharacterName,
		player.Kills,
		player.Assists,
		player.Knockdowns,
		player.DamageDealt,
		player.SurvivalTime,
		player.Hardware,
		player.Headshots,
		player.Shots,
		player.Hits,
		player.Respawns,
		player.Revives,
	).Scan(
		&player.ID,
		&player.CreatedAt,
	)
}

func (r *matchResultPlayerRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.MatchResultPlayer, error) {
	query := `
	SELECT
		id,
		match_result_slot_id,
		nid_hash,
		player_name,
		legend_id,
		character_name,
		kills,
		assists,
		knockdowns,
		damage_dealt,
		survival_time,
		hardware,
		headshots,
		shots,
		hits,
		respawns,
		revives,
		created_at
	FROM match_result_players
	WHERE id = $1
	`

	var player model.MatchResultPlayer
	err := r.db.QueryRow(ctx, query, id).Scan(
		&player.ID,
		&player.MatchResultSlotID,
		&player.NIDHash,
		&player.PlayerName,
		&player.LegendID,
		&player.CharacterName,
		&player.Kills,
		&player.Assists,
		&player.Knockdowns,
		&player.DamageDealt,
		&player.SurvivalTime,
		&player.Hardware,
		&player.Headshots,
		&player.Shots,
		&player.Hits,
		&player.Respawns,
		&player.Revives,
		&player.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *matchResultPlayerRepository) ListBySlotID(ctx context.Context, slotID uuid.UUID) ([]model.MatchResultPlayer, error) {
	query := `
	SELECT
		id,
		match_result_slot_id,
		nid_hash,
		player_name,
		legend_id,
		character_name,
		kills,
		assists,
		knockdowns,
		damage_dealt,
		survival_time,
		hardware,
		headshots,
		shots,
		hits,
		respawns,
		revives,
		created_at
	FROM match_result_players
	WHERE match_result_slot_id = $1
	`

	rows, err := r.db.Query(ctx, query, slotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var players []model.MatchResultPlayer
	for rows.Next() {
		var player model.MatchResultPlayer
		err := rows.Scan(
			&player.ID,
			&player.MatchResultSlotID,
			&player.NIDHash,
			&player.PlayerName,
			&player.LegendID,
			&player.CharacterName,
			&player.Kills,
			&player.Assists,
			&player.Knockdowns,
			&player.DamageDealt,
			&player.SurvivalTime,
			&player.Hardware,
			&player.Headshots,
			&player.Shots,
			&player.Hits,
			&player.Respawns,
			&player.Revives,
			&player.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, rows.Err()
}

func (r *matchResultPlayerRepository) Update(ctx context.Context, player *model.MatchResultPlayer) error {
	query := `
	UPDATE match_result_players
	SET
		match_result_slot_id = $1,
		nid_hash = $2,
		player_name = $3,
		legend_id = $4,
		character_name = $5,
		kills = $6,
		assists = $7,
		knockdowns = $8,
		damage_dealt = $9,
		survival_time = $10,
		hardware = $11,
		headshots = $12,
		shots = $13,
		hits = $14,
		respawns = $15,
		revives = $16,
	WHERE id = $17
	`
	result, err := r.db.Exec(ctx, query,
		player.MatchResultSlotID,
		player.NIDHash,
		player.PlayerName,
		player.LegendID,
		player.CharacterName,
		player.Kills,
		player.Assists,
		player.Knockdowns,
		player.DamageDealt,
		player.SurvivalTime,
		player.Hardware,
		player.Headshots,
		player.Shots,
		player.Hits,
		player.Respawns,
		player.Revives,
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

func (r *matchResultPlayerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE
	FROM match_result_players
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
