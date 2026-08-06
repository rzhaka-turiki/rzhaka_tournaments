package middleware

import (
	"log"
	"net/http"

	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/response"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
