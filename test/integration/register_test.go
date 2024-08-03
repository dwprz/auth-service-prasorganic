package test

import (
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/restful"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-auth-service/src/mock/client"
	"github.com/dwprz/prasorganic-auth-service/src/mock/helper"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	"github.com/dwprz/prasorganic-auth-service/test/util"
	"github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

// *nyalakan database nya terlebih dahulu
// go test -v ./test/integration/... -count=1 -p=1
// go test -run ^TestIntegration_Register$  -v ./test/integration -count=1

type RegisterTestSuite struct {
	suite.Suite
	restfulServer  *restful.Server
	userGrpcClient *client.UserGrpcMock
	redisDB        *redis.ClusterClient
	conf           *config.Config
	logger         *logrus.Logger
	helper         *helper.HelperMock
}

func (r *RegisterTestSuite) SetupSuite() {
	// mock
	r.userGrpcClient = client.NewUserMock()

	restfulServer, redisDB, conf, logger, helper := util.NewRestfulServer(r.userGrpcClient)
	r.restfulServer = restfulServer
	r.redisDB = redisDB
	r.conf = conf
	r.logger = logger
	r.helper = helper
}

func (r *RegisterTestSuite) TearDownSuite() {
	r.redisDB.Close()
}

func (r *RegisterTestSuite) Test_Success() {
	registerReq := &dto.RegisterReq{
		Email:    "johndoe123@gamil.com",
		FullName: "John Doe",
		Password: "rahasia",
	}

	r.MockUserGrpcClient_FindByEmail(registerReq.Email, nil)
	r.MockHelper_GenerateOtp("123456")

	request := r.CreateRegisterRequest(registerReq)
	res, err := r.restfulServer.Test(request)
	assert.NoError(r.T(), err)

	assert.Equal(r.T(), 200, res.StatusCode)

	resBody := util.UnmarshalResponseBody(res.Body)
	assert.NotNil(r.T(), resBody["data"])
}

func (r *RegisterTestSuite) Test_AlreadyExists() {
	registerReq := &dto.RegisterReq{
		Email:    "userexisted@gamil.com",
		FullName: "John Doe",
		Password: "rahasia",
	}

	r.MockUserGrpcClient_FindByEmail(registerReq.Email, new(user.User))

	request := r.CreateRegisterRequest(registerReq)
	res, err := r.restfulServer.Test(request)
	assert.NoError(r.T(), err)

	assert.Equal(r.T(), 409, res.StatusCode)

	resBody := util.UnmarshalResponseBody(res.Body)
	assert.NotNil(r.T(), resBody["errors"])
}

func (r *RegisterTestSuite) Test_InvalidInput() {
	registerReq := &dto.RegisterReq{
		Email:    "12345",
		FullName: "John Doe",
		Password: "rahasia",
	}

	request := r.CreateRegisterRequest(registerReq)
	res, err := r.restfulServer.Test(request)
	assert.NoError(r.T(), err)

	assert.Equal(r.T(), 400, res.StatusCode)

	resBody := util.UnmarshalResponseBody(res.Body)
	assert.NotNil(r.T(), resBody["errors"])
}

func (r *RegisterTestSuite) MockHelper_GenerateOtp(otp string) {
	r.helper.Mock.On("GenerateOtp").Return(otp, nil)
}

func (r *RegisterTestSuite) MockUserGrpcClient_FindByEmail(email string, data *user.User) {
	r.userGrpcClient.Mock.On("FindByEmail", mock.Anything, email).Return(
		&user.FindUserResponse{
			Data: data,
		}, nil)
}

func (r *RegisterTestSuite) CreateRegisterRequest(body *dto.RegisterReq) *http.Request {
	reqBody := util.MarshalRequestBody(body)

	request := httptest.NewRequest("POST", "/api/auth/register", reqBody)
	request.Header.Set("Content-Type", "application/json")
	return request
}

func TestIntegration_Register(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}
