package helper

import (
	"github.com/gofiber/fiber/v2"
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

func (h *HelperMock) GenerateOtp() (string, error) {
	arguments := h.Mock.Called()

	return arguments.String(0), arguments.Error(1)
}

func (h *HelperMock) GenerateOauthState() (string, error) {
	arguments := h.Mock.Called()

	return arguments.String(0), arguments.Error(1)
}

func (h *HelperMock) GenerateAccessToken(userId string, email string, role string) (string, error) {
	arguments := h.Mock.Called(userId, email, role)

	return arguments.String(0), arguments.Error(1)
}

func (h *HelperMock) GenerateRefreshToken() (string, error) {
	arguments := h.Mock.Called()

	return arguments.String(0), arguments.Error(1)
}

func (h *HelperMock) HandlePanic(name string, c *fiber.Ctx) {}
