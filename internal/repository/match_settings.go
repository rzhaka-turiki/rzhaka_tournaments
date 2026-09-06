package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type MatchSettingsRepository interface {
	Create(ctx context.Context, settings *model.MatchSettings) error
	Update(ctx context.Context, settings *model.MatchSettings) error
	GetByID(ctx context.Context, matchID uuid.UUID) (*model.MatchSettings, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type matchSettingsRepository struct {
	db DBTX
}

func NewMatchSettingsRepository(db DBTX) MatchSettingsRepository {
	return &matchSettingsRepository{
		db: db,
	}
}

func (r *matchSettingsRepository) Create(ctx context.Context, settings *model.MatchSettings) error {
	query := `
	INSERT INTO match_settings (
		match_id,
		drop_spots_enabled,
		playlist_name,
		map_name,
		map_id,
		admin_chat,
		team_rename,
		self_assign,
		aim_assist,
		anon_mode,
		fillBotsMode
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING
		created_at,
		updated_at
	`
	return r.db.QueryRow(ctx, query,
		settings.MatchID,
		settings.DropSpotsEnabled,
		settings.PlaylistName,
		settings.MapName,
		settings.MapID,
		settings.AdminChat,
		settings.TeamRename,
		settings.SelfAssign,
		settings.AimAssist,
		settings.AnonMode,
		settings.FillBotsMode,
	).Scan(&settings.CreatedAt, &settings.UpdatedAt)
}

func (r *matchSettingsRepository) Update(ctx context.Context, settings *model.MatchSettings) error {
	query := `
	UPDATE match_settings
	SET
		drop_spots_enabled = $1,
		playlist_name = $2,
		map_name = $3,
		map_id = $4,
		admin_chat = $5,
		team_rename = $6,
		self_assign = $7,
		aim_assist = $8,
		anon_mode = $9,
		fill_bots_mode = $10,
		updated_at = NOW()
	WHERE match_id = $11
	`
	result, err := r.db.Exec(ctx, query,
		settings.DropSpotsEnabled,
		settings.PlaylistName,
		settings.MapName,
		settings.MapID,
		settings.AdminChat,
		settings.TeamRename,
		settings.SelfAssign,
		settings.AimAssist,
		settings.AnonMode,
		settings.FillBotsMode,
		settings.MatchID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *matchSettingsRepository) GetByID(ctx context.Context, matchID uuid.UUID) (*model.MatchSettings, error) {
	query := `
	SELECT
		match_id,
		drop_spots_enabled,
		playlist_name,
		map_id,
		admin_chat,
		team_rename,
		self_assign,
		aim_assist,
		anon_mode,
		fill_bots_mode,
		created_at,
		updated_at
	FROM match_settings
	WHERE match_id = $1
	`
	var match model.MatchSettings
	err := r.db.QueryRow(ctx, query, matchID).Scan(
		&match.MatchID,
		&match.DropSpotsEnabled,
		&match.PlaylistName,
		&match.MapID,
		&match.AdminChat,
		&match.TeamRename,
		&match.SelfAssign,
		&match.AimAssist,
		&match.AnonMode,
		&match.FillBotsMode,
		&match.CreatedAt,
		&match.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *matchSettingsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE
	FROM match_settings
	WHERE match_id = $1
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
