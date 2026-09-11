package app

import "github.com/rzhaka-turiki/rzhaka_tournaments/internal/repository"

type Repositories struct {
	User                   repository.UserRepository
	Role                   repository.RoleRepository
	Event                  repository.EventRepository
	Team                   repository.TeamRepository
	TeamMemberRepository   repository.TeamMemberRepository
	TeamSnapshotRepository repository.TeamSnapshotRepository
	TokenRepository        repository.TokenRepository
	Permission             repository.PermissionRepository
	ApexAccount            repository.ApexAccountRepository
	Organisation           repository.OrganisationsRepository
	Match                  repository.MatchRepository
}
