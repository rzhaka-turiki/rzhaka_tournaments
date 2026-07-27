package repository

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	pool   *pgxpool.Pool
	encKey []byte
}

func NewUserRepo(pool *pgxpool.Pool, encKey []byte) UserRepository {
	return &userRepo{pool: pool, encKey: encKey}
}

func (r *userRepo) CreateOrUpdate(ctx context.Context, discordID string, user *model.User) error {
	// discordID crypt
	encrypted, err := encrypt.Encrypt(discordID, r.encKey)
	if err != nil {
		return fmt.Errorf("encrypt discord id: %w", err)
	}

	hash := sha256.Sum256([]byte(discordID))
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	query := `
	INSERT INTO users (id, discord_id_encrypted, discord_id_hash, username, avatar_url, role, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW())
        ON CONFLICT (discord_id_hash) DO UPDATE SET
            username = EXCLUDED.username,
            avatar_url = EXCLUDED.avatar_url
        RETURNING id, discord_id_encrypted, discord_id_hash, username, avatar_url, role, created_at
	`
	row := r.pool.QueryRow(ctx, query,
		user.ID,
		encrypted,
		hash[:],
		user.Username,
		user.AvatarURL,
		user.Role,
	)

	// res scan
	err = row.Scan(
		&user.ID,
		&user.DiscordIDEncrypted,
		&user.DiscordIDHash,
		&user.Username,
		&user.AvatarURL,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert/update user: %w", err)
	}

	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
	SELECT id, discord_id_encrypted, discord_id_hash, username, avatar_url, role, created_at
        FROM users
        WHERE id = $1
	`
	user := &model.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.DiscordIDEncrypted,
		&user.DiscordIDHash,
		&user.Username,
		&user.AvatarURL,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by id: %w", err)
	}
	return user, nil
}

func (r *userRepo) UpdateProfile(ctx context.Context, id uuid.UUID, username, avatarURL *string) error {
	query := `
	UPDATE users
        SET
            username = COALESCE($2, username),
            avatar_url = CASE WHEN $3 THEN $4 ELSE avatar_url END
        WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id, username, avatarURL != nil, avatarURL)
	return err
}
