package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type TeamSnapshotRepository interface {
	Create(ctx context.Context, snapshot *model.TeamSnapshot) error
	GetByID(ctx context.Context, ID uuid.UUID) (*model.TeamSnapshot, error)
	Update(ctx context.Context, snapshot *model.TeamSnapshot) error
}

type teamSnapshotRepository struct {
	db DBTX
}

func NewTeamSnapshotRepository(db DBTX) TeamSnapshotRepository {
	return &teamSnapshotRepository{
		db: db,
	}
}

func (r *teamSnapshotRepository) Create(ctx context.Context, snapshot *model.TeamSnapshot) error {
	query := `
	INSERT INTO team_snapshots (
		team_id,
		name,
		short_name,
		logo_path,
		logo_dark_path
	)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING
		id,
		created_at
	`
	return r.db.QueryRow(
		ctx,
		query,
		snapshot.TeamID,
		snapshot.Name,
		snapshot.ShortName,
		snapshot.LogoPath,
		snapshot.LogoDarkPath,
	).Scan(&snapshot.ID, &snapshot.CreatedAt)
}

func (r *teamSnapshotRepository) GetByID(ctx context.Context, ID uuid.UUID) (*model.TeamSnapshot, error) {
	query := `
	SELECT
		id,
		team_id,
		name,
		short_name,
		logo_path,
		logo_dark_path,
		created_at
	FROM team_snapshots
	WHERE id = $1
	`
	var snapshot model.TeamSnapshot
	err := r.db.QueryRow(ctx, query, ID).Scan(
		&snapshot.ID,
		&snapshot.TeamID,
		&snapshot.Name,
		&snapshot.ShortName,
		&snapshot.LogoPath,
		&snapshot.LogoDarkPath,
		&snapshot.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &snapshot, nil
}
func (r *teamSnapshotRepository) Update(ctx context.Context, snapshot *model.TeamSnapshot) error {
	query := `
	UPDATE team_snapshots
	SET
		name = $1,
		short_name = $2,
		logo_path = $3,
		logo_dark_path = $4
	WHERE id = $5
	`
	result, err := r.db.Exec(
		ctx,
		query,
		snapshot.Name,
		snapshot.ShortName,
		snapshot.LogoPath,
		snapshot.LogoDarkPath,
		snapshot.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
