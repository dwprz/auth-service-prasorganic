package router

import (
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/middleware"
	"github.com/gofiber/fiber/v2"
)

func AddAuth(app *fiber.App, h *handler.AuthRestful, m *middleware.Middleware) {
	app.Add("POST", "/api/auth/register", h.Register)
	app.Add("POST", "/api/auth/register/verify", h.VerifyRegister)
}
