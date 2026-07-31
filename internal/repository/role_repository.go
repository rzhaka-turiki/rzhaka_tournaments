package repository

import (
	"context"
	"errors"

	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RoleRepository interface {
	// get
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error)
	GetRoleByID(ctx context.Context, id int) (*model.Role, error)
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]model.Permission, error)
	// misc
	HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error)
	AddUserRole(ctx context.Context, userID uuid.UUID, roleID int) error
	RemoveUserRole(ctx context.Context, userID uuid.UUID, roleID int) error
}

type roleRepository struct {
	db DBTX
}

func NewRoleRepository(db DBTX) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error) {
	query := `
		SELECT
			r.id,
			r.name,
			r.role_color
		FROM roles r
		JOIN user_roles ur
			ON ur.role_id = r.id
		WHERE ur.user_id = $1
		ORDER BY r.id
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []model.Role

	for rows.Next() {
		var role model.Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.RoleColor,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (r *roleRepository) AddUserRole(ctx context.Context, userID uuid.UUID, roleID int) error {
	query := `
		INSERT INTO user_roles (
			user_id,
			role_id
		)
		VALUES ($1,$2)
		ON CONFLICT DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, userID, roleID)
	return err
}

func (r *roleRepository) RemoveUserRole(ctx context.Context, userID uuid.UUID, roleID int) error {
	query := `
		DELETE FROM user_roles
		WHERE user_id = $1
		AND role_id = $2
	`
	_, err := r.db.Exec(ctx, query, userID, roleID)
	return err
}

func (r *roleRepository) GetRoleByID(ctx context.Context, id int) (*model.Role, error) {
	query := `
		SELECT
			id,
			name,
			role_color
		FROM roles
		WHERE id = $1
	`
	var role model.Role
	err := r.db.QueryRow(ctx, query, id).Scan(&role.ID, &role.Name, &role.RoleColor)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]model.Permission, error) {
	query := `
		SELECT DISTINCT
			p.id,
			p.code,
			p.description
		FROM permissions p
		JOIN role_permissions rp
			ON rp.permission_id = p.id
		JOIN user_roles ur
			ON ur.role_id = rp.role_id
		WHERE ur.user_id = $1
		ORDER BY p.id
	`
	rows, err := r.db.Query(ctx, query, userID)
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

func (r *roleRepository) HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM permissions p
			JOIN role_permissions rp
				ON rp.permission_id = p.id
			JOIN user_roles ur
				ON ur.role_id = rp.role_id
			WHERE ur.user_id = $1
			AND p.code = $2
		)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, permission).Scan(&exists)
	return exists, err
}
