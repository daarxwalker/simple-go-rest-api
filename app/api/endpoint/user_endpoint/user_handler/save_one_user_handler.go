package user_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gocourse/app/api/endpoint/user_endpoint/user_request"
	"gocourse/app/api/endpoint/user_endpoint/user_response"
	"gocourse/common/database"
	"gocourse/common/logger"
	"gocourse/internal/application/user_usecase"
	"gocourse/internal/domain/user_domain"
)

type SaveOneDeps struct {
	DB     database.DB
	Logger *logger.Logger
}

func SaveOne(deps SaveOneDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body user_request.SaveOne
		if bindErr := c.ShouldBindJSON(&body); bindErr != nil {
			c.Status(http.StatusBadRequest)
			c.Error(bindErr)
			return
		}
		userId, saveOneUserErr := user_usecase.SaveOne(
			c,
			user_usecase.SaveOneDeps{
				DB:     deps.DB,
				Logger: deps.Logger,
			},
			user_domain.UserEntity{
				Id:    body.Id,
				Name:  body.Name,
				Email: body.Email,
			},
		)
		if saveOneUserErr != nil {
			c.Status(http.StatusBadRequest)
			c.Error(saveOneUserErr)
			return
		}
		c.JSON(http.StatusCreated, user_response.SaveOne{Id: userId})
	}
}
