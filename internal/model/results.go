package model

import (
	"time"

	"github.com/google/uuid"
)

type RawPlayerResult struct {
	PlayerHash string `json:"player_hash"`
	Legend     string `json:"legends"`
	Kills      int    `json:"kills"`
	Assists    int    `json:"assists"`
	Damage     int    `json:"damage"`
	Headshots  int    `json:"headshots"`
	Revives    int    `json:"revives"`
	Respawns   int    `json:"respawns"`
	TimeAlive  int    `json:"time_alive"`
	Hits       int    `json:"hits"`
	Shots      int    `json:"shots"`
	Knockdowns int    `json:"knockdowns"`
}

type TeamMatchResult struct {
	TeamID    uuid.UUID           `json:"team_id"`
	TeamName  string              `json:"team_name"`
	Placement int                 `json:"placement"`
	Points    int                 `json:"points"`
	Kills     int                 `json:"kills"`
	Players   []PlayerMatchResult `json:"players"`
}

type PlayerMatchResult struct {
	PlayerHash string `json:"player_hash"`
	Gamertag   string `json:"gamertag"`
	Legend     string `json:"legend"`
	Kills      int    `json:"kills"`
	Assists    int    `json:"assists"`
	Damage     int    `json:"damage"`
	// остальные поля при желании
}

type MatchResult struct {
	MatchID     uuid.UUID         `json:"match_id"`
	GameMatchID string            `json:"game_match_id,omitempty"`
	Map         string            `json:"map"`
	StartedAt   time.Time         `json:"started_at"`
	Teams       []TeamMatchResult `json:"teams"`
}

type StageResult struct {
	StageID uuid.UUID         `json:"stage_id"`
	Name    string            `json:"name"`
	Teams   []TeamStageResult `json:"teams"`
}

type TeamStageResult struct {
	TeamID       uuid.UUID         `json:"team_id"`
	TeamName     string            `json:"team_name"`
	TotalPoints  int               `json:"total_points"`
	TotalKills   int               `json:"total_kills"`
	MatchResults []TeamMatchResult `json:"match_results,omitempty"`
}
