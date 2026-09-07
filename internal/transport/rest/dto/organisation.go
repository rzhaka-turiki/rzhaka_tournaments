package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type CreateOrgRequest struct {
	Name      string    `json:"name" binding:"required"`
	OwnerID   uuid.UUID `json:"ownerID" binding:"requiered"`
	ShortName *string   `json:"short_name"`
	ImageURL  *string   `json:"image_url"`
	BannerURL *string   `json:"banner_url"`
}

type OrgResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name" binding:"required"`
	ShortName *string   `json:"short_name"`
	ImageURL  *string   `json:"image_url"`
	BannerURL *string   `json:"banner_url"`
}

type UpdateOrgRequest struct {
	Name      *string `json:"name" binding:"required"`
	ShortName *string `json:"short_name"`
	ImageURL  *string `json:"image_url"`
	BannerURL *string `json:"banner_url"`
}

type TransferOrgOwnershipRequest struct {
	NewOwnerID uuid.UUID `json:"new_owner_id" binding:"required"`
}

type OrgMemberResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Role     uuid.UUID `json:"role_id"`
	JoinedAt time.Time `json:"joined_at"`
}

func FromOrgMember(member model.OrganisationMember) OrgMemberResponse {
	return OrgMemberResponse{
		UserID:   member.UserID,
		Role:     member.RoleID,
		JoinedAt: member.CreatedAt,
	}
}

func FromOrgMembers(members []model.OrganisationMember) []OrgMemberResponse {
	result := make([]OrgMemberResponse, 0, len(members))
	for _, member := range members {
		result = append(result, FromOrgMember(member))
	}
	return result
}

func FromOrg(org model.Organisation) OrgResponse {
	return OrgResponse{
		ID:        org.ID,
		Name:      org.Name,
		ImageURL:  &org.ImageURL,
		ShortName: &org.ShortName,
		BannerURL: &org.BannerURL,
	}
}

func FromOrgs(orgs []model.Organisation) []OrgResponse {
	OrgsResponse := make([]OrgResponse, 0, len(orgs))
	for _, org := range orgs {
		OrgsResponse = append(OrgsResponse, FromOrg(org))
	}
	return OrgsResponse
}
