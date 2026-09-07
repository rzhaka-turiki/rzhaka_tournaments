package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/service"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/auth"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/dto"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/response"
)

type OrganisationHandler struct {
	orgService service.OrganisationService
}

func NewOrgHandler(orgService service.OrganisationService) *OrganisationHandler {
	return &OrganisationHandler{
		orgService: orgService,
	}
}

func (h *OrganisationHandler) Create(c *gin.Context) {
	var req dto.CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	organisation := &model.Organisation{
		Name:      req.Name,
		OwnerID:   req.OwnerID,
		ShortName: *req.ShortName,
		ImageURL:  *req.ImageURL,
		BannerURL: *req.BannerURL,
	}
	err := h.orgService.Create(c.Request.Context(), organisation.OwnerID, actorID, organisation)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromOrg(*organisation))
}

func (h *OrganisationHandler) GetByID(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ORGANISATION_ID", "invalid organisation id")
		return
	}
	org, err := h.orgService.GetByID(c.Request.Context(), orgID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromOrg(*org))
}

func (h *OrganisationHandler) GetOwnedOrgs(c *gin.Context) {
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	orgs, err := h.orgService.GetOwnedOrganisations(
		c.Request.Context(),
		actorID,
	)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromOrgs(orgs))
}

func (h *OrganisationHandler) GetMemberOrgs(c *gin.Context) {
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	orgs, err := h.orgService.GetMemberOrganisations(
		c.Request.Context(),
		actorID,
	)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromOrgs(orgs))
}

func (h *OrganisationHandler) Update(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ORGANISATION_ID", "invalid organisation id")
		return
	}
	var req dto.UpdateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	orgUPD := service.OrganisationUpdate{
		Name:      req.Name,
		ShortName: req.ShortName,
		ImageURL:  req.ImageURL,
		BannerURL: req.BannerURL,
	}
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	err = h.orgService.Update(
		c.Request.Context(),
		actorID,
		orgID,
		orgUPD,
	)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

func (h *OrganisationHandler) RemoveMember(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ORGANISATION_ID", "invalid organisation id")
		return
	}
	memberID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_USER_ID", "invalid user id")
		return
	}
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	err = h.orgService.RemoveMember(c.Request.Context(), orgID, actorID, memberID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

func (h *OrganisationHandler) TransferOwnership(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ORGANISATION_ID", "invalid organisation id")
		return
	}
	var req dto.TransferOrgOwnershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid reqoest body")
		return
	}
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	err = h.orgService.TransferOwnership(c.Request.Context(), orgID, actorID, req.NewOwnerID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

func (h *OrganisationHandler) GetMembers(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_ORGANISATION_ID", "invalid organisation id")
		return
	}
	members, err := h.orgService.GetMembers(c.Request.Context(), orgID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromOrgMembers(members))
}
