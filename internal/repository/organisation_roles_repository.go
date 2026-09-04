package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type OrganisationRolesRepository interface {
	Create(ctx context.Context, role *model.OrganisationRole) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.OrganisationRole, error)
	GetAllRoles(ctx context.Context) ([]model.OrganisationRole, error)
	Update(ctx context.Context, role *model.OrganisationRole) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type organisationRolesRepository struct {
	db DBTX
}

func NewOrganisationRoleRepository(db DBTX) OrganisationRolesRepository {
	return &organisationRolesRepository{
		db: db,
	}
}

func (r *organisationRolesRepository) Create(ctx context.Context, role *model.OrganisationRole) error {
	query := `
	INSERT INTO organisation_roles (
		name,
		role_color
	) VALUES ($1, $2)
	RETURNING
		id,
		created_at,
		updated_at
	`

	return r.db.QueryRow(ctx, query, role.Name, role.RoleColor).Scan(
		&role.ID,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
}

func (r *organisationRolesRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.OrganisationRole, error) {
	query := `
	SELECT
		id,
		name,
		role_color,
		created_at,
		updated_at
	FROM organisation_roles
	WHERE id = $1
	`
	var role model.OrganisationRole
	err := r.db.QueryRow(ctx, query, id).Scan(
		&role.ID,
		&role.Name,
		&role.RoleColor,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *organisationRolesRepository) Update(ctx context.Context, role *model.OrganisationRole) error {
	query := `
	UPDATE organisation_roles
	SET
		name = $1,
		role_color = $2,
		updated_at = NOW()
	WHERE id = $3
	`
	result, err := r.db.Exec(ctx, query,
		role.Name,
		role.RoleColor,
		role.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *organisationRolesRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM organisation_roles
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

func (r *organisationRolesRepository) GetAllRoles(ctx context.Context) ([]model.OrganisationRole, error) {
	query := `
	SELECT
		id,
		name,
		role_color,
		creatat_at,
		update_at
	FROM organisation_role
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []model.OrganisationRole
	for rows.Next() {
		var role model.OrganisationRole
		err = rows.Scan(
			&role.ID,
			&role.Name,
			&role.RoleColor,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}
