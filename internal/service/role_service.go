package service

import (
	"context"
	"encoding/json"

	"github.com/a1uka/rzhaka_tournaments/internal/database"
	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/a1uka/rzhaka_tournaments/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RoleService interface {
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error)
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]model.Permission, error)
	AssignRole(ctx context.Context, actorID uuid.UUID, targetUserID uuid.UUID, roleID int) error
	RemoveRole(ctx context.Context, actorID uuid.UUID, targetUserID uuid.UUID, roleID int) error
}

type roleService struct {
	txManager       *database.TxManager
	roleRepository  repository.RoleRepository
	eventRepository repository.EventRepository
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

func (s *roleService) hasPermission(
	ctx context.Context,
	roleRepository repository.RoleRepository,
	userID uuid.UUID,
	permission string,
) error {
	hasPermission, err := roleRepository.HasPermission(ctx, userID, permission)
	if err != nil {
		return err
	}
	if !hasPermission {
		return repository.ErrForbidden
	}
	return nil
}

func (s *roleService) canManageRole(ctx context.Context, roleRepository repository.RoleRepository, actorID uuid.UUID, targetRole *model.Role) error {
	actorRole, err := roleRepository.GetHighestUserRole(ctx, actorID)
	if err != nil {
		return err
	}
	if actorRole.Position <= targetRole.Position {
		return repository.ErrForbidden
	}
	return nil
}

func (s *roleService) canAssignRole(ctx context.Context, roleRepository repository.RoleRepository, actorID uuid.UUID, role *model.Role) error {
	if err := s.hasPermission(ctx, roleRepository, actorID, "role.assign"); err != nil {
		return err
	}
	if err := s.canManageRole(ctx, roleRepository, actorID, role); err != nil {
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
