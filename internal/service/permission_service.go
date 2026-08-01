package service

import (
	"context"

	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/a1uka/rzhaka_tournaments/internal/repository"
)

type PermissionService interface {
	GetAll(ctx context.Context) ([]model.Permission, error)
	GetByID(ctx context.Context, ID int) (*model.Permission, error)
	GetByCode(ctx context.Context, code string) (*model.Permission, error)
}

type permissionService struct {
	permissionRepository repository.PermissionRepository
}

func NewPermissionService(permissionRepository repository.PermissionRepository) PermissionService {
	return &permissionService{
		permissionRepository: permissionRepository,
	}
}

func (s *permissionService) GetAll(ctx context.Context) ([]model.Permission, error) {
	return s.permissionRepository.GetAll(ctx)
}

func (s *permissionService) GetByID(ctx context.Context, ID int) (*model.Permission, error) {
	return s.permissionRepository.GetByID(ctx, ID)
}

func (s *permissionService) GetByCode(ctx context.Context, code string) (*model.Permission, error) {
	return s.permissionRepository.GetByCode(ctx, code)
}
