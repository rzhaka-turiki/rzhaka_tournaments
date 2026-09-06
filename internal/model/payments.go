package model

import (
	"time"

	"github.com/google/uuid"
)

type OrganisationSubscription struct {
	ID uuid.UUID

	OrganisationID uuid.UUID
	SubPlan        string
	Status         string

	StartedAt time.Time
	ExpiresAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrganisationPayment struct {
	ID uuid.UUID

	OrganisationID uuid.UUID
	SubscriptionID *uuid.UUID

	Amount   int
	Currency string

	Status string

	PaidAt  time.Time
	Comment string

	CreatedAt time.Time
	UpdatedAt time.Time
}
