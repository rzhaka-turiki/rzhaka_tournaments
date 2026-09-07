package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/shopspring/decimal"
)

type OrganisationPaymentsRepository interface {
	Create(ctx context.Context, payment *model.OrganisationPayment) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.OrganisationPayment, error)
	GetByOrganisationID(ctx context.Context, organisationID uuid.UUID) ([]model.OrganisationPayment, error)
	GetByCurrency(ctx context.Context, currency string) ([]model.OrganisationPayment, error)
	GetBySubscriptionID(ctx context.Context, subscriptionID uuid.UUID) (*model.OrganisationPayment, error)
	GetByStatus(ctx context.Context, status string) ([]model.OrganisationPayment, error)
	Update(ctx context.Context, payment *model.OrganisationPayment) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type organisationPaymentsRepository struct {
	db DBTX
}

func NewOrganisationPaymentsRepository(db DBTX) OrganisationPaymentsRepository {
	return &organisationPaymentsRepository{
		db: db,
	}
}

func (r *organisationPaymentsRepository) Create(ctx context.Context, payment *model.OrganisationPayment) error {
	query := `
	INSERT INTO organisation_payments (
		organisation_id,
		subscription_id,
		amount,
		currency,
		status,
		paid_at,
		comment
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING
		id,
		created_at,
		updated_at
	`
	return r.db.QueryRow(ctx, query,
		payment.OrganisationID,
		payment.SubscriptionID,
		payment.Amount.String(),
		payment.Currency,
		payment.Status,
		payment.PaidAt,
		payment.Comment,
	).Scan(
		&payment.ID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
}

func (r *organisationPaymentsRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.OrganisationPayment, error) {
	query := `
	SELECT
		id,
		organisation_id,
		subscription_id,
		amount,
		currency,
		status,
		paid_at,
		comment,
		created_at,
		updated_at
	FROM organisation_payments
	WHERE id = $1
	`
	var payment model.OrganisationPayment
	var amountStr string
	err := r.db.QueryRow(ctx, query, id).Scan(
		&payment.ID,
		&payment.OrganisationID,
		&payment.SubscriptionID,
		&amountStr,
		&payment.Currency,
		&payment.Status,
		&payment.PaidAt,
		&payment.Comment,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	payment.Amount, err = decimal.NewFromString(amountStr)
	if err != nil {
		return nil, fmt.Errorf("parse points: %w", err)
	}
	return &payment, nil
}

func (r *organisationPaymentsRepository) GetByOrganisationID(ctx context.Context, organisationID uuid.UUID) ([]model.OrganisationPayment, error) {
	query := `
	SELECT
		id,
		organisation_id,
		subscription_id,
		amount,
		currency,
		status,
		paid_at,
		comment,
		created_at,
		updated_at
	FROM organisation_payments
	WHERE organisation_id = $1
	`
	rows, err := r.db.Query(ctx, query, organisationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []model.OrganisationPayment
	for rows.Next() {
		var payment model.OrganisationPayment
		var amountStr string
		err := rows.Scan(
			&payment.ID,
			&payment.OrganisationID,
			&payment.SubscriptionID,
			&amountStr,
			&payment.Currency,
			&payment.Status,
			&payment.PaidAt,
			&payment.Comment,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payment.Amount, err = decimal.NewFromString(amountStr)
		if err != nil {
			return nil, fmt.Errorf("parse points: %w", err)
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (r *organisationPaymentsRepository) GetByCurrency(ctx context.Context, currency string) ([]model.OrganisationPayment, error) {
	query := `
	SELECT
		id,
		organisation_id,
		subscription_id,
		amount,
		currency,
		status,
		paid_at,
		comment,
		created_at,
		updated_at
	FROM organisation_payments
	WHERE currency = $1
	`
	rows, err := r.db.Query(ctx, query, currency)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []model.OrganisationPayment
	for rows.Next() {
		var payment model.OrganisationPayment
		var amountStr string
		err := rows.Scan(
			&payment.ID,
			&payment.OrganisationID,
			&payment.SubscriptionID,
			&amountStr,
			&payment.Currency,
			&payment.Status,
			&payment.PaidAt,
			&payment.Comment,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payment.Amount, err = decimal.NewFromString(amountStr)
		if err != nil {
			return nil, fmt.Errorf("parse points: %w", err)
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (r *organisationPaymentsRepository) GetBySubscriptionID(ctx context.Context, subscriptionID uuid.UUID) (*model.OrganisationPayment, error) {
	query := `
	SELECT
		id,
		organisation_id,
		subscription_id,
		amount,
		currency,
		status,
		paid_at,
		comment,
		created_at,
		updated_at
	FROM organisation_payments
	WHERE subscription_id = $1
	`
	var payment model.OrganisationPayment
	var amountStr string
	err := r.db.QueryRow(ctx, query, subscriptionID).Scan(
		&payment.ID,
		&payment.OrganisationID,
		&payment.SubscriptionID,
		&amountStr,
		&payment.Currency,
		&payment.Status,
		&payment.Status,
		&payment.PaidAt,
		&payment.Comment,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	payment.Amount, err = decimal.NewFromString(amountStr)
	if err != nil {
		return nil, fmt.Errorf("parse points: %w", err)
	}
	return &payment, nil
}

func (r *organisationPaymentsRepository) GetByStatus(ctx context.Context, status string) ([]model.OrganisationPayment, error) {
	query := `
	SELECT
		id,
		organisation_id,
		subscription_id,
		amount,
		currency,
		status,
		paid_at,
		comment,
		created_at,
		updated_at
	FROM organisation_payments
	WHERE status = $1
	`

	var payments []model.OrganisationPayment
	rows, err := r.db.Query(ctx, query, status)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var payment model.OrganisationPayment
		var amountStr string
		err := rows.Scan(
			&payment.ID,
			&payment.OrganisationID,
			&payment.SubscriptionID,
			&amountStr,
			&payment.Currency,
			&payment.Status,
			&payment.PaidAt,
			&payment.Comment,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (r *organisationPaymentsRepository) Update(ctx context.Context, payment *model.OrganisationPayment) error {
	query := `
	UPDATE organisation_payments
	SET
		organisation_id = $1,
		subscription_id = $2,
		amount = $3,
		currency = $4,
		status = $5,
		paid_at = $6,
		comment = $7,
	WHERE id = $8
	`

	result, err := r.db.Exec(ctx, query,
		payment.OrganisationID,
		payment.SubscriptionID,
		payment.Amount.String(),
		payment.Currency,
		payment.Status,
		payment.PaidAt,
		payment.Comment,
		payment.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *organisationPaymentsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE
	FROM organisation_payments
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
