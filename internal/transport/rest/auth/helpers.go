package auth

import (
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func CurrentUser(c *gin.Context) *User {
	value, exists := c.Get(string(UserContextKey))
	if !exists {
		return nil
	}
	user, ok := value.(*User)
	if !ok {
		return nil
	}
	return user
}

func UserID(c *gin.Context) uuid.UUID {
	user := CurrentUser(c)
	if user == nil {
		return uuid.Nil
	}
	return user.ID
}
