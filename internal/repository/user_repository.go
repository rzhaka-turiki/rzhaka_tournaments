package repository

import (
	"context"
	"errors"

	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByDiscordHash(ctx context.Context, hash []byte) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
}

type userRepository struct {
	db DBTX
}

func NewUserRepository(db DBTX) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
		SELECT
			id,
			discord_id_encrypted,
			discord_id_hash,
			username,
			avatar_url,
			status,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`
	var user model.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.DiscordIDEncrypted,
		&user.DiscordIDHash,
		&user.Username,
		&user.AvatarURL,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByDiscordHash(ctx context.Context, hash []byte) (*model.User, error) {
	query := `
		SELECT
			id,
			discord_id_encrypted,
			discord_id_hash,
			username,
			avatar_url,
			status,
			created_at,
			updated_at
		FROM users
		WHERE discord_id_hash = $1
	`
	var user model.User
	err := r.db.QueryRow(ctx, query, hash).Scan(
		&user.ID,
		&user.DiscordIDEncrypted,
		&user.DiscordIDHash,
		&user.Username,
		&user.AvatarURL,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (
			discord_id_encrypted,
			discord_id_hash,
			username,
			avatar_url
		)
		VALUES ($1, $2, $3, $4)
		RETURNING 
			id,
			created_at,
			updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		user.DiscordIDEncrypted,
		user.DiscordIDHash,
		user.Username,
		user.AvatarURL,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET
			username = $1,
			avatar_url = $2,
			updated_at = NOW()
		WHERE id = $3
	`
	cmd, err := r.db.Exec(ctx, query, user.Username, user.AvatarURL, user.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}
