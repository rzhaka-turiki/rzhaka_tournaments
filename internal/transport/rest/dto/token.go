package dto

type CreateTokenRequest struct {
	StatsToken  string `json:"stats_token" binding:"required"`
	PlayerToken string `json:"player_token" binding:"required"`
	AdminToken  string `json:"admin_token" binding:"required"`
}
