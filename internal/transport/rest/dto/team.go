package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type CreateTeamRequest struct {
	Name         string  `json:"name" binding:"required"`
	ShortName    *string `json:"short_name"`
	LogoPath     *string `json:"logo_path"`
	LogoDarkPath *string `json:"logo_dark_path"`
}

type TeamMemberResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type TeamResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name" binding:"required"`
	ShortName    *string   `json:"short_name"`
	LogoPath     *string   `json:"logo_path"`
	LogoDarkPath *string   `json:"logo_dark_path"`
}

type UpdateTeamRequest struct {
	Name         *string `json:"name" binding:"required"`
	ShortName    *string `json:"short_name"`
	LogoPath     *string `json:"logo_path"`
	LogoDarkPath *string `json:"logo_dark_path"`
}

type TransferOwnershipRequest struct {
	NewOwnerID uuid.UUID `json:"new_owner_id" binding:"required"`
}

func FromTeam(team model.Team) TeamResponse {
	return TeamResponse{
		ID:           team.ID,
		Name:         team.Name,
		ShortName:    &team.ShortName,
		LogoPath:     &team.LogoPath,
		LogoDarkPath: &team.LogoDarkPath,
	}
}

func FromTeams(teams []model.Team) []TeamResponse {
	TeamsResponse := make([]TeamResponse, 0, len(teams))
	for _, team := range teams {
		TeamsResponse = append(TeamsResponse, FromTeam(team))
	}
	return TeamsResponse
}

func FromTeamMember(member model.TeamMember) TeamMemberResponse {
	return TeamMemberResponse{
		UserID:   member.UserID,
		Role:     member.Role,
		JoinedAt: member.JoinedAt,
	}
}

func FromTeamMembers(members []model.TeamMember) []TeamMemberResponse {
	result := make([]TeamMemberResponse, 0, len(members))

	for _, member := range members {
		result = append(result, FromTeamMember(member))
	}

	return result
}
