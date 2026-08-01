package service

import (
	"context"
	"encoding/json"

	"github.com/a1uka/rzhaka_tournaments/internal/database"
	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/a1uka/rzhaka_tournaments/internal/permission"
	"github.com/a1uka/rzhaka_tournaments/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RoleService interface {
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error)
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]model.Permission, error)
	GetAllRoles(ctx context.Context) ([]model.Role, error)
	GetAllPermissions(ctx context.Context) ([]model.Permission, error)
	AssignRole(ctx context.Context, actorID uuid.UUID, targetUserID uuid.UUID, roleID int) error
	RemoveRole(ctx context.Context, actorID uuid.UUID, targetUserID uuid.UUID, roleID int) error

	CreateRole(ctx context.Context, actorID uuid.UUID, role *model.Role) error

	DeleteRole(ctx context.Context, actorID uuid.UUID, roleID int) error
	RestoreRole(ctx context.Context, actorID uuid.UUID, roleID int) error
}

type roleService struct {
	txManager                *database.TxManager
	roleRepository           repository.RoleRepository
	eventRepository          repository.EventRepository
	permissionRepository     repository.PermissionRepository
	rolePermissionRepository repository.RolePermissionRepository
}

func NewRoleService(
	txManager *database.TxManager,
	roleRepository repository.RoleRepository,
	eventRepository repository.EventRepository,
) RoleService {
	return &roleService{
		txManager:       txManager,
		roleRepository:  roleRepository,
		eventRepository: eventRepository,
	}
}

func (s *roleService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error) {
	return s.roleRepository.GetUserRoles(ctx, userID)
}

func (s *roleService) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]model.Permission, error) {
	return s.roleRepository.GetUserPermissions(ctx, userID)
}

func (s *roleService) canAssignRole(ctx context.Context, roleRepository repository.RoleRepository, actorID uuid.UUID, role *model.Role) error {
	if err := CheckPermission(ctx, roleRepository, actorID, "role.assign"); err != nil {
		return err
	}
	if err := CanManageRole(ctx, roleRepository, actorID, role); err != nil {
		return err
	}
	return nil
}

func (s *roleService) AssignRole(ctx context.Context, actorID uuid.UUID, targetUserID uuid.UUID, roleID int) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		// transactions for adding
		roleRepository := repository.NewRoleRepository(tx)
		eventRepository := repository.NewEventRepository(tx)
		// permission check
		targetRole, err := roleRepository.GetRoleByID(ctx, roleID)
		if err != nil {
			return err
		}
		if err := s.canAssignRole(ctx, roleRepository, actorID, targetRole); err != nil {
			return err
		}
		if err := roleRepository.AddUserRole(ctx, targetUserID, roleID); err != nil {
			return err
		}
		// add event
		payload, err := json.Marshal(map[string]any{
			"role": map[string]any{
				"id":   targetRole.ID,
				"name": targetRole.Name,
			},
		})
		if err != nil {
			return err
		}
		event := &model.Event{
			ActorID:   &actorID,
			UserID:    &targetUserID,
			EventType: "ROLE_ASSIGNED",
			Payload:   payload,
		}
		return eventRepository.Create(ctx, event)
	})
}

func (s *roleService) canRemoveRole(ctx context.Context, roleRepository repository.RoleRepository, actorID uuid.UUID, role *model.Role) error {
	return s.canAssignRole(ctx, roleRepository, actorID, role)
}

func (s *roleService) RemoveRole(ctx context.Context, actorID uuid.UUID, targetUserID uuid.UUID, roleID int) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		// transactions for adding
		roleRepository := repository.NewRoleRepository(tx)
		eventRepository := repository.NewEventRepository(tx)
		// permission check
		targetRole, err := roleRepository.GetRoleByID(ctx, roleID)
		if err != nil {
			return err
		}
		if err := s.canRemoveRole(ctx, roleRepository, actorID, targetRole); err != nil {
			return err
		}
		if err := roleRepository.RemoveUserRole(ctx, targetUserID, roleID); err != nil {
			return err
		}
		// add event
		payload, err := json.Marshal(map[string]any{
			"role": map[string]any{
				"id":   targetRole.ID,
				"name": targetRole.Name,
			},
		})
		if err != nil {
			return err
		}
		event := &model.Event{
			ActorID:   &actorID,
			UserID:    &targetUserID,
			EventType: "ROLE_REMOVED",
			Payload:   payload,
		}
		return eventRepository.Create(ctx, event)
	})
}

func (s *roleService) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	return s.roleRepository.GetAllRoles(ctx)
}

func (s *roleService) GetAllPermissions(ctx context.Context) ([]model.Permission, error) {
	return s.roleRepository.GetAllPermissions(ctx)
}

func (s *roleService) CreateRole(ctx context.Context, actorID uuid.UUID, role *model.Role) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		roleRepo := repository.NewRoleRepository(tx)
		eventRepo := repository.NewEventRepository(tx)

		if err := CheckPermission(ctx, roleRepo, actorID, permission.RoleCreate); err != nil {
			return err
		}

		actorRole, err := roleRepo.GetHighestUserRole(ctx, actorID)
		if err != nil {
			return err
		}
		if role.Position >= actorRole.Position {
			return repository.ErrForbidden
		}

		if err := roleRepo.Create(ctx, role); err != nil {
			return err
		}

		payload, err := json.Marshal(
			map[string]interface{}{
				"role_id": role.ID,
				"name":    role.Name,
			},
		)
		if err != nil {
			return err
		}

		event := &model.Event{
			ActorID:   &actorID,
			EventType: "ROLE_CREATED",
			Payload:   payload,
		}

		return eventRepo.Create(ctx, event)
	})
}

func (s *roleService) DeleteRole(ctx context.Context, actorID uuid.UUID, roleID int) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		roleRepo := repository.NewRoleRepository(tx)
		eventRepo := repository.NewEventRepository(tx)

		if err := CheckPermission(ctx, roleRepo, actorID, permission.RoleDelete); err != nil {
			return err
		}
		role, err := roleRepo.GetRoleByID(ctx, roleID)
		if err != nil {
			return err
		}
		if role.IsSystem {
			return repository.ErrForbidden
		}
		if err := CanManageRole(ctx, roleRepo, actorID, role); err != nil {
			return err
		}
		if err := roleRepo.SoftDelete(ctx, roleID, actorID); err != nil {
			return err
		}
		payload, _ := json.Marshal(
			map[string]interface{}{
				"role_id": roleID,
				"name":    role.Name,
			},
		)
		return eventRepo.Create(ctx, &model.Event{
			ActorID:   &actorID,
			EventType: "ROLE_DELETED",
			Payload:   payload,
		})
	})
}

func (s *roleService) RestoreRole(ctx context.Context, actorID uuid.UUID, roleID int) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		roleRepo := repository.NewRoleRepository(tx)
		eventRepo := repository.NewEventRepository(tx)

		if err := CheckPermission(ctx, roleRepo, actorID, permission.RoleRestore); err != nil {
			return err
		}
		role, err := roleRepo.GetRoleByIDIncludeDeleted(ctx, roleID)
		if err != nil {
			return err
		}
		if err := CanManageRole(ctx, roleRepo, actorID, role); err != nil {
			return err
		}
		if err := roleRepo.Restore(ctx, roleID); err != nil {
			return err
		}
		payload, _ := json.Marshal(
			map[string]interface{}{
				"role_id": roleID,
				"name":    role.Name,
			},
		)
		return eventRepo.Create(ctx, &model.Event{
			ActorID:   &actorID,
			EventType: "ROLE_RESTORED",
			Payload:   payload,
		})
	})
}
