package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type ApexAccountRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.ApexAccount, error)
	GetByUID(ctx context.Context, uid string) (*model.ApexAccount, error)
	GetByNIDHash(ctx context.Context, nidHash string) (*model.ApexAccount, error)
	Create(ctx context.Context, account *model.ApexAccount) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type apexAccountRepository struct {
	db DBTX
}

func NewApexAccountRepository(db DBTX) ApexAccountRepository {
	return &apexAccountRepository{
		db: db,
	}
}

func (r *apexAccountRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.ApexAccount, error) {
	query := `
	SELECT
		id,
		user_id,
		uid,
		nid_hash,
		created_at,
		updated_at
	FROM apex_accounts
	WHERE user_id = $1
	`

	var account model.ApexAccount
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&account.ID,
		&account.UserID,
		&account.UID,
		&account.NIDHash,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *apexAccountRepository) GetByUID(ctx context.Context, uid string) (*model.ApexAccount, error) {
	query := `
	SELECT
		id,
		user_id,
		uid,
		nid_hash,
		created_at,
		updated_at
	FROM apex_accounts
	WHERE uid = $1
	`

	var account model.ApexAccount
	err := r.db.QueryRow(ctx, query, uid).Scan(
		&account.ID,
		&account.UserID,
		&account.UID,
		&account.NIDHash,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *apexAccountRepository) GetByNIDHash(ctx context.Context, nidHash string) (*model.ApexAccount, error) {
	query := `
	SELECT
		id,
		user_id,
		uid,
		nid_hash,
		created_at,
		updated_at
	FROM apex_accounts
	WHERE nid_hash = $1
	`

	var account model.ApexAccount
	err := r.db.QueryRow(ctx, query, nidHash).Scan(
		&account.ID,
		&account.UserID,
		&account.UID,
		&account.NIDHash,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *apexAccountRepository) Create(ctx context.Context, account *model.ApexAccount) error {
	query := `
	INSERT INTO apex_accounts (
		user_id,
		uid,
		nid_hash
	)
	VALUES ($1, $2, $3)
	RETURNING
		id,
		created_at,
		updated_at
	`
	return r.db.QueryRow(
		ctx,
		query,
		account.UserID,
		account.UID,
		account.NIDHash,
	).Scan(
		&account.ID,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
}

func (r *apexAccountRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `
	DELETE FROM apex_accounts
	WHERE user_id = $1
	`
	cmd, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
