package service

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/database"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/permission"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/repository"
)

type RolePermissionService interface {
	GetRolePermissions(ctx context.Context, roleID int) ([]model.Permission, error)
	AddPermission(ctx context.Context, actorID uuid.UUID, roleID int, permissionID int) error
	RemovePermission(ctx context.Context, actorID uuid.UUID, roleID int, permissionID int) error
}

type rolePermissionService struct {
	txManager                *database.TxManager
	roleRepository           repository.RoleRepository
	permissionRepository     repository.PermissionRepository
	rolePermissionRepository repository.RolePermissionRepository
	eventRepository          repository.EventRepository
}

func NewRolePermissionService(
	txManager *database.TxManager,
	roleRepository repository.RoleRepository,
	permissionRepository repository.PermissionRepository,
	rolePermissionRepository repository.RolePermissionRepository,
	eventRepository repository.EventRepository,
) RolePermissionService {
	return &rolePermissionService{
		txManager:                txManager,
		roleRepository:           roleRepository,
		permissionRepository:     permissionRepository,
		rolePermissionRepository: rolePermissionRepository,
		eventRepository:          eventRepository,
	}
}

func (s *rolePermissionService) GetRolePermissions(ctx context.Context, roleID int) ([]model.Permission, error) {
	return s.rolePermissionRepository.GetRolePermissions(ctx, roleID)
}

func (s *rolePermissionService) AddPermission(ctx context.Context, actorID uuid.UUID, roleID int, permissionID int) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		roleRepo := repository.NewRoleRepository(tx)
		permissionRepo := repository.NewPermissionRepository(tx)
		rolePermissionRepo := repository.NewRolePermissionRepository(tx)
		eventRepo := repository.NewEventRepository(tx)
		if err := CheckPermission(ctx, roleRepo, actorID, permission.RolePermissionsManage); err != nil {
			return err
		}
		role, err := roleRepo.GetRoleByID(ctx, roleID)
		if err != nil {
			return err
		}
		if err := CanManageRole(ctx, roleRepo, actorID, role); err != nil {
			return err
		}
		permissionExists, err := s.permissionRepository.Exists(ctx, permissionID)
		if err != nil {
			return err
		}
		if !permissionExists {
			return repository.ErrNotFound
		}
		perm, err := permissionRepo.GetByID(ctx, permissionID)
		if err != nil {
			return err
		}
		err = rolePermissionRepo.AddPermission(ctx, role.ID, perm.ID)
		if err != nil {
			return err
		}
		payload, err := json.Marshal(
			map[string]interface{}{
				"role_id":         role.ID,
				"role_name":       role.Name,
				"permission_id":   perm.ID,
				"permission_code": perm.Code,
			},
		)
		if err != nil {
			return err
		}
		event := &model.Event{
			ActorID:   &actorID,
			EventType: "ROLE_PERMISSION_ADDED",
			Payload:   payload,
		}
		return eventRepo.Create(ctx, event)
	})
}

func (s *rolePermissionService) RemovePermission(ctx context.Context, actorID uuid.UUID, roleID int, permissionID int) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		roleRepo := repository.NewRoleRepository(tx)
		permissionRepo := repository.NewPermissionRepository(tx)
		rolePermissionRepo := repository.NewRolePermissionRepository(tx)
		eventRepo := repository.NewEventRepository(tx)
		if err := CheckPermission(ctx, roleRepo, actorID, permission.RolePermissionsManage); err != nil {
			return err
		}
		role, err := roleRepo.GetRoleByID(ctx, roleID)
		if err != nil {
			return err
		}
		if err := CanManageRole(ctx, roleRepo, actorID, role); err != nil {
			return err
		}
		perm, err := permissionRepo.GetByID(ctx, permissionID)
		if err != nil {
			return err
		}
		err = rolePermissionRepo.RemovePermission(ctx, role.ID, perm.ID)
		if err != nil {
			return err
		}
		payload, err := json.Marshal(
			map[string]interface{}{
				"role_id":         role.ID,
				"role_name":       role.Name,
				"permission_id":   perm.ID,
				"permission_code": perm.Code,
			},
		)
		if err != nil {
			return err
		}
		event := &model.Event{
			ActorID:   &actorID,
			EventType: "ROLE_PERMISSION_REMOVED",
			Payload:   payload,
		}
		return eventRepo.Create(ctx, event)
	})
}
