package service

import (
	"context"

	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/a1uka/rzhaka_tournaments/internal/repository"
	"github.com/google/uuid"
)

type RoleService interface {
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error)
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]model.Permission, error)
	HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error)
	AssignRole(ctx context.Context, actorID uuid.UUID, targetUserID uuid.UUID, roleID int) error
	RemoveRole(ctx context.Context, actorID uuid.UUID, targetUserID uuid.UUID, roleID int) error
}

type roleService struct {
	roleRepository repository.RoleRepository
}

func NewRoleService(roleRepository repository.RoleRepository) RoleService {
	return &roleService{
		roleRepository: roleRepository,
	}
}

func (s *roleService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error) {
	return s.roleRepository.GetUserRoles(ctx, userID)
}

func (s *roleService) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]model.Permission, error) {
	return s.roleRepository.GetUserPermissions(ctx, userID)
}

func (s *roleService) HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) {
	return s.roleRepository.HasPermission(ctx, userID, permission)
}

func (s *roleService) AssignRole(ctx context.Context, actorID uuid.UUID, targetUserID uuid.UUID, roleID int) error {
	actorRole, err := s.roleRepository.GetHighestUserRole(ctx, actorID)
	if err != nil {
		return err
	}
	targetRole, err := s.roleRepository.GetHighestUserRole(ctx, targetUserID)
	if err != nil {
		return err
	}

	if actorRole.Position <= targetRole.Position {
		return repository.ErrForbidden
	}
	return s.roleRepository.AddUserRole(ctx, targetUserID, roleID)
}

func (s *roleService) RemoveRole(ctx context.Context, actorID uuid.UUID, targetUserID uuid.UUID, roleID int) error {
	actorRole, err := s.roleRepository.GetHighestUserRole(ctx, actorID)
	if err != nil {
		return err
	}
	targetRole, err := s.roleRepository.GetHighestUserRole(ctx, targetUserID)
	if err != nil {
		return err
	}

	if actorRole.Position <= targetRole.Position {
		return repository.ErrForbidden
	}
	return s.roleRepository.RemoveUserRole(ctx, targetUserID, roleID)
}
