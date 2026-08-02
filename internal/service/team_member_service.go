package service

import (
	"context"

	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/google/uuid"
)

type TeamMemberService interface {
	GetMembers(ctx context.Context, teamID uuid.UUID) ([]model.TeamMember, error)
	RemoveMember(ctx context.Context, teamID, actorID, memberID uuid.UUID) error
	TransferOwnership(ctx context.Context, teamID, ownerID, newOwnerID uuid.UUID) error
}
