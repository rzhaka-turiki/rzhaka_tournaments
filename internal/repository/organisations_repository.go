package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type OrganisationsRepository interface {
	Create(ctx context.Context, organisation *model.Organisation) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Organisation, error)
	Update(ctx context.Context, organisation *model.Organisation) error
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
		name,
		short_name,
		image_url,
		banner_url
	) VALUES ($1, $2, $3, $4)
	RETURNING
		id,
		created_at,
		updated_at
	`

	return r.db.QueryRow(ctx, query, organisation.Name, organisation.ShortName, organisation.ImageURL, organisation.BannerURL).Scan(
		&organisation.ID,
		&organisation.CreatedAt,
		&organisation.UpdatedAt,
	)
}

func (r *organisationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Organisation, error) {
	query := `
	SELECT
		id,
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

func (r *organisationRepository) Update(ctx context.Context, organisation *model.Organisation) error {
	query := `
	UPDATE organisations 
	SET
		name = $1,
		short_name = $2,
		image_url = $3,
		banner_url = $4,
		updated_at = NOW()
	WHERE id = $5
	`

	result, err := r.db.Exec(ctx, query,
		organisation.Name,
		organisation.ShortName,
		organisation.ImageURL,
		organisation.BannerURL,
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
