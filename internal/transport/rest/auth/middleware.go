package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func MockAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &User{
			ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		}

		c.Set(string(UserContextKey), user)
		c.Next()
	}
}
