package user_endpoint

import (
	"github.com/gin-gonic/gin"

	"gocourse/app/api/endpoint/user_endpoint/user_handler"
	"gocourse/common/database"
	"gocourse/common/logger"
	"gocourse/common/middleware"
	"gocourse/common/storage"
)

type Deps struct {
	DB     database.DB
	Cache  *storage.Storage
	Logger *logger.Logger
}

func Register(router gin.IRouter, deps Deps) {
	router.GET(
		"/:id",
		middleware.UUID("id"),
		user_handler.FindOneUser(
			user_handler.FindOneDeps{
				DB:     deps.DB,
				Cache:  deps.Cache,
				Logger: deps.Logger,
			},
		),
	)
	router.POST(
		"",
		user_handler.SaveOne(
			user_handler.SaveOneDeps{
				DB:     deps.DB,
				Logger: deps.Logger,
			},
		),
	)
}
