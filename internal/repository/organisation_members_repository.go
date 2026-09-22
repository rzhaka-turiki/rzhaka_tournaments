package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type OrganisationMembersRepository interface {
	Create(ctx context.Context, member *model.OrganisationMember) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.OrganisationMember, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.OrganisationMember, error)
	GetByOrganisation(ctx context.Context, organisationID uuid.UUID) ([]model.OrganisationMember, error)
	HasRole(ctx context.Context, organisationID, userID uuid.UUID, roleCode string) (bool, error)
	GetRole(ctx context.Context, organisationID, userId uuid.UUID) (*string, error)
	GetByMember(ctx context.Context, memberID uuid.UUID) ([]model.OrganisationMember, error)
	Update(ctx context.Context, member *model.OrganisationMember) error
	UpdateRole(ctx context.Context, organisationID, userID, roleID uuid.UUID) error
	Exists(ctx context.Context, userID, organisationID uuid.UUID) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Remove(ctx context.Context, organisationID, userID uuid.UUID) error
}

type organisationMembersRepository struct {
	db DBTX
}

func NewOrganisationMembersRepository(db DBTX) OrganisationMembersRepository {
	return &organisationMembersRepository{
		db: db,
	}
}

func (r *organisationMembersRepository) Create(ctx context.Context, member *model.OrganisationMember) error {
	query := `
	INSERT INTO organisation_members (
		user_id,
		role_id,
		organisation_id,
		added_by,
	) VALUES ($1, $2, $3, $4)
	RETURNING
		id,
		created_at,
		updated_at
	`

	return r.db.QueryRow(ctx, query, member.UserID, member.RoleID, member.OrganisationID, member.AddedBy).Scan(
		&member.ID,
		&member.CreatedAt,
		&member.UpdatedAt,
	)
}

func (r *organisationMembersRepository) Remove(ctx context.Context, organisationID, userID uuid.UUID) error {
	query := `
	DELETE FROM organisation_members
	WHERE organisation_id = $1
		AND user_id = $2
	`
	result, err := r.db.Exec(ctx, query, organisationID, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *organisationMembersRepository) GetRole(ctx context.Context, organisationID, userId uuid.UUID) (*string, error) {
	query := `
	SELECT orole.code
	FROM organisation_roles AS orole
	INNER JOIN organisation_members AS omem 
    	ON orole.id = omem.role_id
	WHERE omem.organisation_id = $1
  		AND omem.user_id = $2
	`
	var roleCode string
	err := r.db.QueryRow(ctx, query, organisationID, userId).Scan(&roleCode)
	if err != nil {
		return nil, err
	}
	return &roleCode, nil
}

func (r *organisationMembersRepository) Exists(ctx context.Context, userID, organisationID uuid.UUID) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM organisation_members
		WHERE user_id = $1
			AND organisation_id = $2
	)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, organisationID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *organisationMembersRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.OrganisationMember, error) {
	query := `
	SELECT
		id,
		user_id,
		role_id,
		organisation_id,
		created_at,
		updated_at,
		added_by
	FROM organisation_members
	WHERE id = $1
	`
	var member model.OrganisationMember
	err := r.db.QueryRow(ctx, query, id).Scan(
		&member.ID,
		&member.UserID,
		&member.RoleID,
		&member.OrganisationID,
		&member.CreatedAt,
		&member.UpdatedAt,
		&member.AddedBy,
	)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *organisationMembersRepository) GetAll(ctx context.Context, limit, offset int) ([]model.OrganisationMember, error) {
	query := `
	SELECT
		id,
		user_id,
		organisation_id,
		role_id,
		created_at,
		updated_at
	FROM organisation_members
	ORDER BY created_at
	LIMIT $1
	OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []model.OrganisationMember
	for rows.Next() {
		var member model.OrganisationMember
		err = rows.Scan(
			&member.ID,
			&member.UserID,
			&member.OrganisationID,
			&member.RoleID,
			&member.CreatedAt,
			&member.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

func (r *organisationMembersRepository) HasRole(ctx context.Context, organisationID, userID uuid.UUID, roleCode string) (bool, error) {
	query := `
	SELECT EXISTS (
    	SELECT 1 
    	FROM organisation_members m
    	JOIN organisation_roles r ON m.role_id = r.id
    	WHERE m.user_id = $1 
    		AND m.organisation_id = $2
      		AND r.code = $3
	) AS has_role;
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, userID, organisationID, roleCode).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *organisationMembersRepository) UpdateRole(ctx context.Context, organisationID, userID, roleID uuid.UUID) error {
	query := `
	UPDATE organisation_members
	SET
		role_id = $1
	WHERE user_id = $2
		AND organisation_id = $3
	`
	result, err := r.db.Exec(ctx, query, roleID, userID, organisationID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *organisationMembersRepository) GetByOrganisation(ctx context.Context, organisationID uuid.UUID) ([]model.OrganisationMember, error) {
	query := `
	SELECT
		id,
		user_id,
		organisation_id,
		role_id,
		created_at,
		updated_at,
		added_by
	FROM organisation_members
	WHERE organisation_id = $1
	`
	rows, err := r.db.Query(ctx, query, organisationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []model.OrganisationMember
	for rows.Next() {
		var member model.OrganisationMember
		err = rows.Scan(
			&member.ID,
			&member.UserID,
			&member.OrganisationID,
			&member.RoleID,
			&member.CreatedAt,
			&member.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

func (r *organisationMembersRepository) GetByMember(ctx context.Context, memberID uuid.UUID) ([]model.OrganisationMember, error) {
	query := `
	SELECT
		id,
		user_id,
		organisation_id,
		role_id,
		created_at,
		updated_at,
		added_by
	FROM organisation_members
	WHERE user_id = $1
	`
	rows, err := r.db.Query(ctx, query, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []model.OrganisationMember
	for rows.Next() {
		var member model.OrganisationMember
		err = rows.Scan(
			&member.ID,
			&member.UserID,
			&member.OrganisationID,
			&member.RoleID,
			&member.CreatedAt,
			&member.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

func (r *organisationMembersRepository) Update(ctx context.Context, member *model.OrganisationMember) error {
	query := `
	UPDATE organisation_members
	SET
		user_id = $1,
		organisation_id = $2,
		role_id = $3,
		added_by = $4,
		updated_at = NOW()
	WHERE id = $5
	`
	result, err := r.db.Exec(ctx, query,
		member.UserID,
		member.OrganisationID,
		member.RoleID,
		member.AddedBy,
		member.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *organisationMembersRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE
	FROM organisation_members
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
