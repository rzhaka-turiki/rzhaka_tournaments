package repository

import (
	"context"

	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	// Создает нового пользователя или обновляет текущего по DiscordID
	// Принимает discordID, user с заполненными полями
	// Возвращает ошибку, если не смог
	CreateOrUpdate(ctx context.Context, discordID string, user *model.User)

	// Ищет пользователя по SHA256 хешу DiscordID
	GetByDiscordHash(ctx context.Context, hash []byte) (*model.User, error)

	// Возвращает пользователя по UUID
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)

	// Обновляет отображаемые имя пользователя и аватарку
	UpdateProfile(ctx context.Context, id uuid.UUID, username, avatarURL *string) error
}
