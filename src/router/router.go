package router

import (
	"github.com/dwprz/prasorganic-auth-service/src/handler"
	"github.com/dwprz/prasorganic-auth-service/src/middleware"
	"github.com/gofiber/fiber/v2"
)

func AddAuthRouter(app *fiber.App, authHandler handler.AuthHandler, m middleware.Middleware) {
	app.Add("POST", "/api/auth/register", authHandler.Register)
}
