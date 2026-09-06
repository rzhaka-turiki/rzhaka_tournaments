package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type MapRepository interface {
	Create(ctx context.Context, apexMap *model.Map) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Map, error)
	GetByInGameName(ctx context.Context, inGameName string) (*model.Map, error)
	List(ctx context.Context, limit, offset int) ([]model.Map, error)
	Update(ctx context.Context, apexMap *model.Map) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type mapRepository struct {
	db DBTX
}

func NewMapRepository(db DBTX) MapRepository {
	return &mapRepository{
		db: db,
	}
}

func (r *mapRepository) Create(ctx context.Context, apexMap *model.Map) error {
	query := `
	INSERT INTO maps (
		name,
		in_game_name,
		image_url,
		minimap_image_url,
		supports_drop_spots,
	) VALUES ($1, $2, $3, $4, $5)
	RETURNING
		id,
		created_at,
		updated_at
	`

	return r.db.QueryRow(ctx, query,
		apexMap.Name,
		apexMap.InGameName,
		apexMap.ImageURL,
		apexMap.MinimapImageURL,
		apexMap.SupportsDropSpots,
	).Scan(
		&apexMap.ID,
		&apexMap.CreatedAt,
		&apexMap.UpdatedAt,
	)
}

func (r *mapRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Map, error) {
	query := `
	SELECT
		id,
		name,
		in_game_name,
		image_url,
		minimap_image_url,
		supports_drop_spots,
		created_at,
		updated_at
	FROM maps
	WHERE id = $1
	`
	var apexMap model.Map
	err := r.db.QueryRow(ctx, query, id).Scan(
		&apexMap.ID,
		&apexMap.Name,
		&apexMap.InGameName,
		&apexMap.ImageURL,
		&apexMap.MinimapImageURL,
		&apexMap.SupportsDropSpots,
		&apexMap.CreatedAt,
		&apexMap.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &apexMap, nil
}

func (r *mapRepository) List(ctx context.Context, limit, offset int) ([]model.Map, error) {
	query := `
	SELECT
		id,
		name,
		in_game_name,
		image_url,
		minimap_image_url,
		supports_drop_spots,
		created_at,
		updated_at
	FROM maps
	ORDER BY created_at
	LIMIT $1
	OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var maps []model.Map
	for rows.Next() {
		var apexMap model.Map
		err := rows.Scan(
			&apexMap.ID,
			&apexMap.Name,
			&apexMap.InGameName,
			&apexMap.MinimapImageURL,
			&apexMap.SupportsDropSpots,
			&apexMap.CreatedAt,
			&apexMap.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		maps = append(maps, apexMap)
	}
	return maps, nil
}

func (r *mapRepository) Update(ctx context.Context, apexMap *model.Map) error {
	query := `
	UPDATE maps
	SET
		name = $1,
		in_game_name = $2,
		image_url = $3,
		minimap_image_url = $4,
		supports_drop_spots = $5,
		updated_at = NOW()
	WHERE id = $6
	`

	result, err := r.db.Exec(
		ctx,
		query,
		apexMap.Name,
		apexMap.InGameName,
		apexMap.ImageURL,
		apexMap.MinimapImageURL,
		apexMap.SupportsDropSpots,
		apexMap.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *mapRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE
	FROM maps
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

func (r *mapRepository) GetByInGameName(ctx context.Context, inGameName string) (*model.Map, error) {
	query := `
	SELECT
		id,
		name,
		in_game_name,
		image_url,
		minimap_image_url,
		supports_drop_spots,
		created_at,
		updated_at
	FROM maps
	WHERE in_game_name = $1
	`

	var apexMap model.Map
	err := r.db.QueryRow(ctx, query, inGameName).Scan(
		&apexMap.ID,
		&apexMap.Name,
		&apexMap.InGameName,
		&apexMap.ImageURL,
		&apexMap.MinimapImageURL,
		&apexMap.SupportsDropSpots,
		&apexMap.CreatedAt,
		&apexMap.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &apexMap, nil
}
