package model

import (
	"github.com/google/uuid"
)

type ApexAccount struct {
	UserID     uuid.UUID
	Platform   string
	PlayerHash string
	PlayerUID  string
}
