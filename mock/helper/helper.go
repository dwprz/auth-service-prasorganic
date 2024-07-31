package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type HelperMock struct {
	mock.Mock
}

func NewMock() *HelperMock {
	return &HelperMock{
		Mock: mock.Mock{},
	}
}

func (h *HelperMock) GenerateOtp() string {
	arguments := h.Mock.Called()

	return arguments.String(0)
}

func (h *HelperMock) HandlePanic(name string, c *fiber.Ctx, logger *logrus.Logger) {}
