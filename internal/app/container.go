package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/client/apexverifier"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/client/matchapi"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/config"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/database"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/repository"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/service"
)

type Container struct {
	// repos
	Repositories Repositories

	//services
	Services Services
	// transaction man
	TxManager *database.TxManager
	Clients   Clients
}

func NewContainer(db *pgxpool.Pool, cfg *config.Config) (*Container, error) {
	//tx man
	txManager := database.NewTxManager(db)
	// repos
	userRepository := repository.NewUserRepository(db)
	roleRepository := repository.NewRoleRepository(db)
	eventRepository := repository.NewEventRepository(db)
	permissionRepository := repository.NewPermissionRepository(db)
	rolePermissionRepository := repository.NewRolePermissionRepository(db)
	teamRepository := repository.NewTeamRepository(db)
	teamMemberRepository := repository.NewTeamMemberRepository(db)
	tokenRepository := repository.NewTokensRepository(db)
	teamSnapshotRepository := repository.NewTeamSnapshotRepository(db)
	apexAccountRepository := repository.NewApexAccountRepository(db)
	organisationRepository := repository.NewOrganisationRepository(db)
	organisationMemberRepository := repository.NewOrganisationMembersRepository(db)
	matchRepository := repository.NewMatchRepository(db)
	// cleints
	apexVerifierClient, err := apexverifier.NewClient(cfg.ApexVerifier.GRPCAddr)
	matchAPI, err := matchapi.NewClient(cfg.MatchAPI.GRPCAddr, cfg.MatchAPI.APIKey)
	// services
	teamService := service.NewTeamService(txManager, teamRepository, teamMemberRepository, teamSnapshotRepository)
	userService := service.NewUserService(userRepository, roleRepository)
	organisationService := service.NewOrganisationService(
		txManager,
		organisationRepository,
		organisationMemberRepository,
	)
	tokenService := service.NewTokenService(txManager, matchAPI, tokenRepository, organisationRepository, organisationMemberRepository)
	roleService := service.NewRoleService(txManager, roleRepository, eventRepository)
	permissionService := service.NewPermissionService(permissionRepository)
	rolePermissionService := service.NewRolePermissionService(
		txManager,
		roleRepository,
		permissionRepository,
		rolePermissionRepository,
		eventRepository,
	)
	apexAccountService := service.NewApexAccountService(apexAccountRepository, apexVerifierClient)
	matchService := service.NewMatchService(txManager, matchRepository)
	if err != nil {
		return nil, err
	}
	return &Container{
		Repositories: Repositories{
			User:                   userRepository,
			Role:                   roleRepository,
			Event:                  eventRepository,
			Team:                   teamRepository,
			TeamMemberRepository:   teamMemberRepository,
			TeamSnapshotRepository: teamSnapshotRepository,
			TokenRepository:        tokenRepository,
			Permission:             permissionRepository,
			Organisation:           organisationRepository,
			ApexAccount:            apexAccountRepository,
			Match:                  matchRepository,
		},
		Services: Services{
			User:           userService,
			Role:           roleService,
			Permission:     permissionService,
			Team:           teamService,
			Token:          tokenService,
			Organisation:   organisationService,
			RolePermission: rolePermissionService,
			ApexAccount:    apexAccountService,
			Match:          matchService,
		},
		Clients: Clients{
			ApexVerifier: apexVerifierClient,
			MatchAPI:     matchAPI,
		},
		TxManager: txManager,
	}, nil
}
