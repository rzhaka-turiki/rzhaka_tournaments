package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type OrganisationSubscriptionsRepository interface {
	Create(ctx context.Context, subscription *model.OrganisationSubscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.OrganisationSubscription, error)
	GetByOrganisationID(ctx context.Context, organisationID uuid.UUID, isActive bool) ([]model.OrganisationSubscription, error)
	GetActive(ctx context.Context) ([]model.OrganisationSubscription, error)
	GetByPlan(ctx context.Context, subPlan string) ([]model.OrganisationSubscription, error)
	Update(ctx context.Context, subscription *model.OrganisationSubscription) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type organisationSubscriptionRepository struct {
	db DBTX
}

func NewOrganisationSubscriptionRepository(db DBTX) OrganisationSubscriptionsRepository {
	return &organisationSubscriptionRepository{
		db: db,
	}
}

func (r *organisationSubscriptionRepository) Create(ctx context.Context, subscription *model.OrganisationSubscription) error {
	query := `
	INSERT INTO organisation_subscriptions (
		organisation_id,
		sub_plan,
		status,
		started_at,
		expires_at
	) VALUES ($1, $2, $3, $4, $5)
	RETURNING
		id,
		created_at,
		updated_at
	`
	return r.db.QueryRow(ctx, query,
		subscription.OrganisationID,
		subscription.SubPlan,
		subscription.Status,
		subscription.StartedAt,
		subscription.ExpiresAt,
	).Scan(
		&subscription.ID,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
}

func (r *organisationSubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.OrganisationSubscription, error) {
	query := `
	SELECT
		id,
		organisation_id,
		sub_plan,
		status,
		started_at,
		expires_at,
		created_at,
		updated_at
	FROM organisation_subscriptions
	WHERE id = $1
	`
	var sub model.OrganisationSubscription
	err := r.db.QueryRow(ctx, query, id).Scan(
		&sub.ID,
		&sub.OrganisationID,
		&sub.SubPlan,
		&sub.Status,
		&sub.StartedAt,
		&sub.ExpiresAt,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *organisationSubscriptionRepository) GetByOrganisationID(ctx context.Context, organisationID uuid.UUID, isActive bool) ([]model.OrganisationSubscription, error) {
	query := `
	SELECT
		id,
		organisation_id,
		sub_plan,
		status,
		started_at,
		expires_at,
		created_at,
		updated_at
	FROM organisation_subscriptions
	WHERE organisation_id = $1
	`

	if isActive {
		query += `
		AND expires_at > NOW() AND started_at <= NOW()
		`
	}
	rows, err := r.db.Query(ctx, query, organisationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subs []model.OrganisationSubscription
	for rows.Next() {
		var sub model.OrganisationSubscription
		err := rows.Scan(
			&sub.ID,
			&sub.OrganisationID,
			&sub.SubPlan,
			&sub.Status,
			&sub.StartedAt,
			&sub.ExpiresAt,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

func (r *organisationSubscriptionRepository) GetActive(ctx context.Context) ([]model.OrganisationSubscription, error) {
	query := `
	SELECT
		id,
		organisation_id,
		sub_plan,
		status,
		started_at,
		expires_at,
		created_at,
		updated_at
	FROM organisation_subscriptions
	WHERE expires_at > NOW()
		AND started_at <= NOW()
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subs []model.OrganisationSubscription
	for rows.Next() {
		var sub model.OrganisationSubscription
		err := rows.Scan(
			&sub.ID,
			&sub.OrganisationID,
			&sub.SubPlan,
			&sub.Status,
			&sub.StartedAt,
			&sub.ExpiresAt,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

func (r *organisationSubscriptionRepository) GetByPlan(ctx context.Context, subPlan string) ([]model.OrganisationSubscription, error) {
	query := `
	SELECT
		id,
		organisation_id,
		sub_plan,
		status,
		started_at,
		expires_at,
		created_at,
		updated_at
	FROM organisation_subscriptions
	WHERE sub_plan = $1
	`

	rows, err := r.db.Query(ctx, query, subPlan)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subs []model.OrganisationSubscription
	for rows.Next() {
		var sub model.OrganisationSubscription
		err := rows.Scan(
			&sub.ID,
			&sub.OrganisationID,
			&sub.SubPlan,
			&sub.Status,
			&sub.StartedAt,
			&sub.ExpiresAt,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

func (r *organisationSubscriptionRepository) Update(ctx context.Context, subscription *model.OrganisationSubscription) error {
	query := `
	UPDATE organisation_subscriptions
	SET
		organisation_id = $1,
		sub_plan = $2,
		status = $3,
		started_at = $4,
		expires_at = $5
		updated_at = NOW()
	WHERE id = $6
	`

	result, err := r.db.Exec(ctx, query,
		subscription.OrganisationID,
		subscription.SubPlan,
		subscription.Status,
		subscription.StartedAt,
		subscription.ExpiresAt,
		subscription.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *organisationSubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE
	FROM organisation_subscriptions
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
