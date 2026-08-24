package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/response"
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
