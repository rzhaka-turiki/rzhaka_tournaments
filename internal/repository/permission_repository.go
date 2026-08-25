package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type PermissionRepository interface {
	GetAll(ctc context.Context) ([]model.Permission, error)
	GetByID(ctx context.Context, ID int) (*model.Permission, error)
	GetByCode(ctx context.Context, code string) (*model.Permission, error)
	Exists(ctx context.Context, id int) (bool, error)
}

type permissionRepository struct {
	db DBTX
}

func NewPermissionRepository(db DBTX) PermissionRepository {
	return &permissionRepository{
		db: db,
	}
}

func (r *permissionRepository) GetAll(ctx context.Context) ([]model.Permission, error) {
	query := `
	SELECT
		id,
		code,
		description
	FROM permissions
	ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var permissions []model.Permission
	for rows.Next() {
		var permission model.Permission
		err := rows.Scan(
			&permission.ID,
			&permission.Code,
			&permission.Description,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, rows.Err()
}

func (r *permissionRepository) GetByID(ctx context.Context, ID int) (*model.Permission, error) {
	query := `
	SELECT
		id,
		code, 
		description
	FROM permissions
	WHERE id = $1
	`

	var permission model.Permission

	err := r.db.QueryRow(ctx, query, ID).Scan(
		&permission.ID,
		&permission.Code,
		&permission.Description,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) GetByCode(ctx context.Context, code string) (*model.Permission, error) {
	query := `
	SELECT
		id,
		code,
		description
	FROM permissions
	WHERE code = $1
	`
	var permission model.Permission
	err := r.db.QueryRow(ctx, query, code).Scan(
		&permission.ID,
		&permission.Code,
		&permission.Description,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) Exists(ctx context.Context, id int) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM permissions
		WHERE id=$1
		AND deleted_at IS NULL
	)
	`

	var exists bool

	err := r.db.QueryRow(ctx, query, id).Scan(&exists)
	return exists, err
}
