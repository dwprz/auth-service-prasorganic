package test

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"github.com/dwprz/prasorganic-auth-service/mock/helper"
	"github.com/dwprz/prasorganic-auth-service/mock/service"
	"github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/common/logger"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/middleware"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/restful"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	"github.com/dwprz/prasorganic-auth-service/test/util"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// go test -v ./src/core/restful/handler/test/... -count=1 -p=1 
// go test -run ^TestHandler_VerifyRegister$  -v ./src/core/restful/handler/test/ -count=1

type VerifyRegisterTestSuite struct {
	suite.Suite
	restfulServer *restful.Server
	authService   *service.AuthMock
	logger        *logrus.Logger
	helper        *helper.HelperMock
}

func (v *VerifyRegisterTestSuite) SetupSuite() {
	v.logger = logger.New()
	conf := config.New("DEVELOPMENT", v.logger)
	v.helper = helper.NewMock()

	// mock
	v.authService = service.NewAuthMock()
	authHandler := handler.NewAuthRestful(v.authService, v.logger, v.helper)

	middleware := middleware.New(conf, v.logger)
	v.restfulServer = restful.NewServer(authHandler, middleware, conf)
}

func (v *VerifyRegisterTestSuite) Test_Success() {

	data := &dto.VerifyRegisterReq{
		Email: "johndoe123@gmail.com",
		Otp:   "123456",
	}

	v.authService.Mock.On("VerifyRegister", context.Background(), data).Return(nil)

	reqBody := util.MarshalRequestBody(data)

	request := httptest.NewRequest("POST", "/api/auth/register/verify", reqBody)

	request.AddCookie(&http.Cookie{
		Name:     "pending_register",
		Value:    base64.StdEncoding.EncodeToString([]byte(data.Email)),
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(30 * time.Minute),
	})

	request.Header.Set("Content-Type", "application/json")

	res, err := v.restfulServer.Test(request)
	assert.NoError(v.T(), err)

	assert.Equal(v.T(), 200, res.StatusCode)

	resBody := util.UnmarshalResponseBody(res.Body)
	assert.NotEmpty(v.T(), resBody["data"])
}

func (v *VerifyRegisterTestSuite) Test_InvalidOtp() {
	data := &dto.VerifyRegisterReq{
		Email: "johndoe123@gmail.com",
		Otp:   "invalid otp",
	}

	errorRes := &errors.Response{Code: 400, Message: "otp is invalid"}
	v.authService.Mock.On("VerifyRegister", context.Background(), data).Return(errorRes)

	reqBody := util.MarshalRequestBody(data)

	request := httptest.NewRequest("POST", "/api/auth/register/verify", reqBody)

	request.AddCookie(&http.Cookie{
		Name:     "pending_register",
		Value:    base64.StdEncoding.EncodeToString([]byte(data.Email)),
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(30 * time.Minute),
	})

	request.Header.Set("Content-Type", "application/json")

	res, err := v.restfulServer.Test(request)
	assert.NoError(v.T(), err)

	assert.Equal(v.T(), 400, res.StatusCode)

	resBody := util.UnmarshalResponseBody(res.Body)
	assert.NotEmpty(v.T(), resBody["errors"])
}

func (v *VerifyRegisterTestSuite) Test_InvalidEmail() {

	data := &dto.VerifyRegisterReq{
		Email: "invalid email",
		Otp:   "123456",
	}

	errorRes := &errors.Response{Code: 400, Message: "email is invalid"}
	v.authService.Mock.On("VerifyRegister", context.Background(), data).Return(errorRes)

	reqBody := util.MarshalRequestBody(data)

	request := httptest.NewRequest("POST", "/api/auth/register/verify", reqBody)

	request.AddCookie(&http.Cookie{
		Name:     "pending_register",
		Value:    base64.StdEncoding.EncodeToString([]byte(data.Email)),
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(30 * time.Minute),
	})

	request.Header.Set("Content-Type", "application/json")

	res, err := v.restfulServer.Test(request)
	assert.NoError(v.T(), err)

	assert.Equal(v.T(), 400, res.StatusCode)

	resBody := util.UnmarshalResponseBody(res.Body)
	assert.NotEmpty(v.T(), resBody["errors"])
}

func TestHandler_VerifyRegister(t *testing.T) {
	suite.Run(t, new(VerifyRegisterTestSuite))
}
