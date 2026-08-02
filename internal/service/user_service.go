package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/a1uka/rzhaka_tournaments/internal/repository"
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/dto"
)

type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByDiscordHash(ctx context.Context, hash []byte) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	GetMe(ctx context.Context, userID uuid.UUID) (*dto.MeResponse, error)
}

type userService struct {
	userRepository repository.UserRepository
	roleRepository repository.RoleRepository
}

func NewUserService(userRepository repository.UserRepository, roleRepository repository.RoleRepository) UserService {
	return &userService{
		userRepository: userRepository,
		roleRepository: roleRepository,
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

func (s *userService) GetMe(ctx context.Context, userID uuid.UUID) (*dto.MeResponse, error) {
	roles, err := s.roleRepository.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	permissions, err := s.roleRepository.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}
	codes := make([]string, 0, len(permissions))
	for _, permission := range permissions {
		codes = append(codes, permission.Code)
	}
	return &dto.MeResponse{
		ID:          userID,
		Roles:       dto.FromRoles(roles),
		Permissions: codes,
	}, nil
}
