package handler

import "github.com/gofiber/fiber/v2"

type AuthHandler interface {
	Register(c *fiber.Ctx) error
}
