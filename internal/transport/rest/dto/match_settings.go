package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type CreateMatchSettingsRequest struct {
	MapID            *uuid.UUID `json:"map_id"`
	PlaylistName     *string    `json:"playlist_name"`
	AdminChat        *bool      `json:"admin_chat"`
	TeamRename       *bool      `json:"team_rename"`
	SelfAssign       *bool      `json:"self_assign"`
	AimAssist        *bool      `json:"aim_asist"`
	AnonMode         *bool      `json:"anon_mode"`
	DropSpotsEnabled *bool      `json:"drop_spots_enabled"`
	FillBotsMode     *bool      `json:"fill_bots_mode"`
}

type MatchSettingsResponse struct {
	MatchID          uuid.UUID `json:"match_id" binding:"required"`
	MapID            uuid.UUID `json:"map_id"`
	PlaylistName     string    `json:"playlist_name"`
	AdminChat        bool      `json:"admin_chat"`
	TeamRename       bool      `json:"team_rename"`
	SelfAssign       bool      `json:"self_assign"`
	AimAssist        bool      `json:"aim_assist"`
	AnonMode         bool      `json:"anon_mode"`
	DropSpotsEnabled bool      `json:"drop_spots_enabled"`
	FillBotsMode     bool      `json:"fill_bots_mode"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type UpdateMatchSettingsRequest struct {
	MapID            *uuid.UUID `json:"map_id"`
	PlaylistName     *string    `json:"playlist_name"`
	AdminChat        *bool      `json:"admin_chat"`
	TeamRename       *bool      `json:"team_rename"`
	SelfAssign       *bool      `json:"self_assign"`
	AimAssist        *bool      `json:"aim_assist"`
	AnonMode         *bool      `json:"anon_mode"`
	DropSpotsEnabled *bool      `json:"drop_spots_enabled"`
	FillBotsMode     *bool      `json:"fill_bots_mode"`
}

func FromMatchSettings(matchSettings *model.MatchSettings) MatchSettingsResponse {
	return MatchSettingsResponse{
		MatchID:          matchSettings.MatchID,
		MapID:            matchSettings.MapID,
		PlaylistName:     matchSettings.PlaylistName,
		AdminChat:        matchSettings.AdminChat,
		TeamRename:       matchSettings.TeamRename,
		SelfAssign:       matchSettings.SelfAssign,
		AimAssist:        matchSettings.AimAssist,
		AnonMode:         matchSettings.AnonMode,
		DropSpotsEnabled: matchSettings.DropSpotsEnabled,
		FillBotsMode:     matchSettings.FillBotsMode,
		CreatedAt:        matchSettings.CreatedAt,
		UpdatedAt:        matchSettings.UpdatedAt,
	}
}

func FromMatchesSettings(matchesSettings []model.MatchSettings) []MatchSettingsResponse {
	MatchesSettingsResponse := make([]MatchSettingsResponse, 0, len(matchesSettings))
	for _, matchSettings := range matchesSettings {
		MatchesSettingsResponse = append(MatchesSettingsResponse, FromMatchSettings(&matchSettings))
	}
	return MatchesSettingsResponse
}
