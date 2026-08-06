package handlers

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error string `json:"error"`
}

func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, ErrorResponse{
		Error: message,
	})
}
