package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type OrganisationsRepository interface {
	Create(ctx context.Context, organisation *model.Organisation) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Organisation, error)
	GetByOwner(ctx context.Context, ownerID uuid.UUID) ([]model.Organisation, error)
	GetMemberOrganisations(ctx context.Context, userID uuid.UUID) ([]model.Organisation, error)
	Update(ctx context.Context, organisation *model.Organisation) error
	IsOwner(ctx context.Context, userID, organisationID uuid.UUID) (bool, error)
	UpdateOwner(ctx context.Context, organisationID, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type organisationRepository struct {
	db DBTX
}

func NewOrganisationRepository(db DBTX) OrganisationsRepository {
	return &organisationRepository{
		db: db,
	}
}

func (r *organisationRepository) Create(ctx context.Context, organisation *model.Organisation) error {
	query := `
	INSERT INTO organisations (
		owner_id,
		name,
		short_name,
		image_url,
		banner_url
	) VALUES ($1, $2, $3, $4, $5)
	RETURNING
		id,
		created_at,
		updated_at
	`

	return r.db.QueryRow(ctx, query, organisation.OwnerID, organisation.Name, organisation.ShortName, organisation.ImageURL, organisation.BannerURL).Scan(
		&organisation.ID,
		&organisation.CreatedAt,
		&organisation.UpdatedAt,
	)
}

func (r *organisationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Organisation, error) {
	query := `
	SELECT
		id,
		owner_id,
		name,
		short_name,
		image_url,
		banner_url,
		created_at,
		updated_at
	FROM organisations
	WHERE id = $1
	`
	var organisation model.Organisation
	err := r.db.QueryRow(ctx, query, id).Scan(
		&organisation.ID,
		&organisation.OwnerID,
		&organisation.Name,
		&organisation.ShortName,
		&organisation.ImageURL,
		&organisation.BannerURL,
		&organisation.CreatedAt,
		&organisation.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &organisation, nil
}

func (r *organisationRepository) IsOwner(ctx context.Context, userID, organisationID uuid.UUID) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM organisations
		WHERE id = $1
			AND owner_id = $2
	)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, organisationID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *organisationRepository) GetMemberOrganisations(ctx context.Context, userID uuid.UUID) ([]model.Organisation, error) {
	query := `
	SELECT
		o.id,
		o.owner_id,
		o.name,
		o.short_name,
		o.image_url,
		o.banner_url,
		o.created_at,
		o.updated_at
	FROM organisations o
	JOIN organisation_members om
		ON om.organisation_id = o.id
	WHERE om.user_id = $1
	ORDER BY o.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var organisations []model.Organisation
	for rows.Next() {
		var organisation model.Organisation
		err := rows.Scan(
			&organisation.ID,
			&organisation.OwnerID,
			&organisation.Name,
			&organisation.ShortName,
			&organisation.ImageURL,
			&organisation.BannerURL,
			&organisation.CreatedAt,
			&organisation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		organisations = append(organisations, organisation)
	}
	return organisations, nil
}

func (r *organisationRepository) UpdateOwner(ctx context.Context, organisationID, userID uuid.UUID) error {
	query := `
	UPDATE organisations
	SET
		owner_id = $1,
		updated_at = NOW()
	WHERE id = $2
	`
	result, err := r.db.Exec(ctx, query, userID, organisationID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *organisationRepository) GetByOwner(ctx context.Context, ownerID uuid.UUID) ([]model.Organisation, error) {
	query := `
	SELECT
		id,
		owner_id,
		name,
		short_name,
		image_url,
		banner_url,
		created_at,
		updated_at
	FROM organisations
	WHERE owner_id = $1
	`

	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var organisations []model.Organisation
	for rows.Next() {
		var organisation model.Organisation
		err := rows.Scan(
			&organisation.ID,
			&organisation.OwnerID,
			&organisation.Name,
			&organisation.ShortName,
			&organisation.ImageURL,
			&organisation.BannerURL,
			&organisation.CreatedAt,
			&organisation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		organisations = append(organisations, organisation)
	}
	return organisations, nil
}

func (r *organisationRepository) Update(ctx context.Context, organisation *model.Organisation) error {
	query := `
	UPDATE organisations 
	SET
		name = $1,
		short_name = $2,
		image_url = $3,
		banner_url = $4,
		owner_id = $5,
		updated_at = NOW()
	WHERE id = $6
	`

	result, err := r.db.Exec(ctx, query,
		organisation.Name,
		organisation.ShortName,
		organisation.ImageURL,
		organisation.BannerURL,
		organisation.OwnerID,
		organisation.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *organisationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE
	FROM organisations 
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
