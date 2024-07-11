package main

import (
	"github.com/dwprz/prasorganic-auth-service/src/app"
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
)

func main() {
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

	server := app.NewApp(middleware)
	router.AddAuthRouter(server, authHandler, middleware)

	server.Listen(conf.App.Address)
}
