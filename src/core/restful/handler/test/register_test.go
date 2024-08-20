package test

import (
	"net/http/httptest"
	"testing"

	"github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/middleware"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/server"
	"github.com/dwprz/prasorganic-auth-service/src/mock/service"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	"github.com/dwprz/prasorganic-auth-service/test/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// go test -v ./src/core/restful/handler/test/... -count=1 -p=1
// go test -run ^TestHandler_Register$  -v ./src/core/restful/handler/test/ -count=1

type RegisterTestSuite struct {
	suite.Suite
	restfulServer *server.Restful
	authService   *service.AuthMock
}

func (r *RegisterTestSuite) SetupSuite() {
	// mock
	r.authService = service.NewAuthMock()

	authHandler := handler.NewAuthRESTful(r.authService)

	middleware := middleware.New()
	r.restfulServer = server.NewRestful(authHandler, middleware)
}

func (r *RegisterTestSuite) Test_Success() {

	data := &dto.RegisterReq{
		Email:    "johndoe123@gamil.com",
		FullName: "John Doe",
		Password: "rahasia",
	}

	r.authService.Mock.On("Register", mock.Anything, data).Return("johndoe123@gamil.com", nil)

	reqBody := util.MarshalRequestBody(data)

	request := httptest.NewRequest("POST", "/api/auth/register", reqBody)
	request.Header.Set("Content-Type", "application/json")

	res, err := r.restfulServer.Test(request)
	assert.NoError(r.T(), err)

	assert.Equal(r.T(), 200, res.StatusCode)

	resBody := util.UnmarshalResponseBody(res.Body)
	assert.NotEmpty(r.T(), resBody["data"])
}

func (r *RegisterTestSuite) Test_AlreadyExists() {

	data := &dto.RegisterReq{
		Email:    "userexisted@gamil.com",
		FullName: "John Doe",
		Password: "rahasia",
	}

	errorRes := &errors.Response{HttpCode: 409, Message: "user already exists"}
	r.authService.Mock.On("Register", mock.Anything, data).Return("", errorRes)

	reqBody := util.MarshalRequestBody(data)

	request := httptest.NewRequest("POST", "/api/auth/register", reqBody)
	request.Header.Set("Content-Type", "application/json")

	res, err := r.restfulServer.Test(request)
	assert.NoError(r.T(), err)

	assert.Equal(r.T(), 409, res.StatusCode)

	resBody := util.UnmarshalResponseBody(res.Body)
	assert.NotEmpty(r.T(), resBody["errors"])
}

func (r *RegisterTestSuite) Test_InvalidEmail() {

	data := &dto.RegisterReq{
		Email:    "12345",
		FullName: "John Doe",
		Password: "rahasia",
	}

	errorRes := &errors.Response{HttpCode: 400, Message: "email is invalid"}
	r.authService.Mock.On("Register", mock.Anything, data).Return("", errorRes)

	reqBody := util.MarshalRequestBody(data)

	request := httptest.NewRequest("POST", "/api/auth/register", reqBody)
	request.Header.Set("Content-Type", "application/json")

	res, err := r.restfulServer.Test(request)
	assert.NoError(r.T(), err)

	assert.Equal(r.T(), 400, res.StatusCode)

	resBody := util.UnmarshalResponseBody(res.Body)
	assert.NotEmpty(r.T(), resBody["errors"])
}

func TestHandler_Register(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}
