package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/database"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/repository"
)

type TeamService interface {
	Create(ctx context.Context, ownerID uuid.UUID, team *model.Team) error
	GetByID(ctx context.Context, teamID uuid.UUID) (*model.Team, error)
	GetOwnedTeams(ctx context.Context, userID uuid.UUID, includeArchived bool) ([]model.Team, error)
	GetMemberTeams(ctx context.Context, userID uuid.UUID, includeArchived bool) ([]model.Team, error)
	Update(ctx context.Context, actorID, teamID uuid.UUID, req TeamUpdate) error
	Archive(ctx context.Context, teamID, actorID uuid.UUID) error
	Restore(ctx context.Context, teamID, actorID uuid.UUID) error
	GetMembers(ctx context.Context, teamID uuid.UUID) ([]model.TeamMember, error)
	RemoveMember(ctx context.Context, teamID, actorID, memberID uuid.UUID) error
	TransferOwnership(ctx context.Context, teamID, actorID, newOwnerID uuid.UUID) error
}

type TeamUpdate struct {
	Name         *string
	ShortName    *string
	LogoPath     *string
	LogoDarkPath *string
}

type teamService struct {
	txManager              *database.TxManager
	teamRepository         repository.TeamRepository
	teamMemberRepository   repository.TeamMemberRepository
	teamSnapshotRepository repository.TeamSnapshotRepository
}

func NewTeamService(
	txManager *database.TxManager,
	teamRepository repository.TeamRepository,
	teamMemberRepository repository.TeamMemberRepository,
	teamSnapshotRepository repository.TeamSnapshotRepository,
) TeamService {
	return &teamService{
		txManager:              txManager,
		teamRepository:         teamRepository,
		teamMemberRepository:   teamMemberRepository,
		teamSnapshotRepository: teamSnapshotRepository,
	}
}

func (s *teamService) Create(ctx context.Context, ownerID uuid.UUID, team *model.Team) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		teamRepo := repository.NewTeamRepository(tx)
		memberRepo := repository.NewTeamMemberRepository(tx)
		snapshotRepo := repository.NewTeamSnapshotRepository(tx)

		count, err := teamRepo.CountActiveTeams(ctx, ownerID)
		if err != nil {
			return err
		}
		if count >= 5 {
			return ErrTeamLimitReached
		}
		team.OwnerID = ownerID
		if err := teamRepo.Create(ctx, team); err != nil {
			return err
		}
		if err := memberRepo.Create(ctx, &model.TeamMember{
			TeamID: team.ID,
			UserID: ownerID,
			Role:   model.TeamRoleOwner,
		}); err != nil {
			return err
		}
		return snapshotRepo.Create(ctx, &model.TeamSnapshot{
			TeamID:       team.ID,
			Name:         team.Name,
			ShortName:    team.ShortName,
			LogoPath:     team.LogoPath,
			LogoDarkPath: team.LogoDarkPath,
		})
	})
}

func (s *teamService) GetByID(ctx context.Context, teamID uuid.UUID) (*model.Team, error) {
	return s.teamRepository.GetByID(ctx, teamID)
}

func (s *teamService) GetOwnedTeams(ctx context.Context, userID uuid.UUID, includeArchived bool) ([]model.Team, error) {
	return s.teamRepository.GetOwnedTeams(ctx, userID, includeArchived)
}

func (s *teamService) GetMemberTeams(ctx context.Context, userID uuid.UUID, includeArchived bool) ([]model.Team, error) {
	return s.teamRepository.GetMemberTeams(ctx, userID, includeArchived)
}

func (s *teamService) Update(ctx context.Context, actorID, teamID uuid.UUID, req TeamUpdate) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		teamRepo := repository.NewTeamRepository(tx)
		snapshotRepo := repository.NewTeamSnapshotRepository(tx)
		isOwner, err := teamRepo.IsOwner(ctx, teamID, actorID)
		if err != nil {
			return err
		}
		if !isOwner {
			return ErrForbidden
		}
		team, err := teamRepo.GetByID(ctx, teamID)
		if err != nil {
			return err
		}
		if req.Name != nil {
			team.Name = *req.Name
		}
		if req.ShortName != nil {
			team.ShortName = *req.ShortName
		}

		if req.LogoPath != nil {
			team.LogoPath = *req.LogoPath
		}

		if req.LogoDarkPath != nil {
			team.LogoDarkPath = *req.LogoDarkPath
		}
		err = teamRepo.Update(ctx, team)
		if err != nil {
			return err
		}
		return snapshotRepo.Create(ctx, &model.TeamSnapshot{
			TeamID:       team.ID,
			Name:         team.Name,
			ShortName:    team.ShortName,
			LogoPath:     team.LogoPath,
			LogoDarkPath: team.LogoDarkPath,
		})
	})
}

func (s *teamService) Archive(ctx context.Context, teamID, actorID uuid.UUID) error {
	isOwner, err := s.teamRepository.IsOwner(
		ctx,
		teamID,
		actorID,
	)

	if err != nil {
		return err
	}

	if !isOwner {
		return ErrForbidden
	}

	return s.teamRepository.Archive(
		ctx,
		teamID,
	)
}

func (s *teamService) Restore(ctx context.Context, teamID, actorID uuid.UUID) error {
	isOwner, err := s.teamRepository.IsOwner(
		ctx,
		teamID,
		actorID,
	)

	if err != nil {
		return err
	}

	if !isOwner {
		return ErrForbidden
	}

	count, err := s.teamRepository.CountActiveTeams(ctx, actorID)
	if err != nil {
		return err

	}
	if count >= 5 {
		return ErrTeamLimitReached
	}

	return s.teamRepository.Restore(
		ctx,
		teamID,
	)
}

func (s *teamService) GetMembers(ctx context.Context, teamID uuid.UUID) ([]model.TeamMember, error) {
	return s.teamMemberRepository.GetMembers(ctx, teamID)
}

func (s *teamService) RemoveMember(ctx context.Context, teamID, actorID, memberID uuid.UUID) error {
	isOwner, err := s.teamRepository.IsOwner(ctx, teamID, actorID)
	if err != nil {
		return err
	}
	if !isOwner {
		return ErrForbidden
	}
	isOwner, err = s.teamRepository.IsOwner(ctx, teamID, memberID)
	if err != nil {
		return err
	}
	if isOwner {
		return ErrCannotRemoveOwner
	}
	return s.teamMemberRepository.Remove(ctx, teamID, memberID)
}

func (s *teamService) TransferOwnership(ctx context.Context, teamID, actorID, newOwnerID uuid.UUID) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		teamRepo := repository.NewTeamRepository(tx)
		memberRepo := repository.NewTeamMemberRepository(tx)
		isOwner, err := teamRepo.IsOwner(ctx, teamID, actorID)
		if err != nil {
			return err
		}
		if !isOwner {
			return ErrForbidden
		}
		isMember, err := memberRepo.Exists(ctx, teamID, newOwnerID)
		if err != nil {
			return err
		}
		if !isMember {
			return repository.ErrNotTeamMember
		}
		if err := teamRepo.UpdateOwner(ctx, teamID, newOwnerID); err != nil {
			return err
		}
		if err := memberRepo.UpdateRole(ctx, teamID, actorID, model.TeamRoleManager); err != nil {
			return err
		}
		if err := memberRepo.UpdateRole(ctx, teamID, newOwnerID, model.TeamRoleOwner); err != nil {
			return err
		}
		return nil
	})
}
