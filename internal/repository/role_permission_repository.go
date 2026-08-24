package repository

import (
	"context"

	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type RolePermissionRepository interface {
	GetRolePermissions(ctx context.Context, roleID int) ([]model.Permission, error)
	AddPermission(ctx context.Context, roleID int, permissionID int) error
	RemovePermission(ctx context.Context, roleID int, permissionID int) error
}

type rolePermissionRepository struct {
	db DBTX
}

func NewRolePermissionRepository(db DBTX) RolePermissionRepository {
	return &rolePermissionRepository{
		db: db,
	}
}

func (r *rolePermissionRepository) GetRolePermissions(ctx context.Context, roleID int) ([]model.Permission, error) {
	query := `
	SELECT
		p.id,
		p.code,
		p.description
	FROM permissions p
	JOIN role_permissions rp
	ON rp.permission_id = p.id
	WHERE rp.role_id = $1
	ORDER BY p.id
	`
	rows, err := r.db.Query(ctx, query, roleID)
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

func (r *rolePermissionRepository) AddPermission(ctx context.Context, roleID int, permissionID int) error {
	query := `
	INSERT INTO role_permissions(
		role_id,
		permission_id
	)
	VALUES($1,$2)
	ON CONFLICT DO NOTHING
	`
	cmd, err := r.db.Exec(ctx, query, roleID, permissionID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrConflict
	}
	return nil
}

func (r *rolePermissionRepository) RemovePermission(ctx context.Context, roleID int, permissionID int) error {
	query := `
	DELETE FROM role_permissions
	WHERE role_id=$1
	AND permission_id=$2
	`
	cmd, err := r.db.Exec(ctx, query, roleID, permissionID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
