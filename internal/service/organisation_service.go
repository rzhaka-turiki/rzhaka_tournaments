package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/database"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/repository"
)

type OrganisationService interface {
	Create(ctx context.Context, ownerID, actorID uuid.UUID, organisation *model.Organisation) error
	GetByID(ctx context.Context, organisationID uuid.UUID) (*model.Organisation, error)
	GetOwnedOrganisations(ctx context.Context, ownerID uuid.UUID) ([]model.Organisation, error)
	GetMemberOrganisations(ctx context.Context, userID uuid.UUID) ([]model.Organisation, error)
	Update(ctx context.Context, actorID, organisationID uuid.UUID, req OrganisationUpdate) error
	GetMembers(ctx context.Context, organisationID uuid.UUID) ([]model.OrganisationMember, error)
	RemoveMember(ctx context.Context, organisationID, actorID, memberID uuid.UUID) error
	TransferOwnership(ctx context.Context, organisationID, actorID, newOwnerID uuid.UUID) error
}

type organisationService struct {
	txManager                     *database.TxManager
	organisationRepository        repository.OrganisationsRepository
	organisationMembersRepository repository.OrganisationMembersRepository
}

type OrganisationUpdate struct {
	Name      *string
	ShortName *string
	ImageURL  *string
	BannerURL *string
}

func NewOrganisationService(
	txManager *database.TxManager,
	organisationRepository repository.OrganisationsRepository,
	organisationsMembersRepository repository.OrganisationMembersRepository,
) OrganisationService {
	return &organisationService{
		txManager:                     txManager,
		organisationRepository:        organisationRepository,
		organisationMembersRepository: organisationsMembersRepository,
	}
}

func (s *organisationService) Create(ctx context.Context, ownerID, actorID uuid.UUID, organisation *model.Organisation) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		organisationRepo := repository.NewOrganisationRepository(tx)
		organisationMemberRepo := repository.NewOrganisationMembersRepository(tx)
		organsationRoleRepo := repository.NewOrganisationRoleRepository(tx)

		organisation.OwnerID = ownerID
		if err := organisationRepo.Create(ctx, organisation); err != nil {
			return err
		}
		role, err := organsationRoleRepo.GetByCode(ctx, "owner")
		if err != nil {
			return err
		}
		return organisationMemberRepo.Create(ctx, &model.OrganisationMember{
			OrganisationID: organisation.ID,
			UserID:         ownerID,
			RoleID:         role.ID,
		})
	})
}

func (s *organisationService) GetByID(ctx context.Context, organisationID uuid.UUID) (*model.Organisation, error) {
	return s.organisationRepository.GetByID(ctx, organisationID)
}

func (s *organisationService) GetOwnedOrganisations(ctx context.Context, ownerID uuid.UUID) ([]model.Organisation, error) {
	return s.organisationRepository.GetByOwner(ctx, ownerID)
}

func (s *organisationService) GetMemberOrganisations(ctx context.Context, userID uuid.UUID) ([]model.Organisation, error) {
	return s.organisationRepository.GetByOwner(ctx, userID)
}

func (s *organisationService) Update(ctx context.Context, actorID, organisationID uuid.UUID, req OrganisationUpdate) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		organisationRepo := repository.NewOrganisationRepository(tx)
		orgMemberRepo := repository.NewOrganisationMembersRepository(tx)
		roleCode, err := orgMemberRepo.GetRole(ctx, organisationID, actorID)
		if err != nil {
			return err
		}
		if *roleCode != "owner" {
			return ErrForbidden
		}
		organisation, err := organisationRepo.GetByID(ctx, organisationID)
		if err != nil {
			return err
		}
		if req.Name != nil {
			organisation.Name = *req.Name
		}
		if req.ShortName != nil {
			organisation.ShortName = *req.ShortName
		}
		if req.ImageURL != nil {
			organisation.ImageURL = *req.ImageURL
		}
		if req.BannerURL != nil {
			organisation.BannerURL = *req.BannerURL
		}
		return organisationRepo.Update(ctx, organisation)
	})
}

func (s *organisationService) GetMembers(ctx context.Context, organisationID uuid.UUID) ([]model.OrganisationMember, error) {
	return s.organisationMembersRepository.GetByOrganisation(ctx, organisationID)
}

func (s *organisationService) TransferOwnership(ctx context.Context, organisationID, actorID, newOwnerID uuid.UUID) error {
	return s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		organisationRepo := repository.NewOrganisationRepository(tx)
		organisationMemberRepo := repository.NewOrganisationMembersRepository(tx)
		orgRoleRepo := repository.NewOrganisationRoleRepository(tx)
		roleCode, err := organisationMemberRepo.GetRole(ctx, organisationID, actorID)
		if err != nil {
			return err
		}
		if *roleCode != "owner" {
			return ErrForbidden
		}
		isMember, err := organisationMemberRepo.Exists(ctx, newOwnerID, organisationID)
		if err != nil {
			return err
		}
		if !isMember {
			return repository.ErrNotOrgMember
		}
		if err := organisationRepo.UpdateOwner(ctx, organisationID, newOwnerID); err != nil {
			return err
		}
		roleOwner, err := orgRoleRepo.GetByCode(ctx, "owner")
		if err != nil {
			return err
		}
		roleManager, err := orgRoleRepo.GetByCode(ctx, "manager")
		if err != nil {
			return err
		}
		if err := organisationMemberRepo.UpdateRole(ctx, organisationID, actorID, roleManager.ID); err != nil {
			return err
		}
		if err := organisationMemberRepo.UpdateRole(ctx, organisationID, newOwnerID, roleOwner.ID); err != nil {
			return err
		}
		return nil
	})
}

func (s *organisationService) RemoveMember(ctx context.Context, organisationID, actorID, memberID uuid.UUID) error {
	roleCode, err := s.organisationMembersRepository.GetRole(ctx, organisationID, actorID)
	if err != nil {
		return err
	}
	if *roleCode != "owner" {
		return ErrForbidden
	}
	roleCode, err = s.organisationMembersRepository.GetRole(ctx, organisationID, memberID)
	if err != nil {
		return err
	}
	if *roleCode == "owner" {
		return ErrCannotRemoveOwner
	}
	return s.organisationMembersRepository.Remove(ctx, organisationID, memberID)
}
