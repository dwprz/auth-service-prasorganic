package service

import (
	"context"
	"encoding/json"

	"github.com/dwprz/prasorganic-auth-service/src/cache"
	"github.com/dwprz/prasorganic-auth-service/src/common/config"
	"github.com/dwprz/prasorganic-auth-service/src/common/custom_error"
	"github.com/dwprz/prasorganic-auth-service/src/common/helper"
	"github.com/dwprz/prasorganic-auth-service/src/common/model/dto"
	"github.com/dwprz/prasorganic-auth-service/src/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthServiceImpl struct {
	AuthRepository repository.AuthRepository
	Validate       *validator.Validate
	Chace          cache.AuthCache
	Conf           *config.Config
	Logger         *logrus.Logger
}

func NewAuthService(r repository.AuthRepository, v *validator.Validate, c cache.AuthCache, conf *config.Config, l *logrus.Logger) AuthService {
	return &AuthServiceImpl{
		AuthRepository: r,
		Validate:       v,
		Chace:          c,
		Conf:           conf,
		Logger:         l,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, data *dto.RegisterReq) (string, error) {
	if err := s.Validate.Struct(data); err != nil {
		return "", &custom_error.ValidationError{
			Name:    "register (validate)",
			Message: err.Error(),
		}
	}

	credential := s.AuthRepository.FindCredentialByEmail(ctx, data.Email)
	if credential != nil {
		return "", &custom_error.ResponseError{Code: 409, Message: "email already exists"}
	}

	otp := helper.GenerateOtp()
	data.Otp = otp

	if err := s.Chace.CacheRegisterReq(ctx, data); err != nil {
		return "", err
	}

	s.SendOtp(ctx, otp, data.Email)
	return data.Email, nil
}

func (s *AuthServiceImpl) SendOtp(ctx context.Context, otp string, email string) {

	client := fiber.AcquireClient()
	defer fiber.ReleaseClient(client)

	agent := client.Post(s.Conf.Other.EmailServiceUrl)
	defer fiber.ReleaseAgent(agent)

	data := map[string]string{
		"email": email,
		"otp":   otp,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		s.Logger.Errorf("send otp (marshal): %v", err.Error())
	}

	req := agent.Request()
	req.SetBody(jsonData)
	req.Header.SetContentType("application/json")
	req.Header.Set("Authorization", "rahasia")

	_, _, errors := agent.String()
	if len(errors) != 0 {
		s.Logger.Errorf("send otp (agent): %v", errors[0])
	}
}
