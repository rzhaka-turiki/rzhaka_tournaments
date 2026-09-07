package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"
)

func registerOrganisations(rg *gin.RouterGroup, h *handlers.OrganisationHandler) {
	orgs := rg.Group("/organisations")
	{
		orgs.POST("", h.Create)
		orgs.GET("/:id", h.GetByID)

		orgs.PATCH("/:id", h.Update)

		orgs.GET("/:id/members", h.GetMembers)

		orgs.DELETE("/:id/members/:userID", h.RemoveMember)

		orgs.POST("/:id/transfer-ownership", h.TransferOwnership)
	}
	me := rg.Group("/users/me")
	{
		me.GET("/organisations/owned", h.GetOwnedOrgs)
		me.GET("/organisations/member", h.GetMemberOrgs)
	}
}
