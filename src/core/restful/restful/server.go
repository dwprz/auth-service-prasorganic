package restful

import (
	"net/http"
	"time"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/middleware"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/router"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/gofiber/fiber/v2"
)

// this main restful server
type Server struct {
	app                *fiber.App
	authRestfulHandler *handler.AuthRestful
	middleware         *middleware.Middleware
	conf               *config.Config
}

func NewServer(arh *handler.AuthRestful, m *middleware.Middleware, conf *config.Config) *Server {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		IdleTimeout:   20 * time.Second,
		ReadTimeout:   20 * time.Second,
		WriteTimeout:  20 * time.Second,
		ErrorHandler:  m.Error,
	})

	router.AddAuth(app, arh, m)

	return &Server{
		app:                app,
		authRestfulHandler: arh,
		middleware:         m,
		conf:               conf,
	}
}

func (r *Server) Run() {
	r.app.Listen(r.conf.CurrentApp.RestfulAddress)
}

func (r *Server) Test(req *http.Request) (*http.Response, error) {
	res, err := r.app.Test(req)

	return res, err
}

func (r *Server) Stop() {
	r.app.Shutdown()
}