package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"gocourse/pkg/httpx"
)

func JSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}
		if c.Request.Header.Get(httpx.HeaderContentType) != httpx.ContentTypeJSON {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid content type, must be json"))
			return
		}
		c.Next()
	}
}
