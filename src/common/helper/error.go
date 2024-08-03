package helper

import "github.com/gofiber/fiber/v2"

func (h *HelperImpl) HandlePanic(name string, c *fiber.Ctx) {
	message := recover()

	if message != nil {
		h.logger.Errorf(name+" | %v", message)

		c.Status(500).JSON(fiber.Map{
			"errors": "sorry, internal server error try again later",
		})
	}
}
