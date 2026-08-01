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
	GetHighestUserRole(ctx context.Context, userID uuid.UUID) (*model.Role, error)
	GetAllRoles(ctx context.Context) ([]model.Role, error)
	GetAllPermissions(ctx context.Context) ([]model.Permission, error)
	GetRoleByIDIncludeDeleted(ctx context.Context, ID int) (*model.Role, error)
	// misc
	HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error)
	AddUserRole(ctx context.Context, userID uuid.UUID, roleID int) error
	RemoveUserRole(ctx context.Context, userID uuid.UUID, roleID int) error

	// Role management
	Create(ctx context.Context, role *model.Role) error
	SoftDelete(ctx context.Context, roleID int, deletedBy uuid.UUID) error
	Restore(ctx context.Context, roleID int) error
}

type RolePermissionRepository interface {
	GetRolePermissions(ctx context.Context, roleID int) ([]model.Permission, error)
	Add(ctx context.Context, roleID int, permissionID int) error
	Remove(ctx context.Context, roleID int, permissionID int) error
}

type roleRepository struct {
	db DBTX
}

type rolePermissionRepository struct {
	db DBTX
}

func NewRoleRepository(db DBTX) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func NewRolePermissionRepository(db DBTX) RolePermissionRepository {
	return &rolePermissionRepository{
		db: db,
	}
}

func (r *roleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error) {
	query := `
		SELECT
			r.id,
			r.name,
			r.role_color
			r.deleted_at
		FROM roles r
		JOIN user_roles ur
			ON ur.role_id = r.id
		WHERE ur.user_id = $1 AND r.deleted_at IS NULL
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
			deleted_at
		FROM roles
		WHERE id = $1
		AND deleted_at IS NULL
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

func (r *roleRepository) GetHighestUserRole(ctx context.Context, userID uuid.UUID) (*model.Role, error) {
	query := `
		SELECT
			r.id,
			r.name,
			r.role_color,
			r.position,
			r.is_system,
			r.created_at
			r.deleted_at
		FROM roles r
		JOIN user_roles ur
			ON ur.role_id = r.id
		WHERE ur.user_id = $1
		AND deleted_at IS NULL
		ORDER BY r.position DESC
		LIMIT 1
	`
	var role model.Role
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&role.ID,
		&role.Name,
		&role.RoleColor,
		&role.Position,
		&role.IsSystem,
		&role.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	query := `
	SELECT
    	id,
    	name,
    	position,
    	role_color
		deleted_at
	FROM roles
	WHERE deleted_at IS NULL
	ORDER BY position DESC;
	`

	rows, err := r.db.Query(ctx, query)
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
			&role.Position,
			&role.RoleColor,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (r *roleRepository) GetAllPermissions(ctx context.Context) ([]model.Permission, error) {
	query := `
	SELECT
    	id,
    	code, 
    	description
		deleted_at
	FROM permissions
	WHERE deleted_at IS NULL
	ORDER BY code;
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

func (r *roleRepository) Create(ctx context.Context, role *model.Role) error {
	query := `
	INSERT INTO roles (
		name,
		position,
		role_color
	)
	VALUES ($1,$2,$3)
	RETURNING id
	`
	return r.db.QueryRow(ctx, query, role.Name, role.Position, role.RoleColor).Scan(&role.ID)
}

func (r *roleRepository) SoftDelete(ctx context.Context, roleID int, deletedBy uuid.UUID) error {
	query := `
		UPDATE roles
		SET
			deleted_at = NOW(),
			deleted_by = $2
		WHERE id = $1
		AND deleted_at IS NULL
	`
	cmd, err := r.db.Exec(ctx, query, roleID, deletedBy)

	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *roleRepository) Restore(ctx context.Context, roleID int) error {
	query := `
		UPDATE roles
		SET
			deleted_at = NULL,
			deleted_by = NULL
		WHERE id = $1
		AND deleted_at IS NOT NULL
	`
	cmd, err := r.db.Exec(ctx, query, roleID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *roleRepository) GetRoleByIDIncludeDeleted(ctx context.Context, ID int) (*model.Role, error) {
	query := `
		SELECT
			id,
			name,
			position,
			role_color,
			system,
			created_at,
			deleted_at,
			deleted_by
		FROM roles
		WHERE id = $1
	`
	var role model.Role

	err := r.db.QueryRow(ctx, query, ID).Scan(
		&role.ID,
		&role.Name,
		&role.Position,
		&role.RoleColor,
		&role.IsSystem,
		&role.CreatedAt,
		&role.DeletedAt,
		&role.DeletedBy,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
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

func (r *rolePermissionRepository) Add(ctx context.Context, roleID int, permissionID int) error {
	query := `
	INSERT INTO role_permissions(
		role_id,
		permission_id
	)
	VALUES($1,$2)
	ON CONFLICT DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, roleID, permissionID)
	return err
}

func (r *rolePermissionRepository) Remove(ctx context.Context, roleID int, permissionID int) error {
	query := `
	DELETE FROM role_permissions
	WHERE role_id=$1
	AND permission_id=$2
	`
	_, err := r.db.Exec(ctx, query, roleID, permissionID)

	return err
}
