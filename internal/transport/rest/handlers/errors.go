package handlers

import (
	"errors"
	"net/http"

	"github.com/a1uka/rzhaka_tournaments/internal/repository"
	"github.com/a1uka/rzhaka_tournaments/internal/service"
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/response"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		response.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())

	case errors.Is(err, service.ErrForbidden):
		response.Fail(c, http.StatusForbidden, "FORBIDDEN", err.Error())

	case errors.Is(err, repository.ErrConflict):
		response.Fail(c, http.StatusConflict, "CONFLICT", err.Error())

	case errors.Is(err, repository.ErrInvalid):
		response.Fail(c, http.StatusBadRequest, "INVALID_DATA", err.Error())

	case errors.Is(err, service.ErrUnauthorized):
		response.Fail(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
	default:
		//log.Println("handler error:", err)
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
}
