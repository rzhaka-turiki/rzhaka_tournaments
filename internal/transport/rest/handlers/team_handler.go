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

type TeamHandler struct {
	teamService service.TeamService
}

func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

func (h *TeamHandler) Create(c *gin.Context) {
	var req dto.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	team := &model.Team{
		Name:         req.Name,
		ShortName:    *req.ShortName,
		LogoPath:     *req.LogoPath,
		LogoDarkPath: *req.LogoDarkPath,
	}
	err := h.teamService.Create(c.Request.Context(), actorID, team)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, dto.FromTeam(*team))
}

func (h *TeamHandler) GetByID(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_TEAM_ID", "invalid team id")
		return
	}
	team, err := h.teamService.GetByID(c.Request.Context(), teamID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromTeam(*team))
}

func (h *TeamHandler) GetOwnedTeams(c *gin.Context) {
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	includeArchived := c.Query("include_archived") == "true"
	teams, err := h.teamService.GetOwnedTeams(
		c.Request.Context(),
		actorID,
		includeArchived,
	)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromTeams(teams))
}

func (h *TeamHandler) GetMemberTeams(c *gin.Context) {
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	includeArchived := c.Query("include_archived") == "true"
	teams, err := h.teamService.GetMemberTeams(
		c.Request.Context(),
		actorID,
		includeArchived,
	)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.FromTeams(teams))
}

func (h *TeamHandler) Update(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_TEAM_ID", "invalid team id")
		return
	}

	var req dto.UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	teamUPD := service.TeamUpdate{
		Name:         req.Name,
		ShortName:    req.ShortName,
		LogoPath:     req.LogoPath,
		LogoDarkPath: req.LogoDarkPath,
	}
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}

	err = h.teamService.Update(
		c.Request.Context(),
		actorID,
		teamID,
		teamUPD,
	)
	if err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, http.StatusNoContent, nil)
}

func (h *TeamHandler) Archive(c *gin.Context) {
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
	err = h.teamService.Archive(c.Request.Context(), teamID, actorID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

func (h *TeamHandler) Restore(c *gin.Context) {
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
	err = h.teamService.Restore(c.Request.Context(), teamID, actorID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

func (h *TeamHandler) RemoveMember(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_TEAM_ID", "invalid team id")
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
	err = h.teamService.RemoveMember(c.Request.Context(), teamID, actorID, memberID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

func (h *TeamHandler) TransferOwnership(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_TEAM_ID", "invalid team id")
		return
	}
	var req dto.TransferOwnershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}
	actorID := auth.UserID(c)
	if actorID == uuid.Nil {
		HandleError(c, service.ErrUnauthorized)
		return
	}
	err = h.teamService.TransferOwnership(c.Request.Context(), teamID, actorID, req.NewOwnerID)
	if err != nil {
		HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

func (h *TeamHandler) GetMembers(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_TEAM_ID", "invalid team id")
		return
	}

	members, err := h.teamService.GetMembers(
		c.Request.Context(),
		teamID,
	)

	if err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, dto.FromTeamMembers(members))
}
