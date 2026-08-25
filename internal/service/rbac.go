package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/repository"
)

func CanManageRole(ctx context.Context, roleRepo repository.RoleRepository, actorID uuid.UUID, role *model.Role) error {
	actorRole, err := roleRepo.GetHighestUserRole(ctx, actorID)
	if err != nil {
		return err
	}
	if actorRole.Position <= role.Position {
		return ErrForbidden
	}
	return nil
}

func CheckPermission(ctx context.Context, repo repository.RoleRepository, userID uuid.UUID, permission string) error {
	ok, err := repo.HasPermission(ctx, userID, permission)
	if err != nil {
		return err
	}
	if !ok {
		return ErrForbidden
	}
	return nil
}
