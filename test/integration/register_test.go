package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dwprz/prasorganic-auth-service/src/core/restful/server"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-auth-service/src/mock/delivery"
	"github.com/dwprz/prasorganic-auth-service/src/mock/util"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	utiltest "github.com/dwprz/prasorganic-auth-service/test/util"
	"github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// *nyalakan database nya terlebih dahulu
// go test -v ./test/integration/... -count=1 -p=1
// go test -run ^TestIntegration_Register$  -v ./test/integration -count=1

type RegisterTestSuite struct {
	suite.Suite
	restfulServer    *server.Restful
	userGrpcDelivery *delivery.UserGrpcMock
	redisDB          *redis.ClusterClient
	redistTestUtil   *utiltest.RedisTest
	util             *util.UtilMock
}

func (r *RegisterTestSuite) SetupSuite() {
	// mock
	r.util = util.NewMock()
	r.userGrpcDelivery = delivery.NewUserGrpcMock()

	r.redisDB = database.NewRedisCluster()

	authCache, otpCache := utiltest.InitCacheTest(r.redisDB)
	rabbitMQClient, _ := utiltest.InitRabbitMQ()
	otpService := utiltest.InitOtpService(rabbitMQClient, otpCache, r.util)

	grpcClient := utiltest.InitGrpcClientTest(r.userGrpcDelivery)
	authService := utiltest.InitAuthServiceTest(grpcClient, otpService, authCache)

	r.restfulServer = utiltest.InitRestfulTest(authService)
	r.redistTestUtil = utiltest.NewRedisTest(r.redisDB)
}

func (r *RegisterTestSuite) TearDownSuite() {
	//r.redistTestUtil.Flushall()
	r.redisDB.Close()

	r.restfulServer.Stop()
}

func (r *RegisterTestSuite) Test_Success() {
	registerReq := &dto.RegisterReq{
		Email:    "johndoe123@gamil.com",
		FullName: "John Doe",
		Password: "rahasia",
	}

	r.MockUserGrpcDelivery_FindByEmail(registerReq.Email, nil)
	r.MockHelper_GenerateOtp("123456")

	request := r.CreateRegisterRequest(registerReq)
	res, err := r.restfulServer.Test(request)
	assert.NoError(r.T(), err)

	assert.Equal(r.T(), 200, res.StatusCode)

	resBody := utiltest.UnmarshalResponseBody(res.Body)
	assert.NotNil(r.T(), resBody["data"])
}

func (r *RegisterTestSuite) Test_AlreadyExists() {
	registerReq := &dto.RegisterReq{
		Email:    "userexisted@gamil.com",
		FullName: "John Doe",
		Password: "rahasia",
	}

	r.MockUserGrpcDelivery_FindByEmail(registerReq.Email, new(user.User))

	request := r.CreateRegisterRequest(registerReq)
	res, err := r.restfulServer.Test(request)
	assert.NoError(r.T(), err)

	assert.Equal(r.T(), 409, res.StatusCode)

	resBody := utiltest.UnmarshalResponseBody(res.Body)
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

	resBody := utiltest.UnmarshalResponseBody(res.Body)
	assert.NotNil(r.T(), resBody["errors"])
}

func (r *RegisterTestSuite) MockHelper_GenerateOtp(otp string) {
	r.util.Mock.On("GenerateOtp").Return(otp, nil)
}

func (r *RegisterTestSuite) MockUserGrpcDelivery_FindByEmail(email string, data *user.User) {
	r.userGrpcDelivery.Mock.On("FindByEmail", mock.Anything, email).Return(
		&user.FindUserRes{
			Data: data,
		}, nil)
}

func (r *RegisterTestSuite) CreateRegisterRequest(body *dto.RegisterReq) *http.Request {
	reqBody := utiltest.MarshalRequestBody(body)

	request := httptest.NewRequest("POST", "/api/auth/register", reqBody)
	request.Header.Set("Content-Type", "application/json")
	return request
}

func TestIntegration_Register(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}
