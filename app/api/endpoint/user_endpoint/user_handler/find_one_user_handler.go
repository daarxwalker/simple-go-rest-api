package user_handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"gocourse/app/api/endpoint/user_endpoint/user_response"
	"gocourse/common/database"
	"gocourse/common/logger"
	"gocourse/common/storage"
	"gocourse/internal/application/user_usecase"
)

type FindOneDeps struct {
	DB     database.DB
	Cache  *storage.Storage
	Logger *logger.Logger
}

func FindOneUser(deps FindOneDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userResponse user_response.FindOne
		userId := c.Param("id")
		cacheKey := "user-" + userId
		if getUserCacheErr := deps.Cache.Get(c, cacheKey, &userResponse); getUserCacheErr != nil {
			deps.Logger.Error(fmt.Errorf("get user cache failed: %w", getUserCacheErr))
		}
		if len(userResponse.Id) > 0 {
			c.JSON(http.StatusOK, userResponse)
			return
		}
		user, findOneUserErr := user_usecase.FindOneUser(
			c,
			user_usecase.FindOneUserDeps{
				DB:     deps.DB,
				Logger: deps.Logger,
			},
			userId,
		)
		if findOneUserErr != nil {
			c.Status(http.StatusInternalServerError)
			c.Error(findOneUserErr)
			return
		}
		if len(user.Id) == 0 {
			c.Status(http.StatusNotFound)
			c.Error(errors.New("user not found"))
			return
		}
		userResponse = user_response.Map(user)
		if setUserCacheErr := deps.Cache.Set(c, cacheKey, userResponse, time.Hour); setUserCacheErr != nil {
			deps.Logger.Error(fmt.Errorf("set user cache failed: %w", setUserCacheErr))
		}
		c.JSON(http.StatusOK, userResponse)
	}
}
