package response

import (
	"github.com/gin-gonic/gin"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Envelope struct {
	Data  any    `json:"data,omitempty"`
	Error *Error `json:"error,omitempty"`
}

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, Envelope{
		Data: data,
	})
}

func Fail(c *gin.Context, status int, code string, message string) {
	c.JSON(status, Envelope{
		Error: &Error{
			Code:    code,
			Message: message,
		},
	})
}
