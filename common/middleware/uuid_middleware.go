package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func UUID(paramKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := uuid.Parse(c.Param(paramKey)); err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid uuid format"))
			return
		}
		c.Next()
	}
}
