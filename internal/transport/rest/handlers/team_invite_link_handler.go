package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/service"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/auth"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/dto"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/response"
)

type TeamInviteLinkHandler struct {
	inviteLinkService service.TeamInviteLinkService
}

func NewTeamInviteLinkHandler(
	inviteLinkService service.TeamInviteLinkService,
) *TeamInviteLinkHandler {
	return &TeamInviteLinkHandler{
		inviteLinkService: inviteLinkService,
	}
}

func (h *TeamInviteLinkHandler) Create(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(
			c,
			http.StatusBadRequest,
			"INVALID_TEAM_ID",
			"invalid team id",
		)
		return
	}

	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}

	var req dto.CreateInviteLinkRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(
			c,
			http.StatusBadRequest,
			"INVALID_BODY",
			"invalid request body",
		)
		return
	}

	link, err := h.inviteLinkService.CreateLink(
		c.Request.Context(),
		actorID,
		teamID,
		req,
	)

	if err != nil {
		HandleError(c, err)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		dto.FromInviteLink(*link),
	)
}

func (h *TeamInviteLinkHandler) PreviewLink(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		response.Fail(c, http.StatusBadRequest, "INVALID_TOKEN", "token is required")
		return
	}
	preview, err := h.inviteLinkService.PreviewLink(c.Request.Context(), token)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromInviteLinkPreview(*preview))
}

func (h *TeamInviteLinkHandler) AcceptByLink(c *gin.Context) {
	token := c.Param("token")

	if token == "" {
		response.Fail(c, http.StatusBadRequest, "INVALID_TOKEN", "token is required")
		return
	}
	userID := auth.UserID(c)
	if userID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	err := h.inviteLinkService.AcceptByLink(c.Request.Context(), token, userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

func (h *TeamInviteLinkHandler) DeleteLink(c *gin.Context) {
	linkID, err := uuid.Parse(c.Param("linkID"))

	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_LINK_ID", "invalid invite link id")
		return
	}
	actorID := auth.UserID(c)

	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	err = h.inviteLinkService.DeleteLink(c.Request.Context(), actorID, linkID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

func (h *TeamInviteLinkHandler) GetTeamLinks(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_TEAM_ID", "invalid team id")
		return
	}

	actorID := auth.UserID(c)

	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}

	links, err := h.inviteLinkService.GetTeamLinks(c.Request.Context(), actorID, teamID)

	if err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, dto.FromInviteLinks(links))
}
