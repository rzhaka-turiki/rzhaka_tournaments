package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type LegendsRepository interface {
	Create(ctx context.Context, legend *model.Legend) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Legend, error)
	GetAll(ctx context.Context) ([]model.Legend, error)
	GetByInGameName(ctx context.Context, inGameName string) (*model.Legend, error)
	Update(ctx context.Context, legend *model.Legend) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type legendsRepository struct {
	db DBTX
}

func NewLegendsRepository(db DBTX) LegendsRepository {
	return &legendsRepository{
		db: db,
	}
}

func (r *legendsRepository) Create(ctx context.Context, legend *model.Legend) error {
	query := `
	INSERT INTO legends (
		name,
		in_game_name,
		image_url,
		profile_image_url,
		class,
		ability,
		ultimate
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING
		id,
		created_at,
		updated_at
	`

	return r.db.QueryRow(ctx, query,
		legend.Name,
		legend.InGameName,
		legend.ImageURL,
		legend.ProfileImageURL,
		legend.Class,
		legend.Ability,
		legend.Ultimate,
	).Scan(
		&legend.ID,
		&legend.CreatedAt,
		&legend.UpdatedAt,
	)
}

func (r *legendsRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Legend, error) {
	query := `
	SELECT
		id,
		name,
		in_game_name,
		image_url,
		profile_image_url,
		class,
		ability,
		ultimate,
		created_at,
		updated_at
	FROM legends
	WHERE id = $1
	`
	var legend model.Legend
	err := r.db.QueryRow(ctx, query, id).Scan(
		&legend.ID,
		&legend.Name,
		&legend.InGameName,
		&legend.ImageURL,
		&legend.ProfileImageURL,
		&legend.Class,
		&legend.Ability,
		&legend.CreatedAt,
		&legend.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &legend, nil
}

func (r *legendsRepository) GetAll(ctx context.Context) ([]model.Legend, error) {
	query := `
	SELECT
		id,
		name,
		in_game_name,
		image_url,
		profile_image_url,
		class,
		ability,
		ultimate,
		created_at,
		updated_at
	FROM legends
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var legends []model.Legend
	for rows.Next() {
		var legend model.Legend
		err = rows.Scan(
			&legend.ID,
			&legend.Name,
			&legend.InGameName,
			&legend.ImageURL,
			&legend.ProfileImageURL,
			&legend.Class,
			&legend.Ability,
			&legend.Ultimate,
			&legend.CreatedAt,
			&legend.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		legends = append(legends, legend)
	}
	return legends, nil
}

func (r *legendsRepository) GetByInGameName(ctx context.Context, inGameName string) (*model.Legend, error) {
	query := `
	SELECT
		id,
		name,
		in_game_name,
		image_url,
		profile_image_url,
		class,
		ability,
		ultimate,
		created_at,
		updated_at
	FROM legends
	WHERE in_game_name = $1
	`
	var legend model.Legend
	err := r.db.QueryRow(ctx, query, inGameName).Scan(
		&legend.ID,
		&legend.Name,
		&legend.InGameName,
		&legend.ImageURL,
		&legend.ProfileImageURL,
		&legend.Class,
		&legend.Ability,
		&legend.CreatedAt,
		&legend.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &legend, nil
}

func (r *legendsRepository) Update(ctx context.Context, legend *model.Legend) error {
	query := `
	UPDATE legends
	SET
		name = $1,
		in_game_name = $2,
		image_url = $3,
		profile_image_url = $4,
		class = $5,
		ability = $6,
		ultimate = $7,
		updated_at = NOW()
	WHERE id = $8
	`
	result, err := r.db.Exec(ctx, query,
		legend.Name,
		legend.InGameName,
		legend.ImageURL,
		legend.ProfileImageURL,
		legend.Class,
		legend.Ability,
		legend.Ultimate,
		legend.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *legendsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE
	FROM legends
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
