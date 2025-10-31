package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				e, ok := err.(error)
				if !ok {
					e = errors.New(fmt.Sprint(err))
				}
				statusCode := c.Writer.Status()
				if statusCode < http.StatusBadRequest {
					statusCode = http.StatusInternalServerError
				}
				c.String(statusCode, e.Error())
				c.Abort()
				return
			}
		}()
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}
		statusCode := c.Writer.Status()
		if statusCode < http.StatusBadRequest {
			statusCode = http.StatusInternalServerError
		}
		c.String(statusCode, err.Error())
		c.Abort()
	}
}
