package server

import (
	"net/http"
	"time"

	"github.com/dwprz/prasorganic-auth-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/middleware"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/router"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// this main restful server
type Restful struct {
	app                *fiber.App
	authRestfulHandler *handler.Auth
	middleware         *middleware.Middleware
}

func NewRestful(arh *handler.Auth, m *middleware.Middleware) *Restful {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		IdleTimeout:   20 * time.Second,
		ReadTimeout:   20 * time.Second,
		WriteTimeout:  20 * time.Second,
		ErrorHandler:  m.Error,
	})

	app.Use(recover.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://restful.local:80",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))

	router.Create(app, arh, m)

	return &Restful{
		app:                app,
		authRestfulHandler: arh,
		middleware:         m,
	}
}

func (r *Restful) Run() {
	r.app.Listen(config.Conf.CurrentApp.RestfulAddress)
}

func (r *Restful) Test(req *http.Request) (*http.Response, error) {
	res, err := r.app.Test(req)

	return res, err
}

func (r *Restful) Stop() {
	r.app.Shutdown()
}