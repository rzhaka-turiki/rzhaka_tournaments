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
	GetByMember(ctx context.Context, memberID uuid.UUID) ([]model.OrganisationMember, error)
	Update(ctx context.Context, member *model.OrganisationMember) error
	Delete(ctx context.Context, id uuid.UUID) error
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
