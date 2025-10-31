package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"gocourse/app/api/endpoint/user_endpoint"
	"gocourse/common/cache"
	"gocourse/common/config"
	"gocourse/common/config/cache_config"
	"gocourse/common/config/database_config"
	"gocourse/common/database"
	"gocourse/common/logger"
	"gocourse/common/middleware"
	"gocourse/common/storage"
	"gocourse/pkg/env"
)

func main() {
	if env.Production() || env.Test() {
		gin.SetMode(gin.ReleaseMode)
	}
	cfg := config.Load()
	app := gin.New()
	{
		// Middleware
		app.Use(middleware.ErrorHandler())
		app.Use(middleware.JSON())
	}
	var customLogger *logger.Logger
	var postgresDB database.DB
	var valkeyClient *redis.Client
	var valkeyStorage *storage.Storage
	{
		// Deps initialization
		var err error
		customLogger = logger.New()
		if postgresDB, err = database.Connect(cfg.GetString(database_config.Uri)); err != nil {
			log.Fatalf("connect to database failed: %s\n", err)
		}
		if valkeyClient, err = cache.Connect(
			cfg.GetString(cache_config.Host),
			cfg.GetString(cache_config.Password),
			cfg.GetInt(cache_config.DB),
		); err != nil {
			log.Fatalf("connect to cache failed: %s\n", err)
		}
		valkeyStorage = storage.New(valkeyClient)
	}
	{
		// Endpoints
		user_endpoint.Register(
			app.Group(
				"/users",
			),
			user_endpoint.Deps{
				DB:     postgresDB,
				Cache:  valkeyStorage,
				Logger: customLogger,
			},
		)
	}
	log.Fatalln(app.Run(":" + os.Getenv("APP_PORT")))
}
