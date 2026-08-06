package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		requestID, _ := c.Get("request_id")

		c.Next()
		duration := time.Since(start)
		status := c.Writer.Status()

		log.Printf(
			"request_id=%s method=%s path=%s status=%d duration=%v",
			requestID,
			method,
			path,
			status,
			duration,
		)
	}
}
