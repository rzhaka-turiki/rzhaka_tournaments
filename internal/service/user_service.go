package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/a1uka/rzhaka_tournaments/internal/repository"
)

type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByDiscordHash(ctx context.Context, hash []byte) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) GetByDiscordHash(ctx context.Context, hash []byte) (*model.User, error) {
	user, err := s.userRepository.GetByDiscordHash(ctx, hash)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, user *model.User) error {
	return s.userRepository.Update(ctx, user)
}
