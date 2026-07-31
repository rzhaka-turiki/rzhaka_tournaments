package routes

import (
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	h *rest.Handlers,
) {
	registerHealth(router, h.Health)
	v1 := router.Group("/api/v1")
	{
		registerUsers(v1, h.User)
		registerRoles(v1, h.Role)
	}
}
