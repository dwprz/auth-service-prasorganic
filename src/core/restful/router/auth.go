package router

import (
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/middleware"
	"github.com/gofiber/fiber/v2"
)

func Create(app *fiber.App, h *handler.AuthRESTful, m *middleware.Middleware) {
	// all
	app.Add("POST", "/api/auth/register", h.Register)
	app.Add("POST", "/api/auth/register/verify", h.VerifyRegister)
	app.Add("POST", "/api/auth/login", h.Login)
	app.Add("POST", "/api/auth/token/refresh", h.RefreshToken)
	app.Add("POST", "/api/auth/logout", h.Logout)
	app.Add("GET", "/api/auth/login/google", h.LoginWithGoogle)
	app.Add("GET", "/api/auth/login/google/callback", m.FetchGoogleUser, h.LoginWithGoogleCallback)
}
