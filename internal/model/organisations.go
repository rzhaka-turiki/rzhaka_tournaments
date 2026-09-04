package model

import (
	"time"

	"github.com/google/uuid"
)

type Organisation struct {
	ID        uuid.UUID
	Name      string
	ShortName string
	ImageURL  string
	BannerURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrganisationRole struct {
	ID        uuid.UUID
	Name      string
	RoleColor string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrganisationMember struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	RoleID         uuid.UUID
	OrganisationID uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
