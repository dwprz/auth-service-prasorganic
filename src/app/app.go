package app

import (
	"time"

	"github.com/dwprz/prasorganic-auth-service/src/cache"
	"github.com/dwprz/prasorganic-auth-service/src/common/config"
	"github.com/dwprz/prasorganic-auth-service/src/common/log"
	"github.com/dwprz/prasorganic-auth-service/src/database"
	"github.com/dwprz/prasorganic-auth-service/src/handler"
	"github.com/dwprz/prasorganic-auth-service/src/middleware"
	"github.com/dwprz/prasorganic-auth-service/src/repository"
	"github.com/dwprz/prasorganic-auth-service/src/router"
	"github.com/dwprz/prasorganic-auth-service/src/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func NewApp(m middleware.Middleware) *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		IdleTimeout:   20 * time.Second,
		ReadTimeout:   20 * time.Second,
		WriteTimeout:  20 * time.Second,
		ErrorHandler:  m.ErrorMiddleware,
	})

	return app
}

func NewAppTest() (*fiber.App, *config.Config, *redis.ClusterClient) {
	logger := log.NewLogger()
	conf := config.NewConfig(logger)
	postgres := database.NewPostgres(conf, logger)
	redis := database.NewRedis(conf)
	authCache := cache.NewAuthCache(redis)
	validate := validator.New()

	authRepository := repository.NewAuthRepository(postgres)
	authService := service.NewAuthService(authRepository, validate, authCache, conf, logger)
	authHandler := handler.NewAuthHandler(authService, logger)
	middleware := middleware.NewMiddleware(conf, logger)

	appServer := NewApp(middleware)
	router.AddAuthRouter(appServer, authHandler, middleware)

	return appServer, conf, redis
}
