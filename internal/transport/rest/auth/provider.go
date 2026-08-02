package auth

import (
	"github.com/gin-gonic/gin"
)

func Middleware() gin.HandlerFunc {
	return MockAuth()
}
