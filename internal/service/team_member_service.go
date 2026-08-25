package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type TeamMemberService interface {
	GetMembers(ctx context.Context, teamID uuid.UUID) ([]model.TeamMember, error)
	RemoveMember(ctx context.Context, teamID, actorID, memberID uuid.UUID) error
	TransferOwnership(ctx context.Context, teamID, ownerID, newOwnerID uuid.UUID) error
}
