package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type OrganisationSubscriptionsRepository interface {
	Create(ctx context.Context, subscription *model.OrganisationSubscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.OrganisationSubscription, error)
	GetByOrganisationID(ctx context.Context, organisationID uuid.UUID) ([]model.OrganisationSubscription, error)
	GetActive(ctx context.Context) ([]model.OrganisationSubscription, error)
	GetActiveByOrganisationID(ctx context.Context, organisationID uuid.UUID) ([]model.OrganisationSubscription, error)
	GetByPlan(ctx context.Context, SubPlan string) ([]model.OrganisationSubscription, error)
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
