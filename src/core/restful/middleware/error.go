package middleware

import (
	"encoding/json"
	"strconv"

	errcustom "github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/common/errors/restful"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/status"
)

func (m *Middleware) Error(c *fiber.Ctx, err error) error {
	restful.LogError(c, err)

	if st, ok := status.FromError(err); ok {
		return restful.HandleGrpcError(c, st)
	}

	if validationError, ok := err.(validator.ValidationErrors); ok {
		return restful.HandleValidationError(c, validationError)
	}

	if responseError, ok := err.(*errcustom.Response); ok {
		return restful.HandleResponseError(c, responseError)
	}

	if jwtError := restful.HanldeJwtError(err); jwtError != nil {
		return c.Status(401).JSON(fiber.Map{"errors": jwtError.Error()})
	}

	if jsonError, ok := err.(*json.UnmarshalTypeError); ok {
		return restful.HandleJsonError(c, jsonError)
	}

	if strconvError, ok := err.(*strconv.NumError); ok {
		return restful.HandleStrconvError(c, strconvError)
	}

	return c.Status(500).JSON(fiber.Map{
		"errors": "sorry, internal server error try again later",
	})
}
