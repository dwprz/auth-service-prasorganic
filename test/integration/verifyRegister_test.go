package test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	"golang.org/x/crypto/bcrypt"
)

// go test -v ./test/integration/... -count=1 -p=1
// go test -run ^TestIntegration_VerifyRegister$  -v ./test/integration -count=1

type VerifyRegisterTestSuite struct {
	suite.Suite
	restfulServer    *server.Restful
	userGrpcDelivery *delivery.UserGrpcMock
	redisDB          *redis.ClusterClient
	redistTestUtil   *utiltest.RedisTest
	util             *util.UtilMock
}

func (v *VerifyRegisterTestSuite) SetupSuite() {
	// mock
	v.util = util.NewMock()
	v.userGrpcDelivery = delivery.NewUserGrpcMock()

	v.redisDB = database.NewRedisCluster()

	authCache, otpCache := utiltest.InitCacheTest(v.redisDB)
	rabbitMQClient, _ := utiltest.InitRabbitMQ()
	otpService := utiltest.InitOtpService(rabbitMQClient, otpCache, v.util)

	grpcClient := utiltest.InitGrpcClientTest(v.userGrpcDelivery)
	authService := utiltest.InitAuthServiceTest(grpcClient, otpService, authCache)

	v.restfulServer = utiltest.InitRestfulTest(authService)
	v.redistTestUtil = utiltest.NewRedisTest(v.redisDB)
}

func (v *VerifyRegisterTestSuite) TearDownSuite() {
	v.redistTestUtil.Flushall()
	v.redisDB.Close()

	v.restfulServer.Stop()
}

func (v *VerifyRegisterTestSuite) Test_Success() {
	// register
	// *hati-hati menggunakan pointer dalam unit test karena bisa jadi value nya berubah setelah function dijalankan
	registerReq := &dto.RegisterReq{
		Email:    "johndoe123@gmail.com",
		FullName: "John Doe",
		Password: "rahasia",
	}
	const otp = "123456"

	v.MockUserGrpcDelivery_FindByEmail(registerReq.Email)
	v.MockHelper_GenerateOtp(otp)

	request := v.CreateRegisterRequest(registerReq)
	_, err := v.restfulServer.Test(request)
	assert.NoError(v.T(), err)

	time.Sleep(500 * time.Millisecond) // meberi waktu redis untuk menyebarkan data ke node lain dalam cluster

	// verify register
	verifyRegisterReq := &dto.VerifyOtpReq{
		Otp: otp,
	}

	v.MockUserGrpcDelivery_Create(registerReq)

	request = v.CreateVerifyRegisterRequest(verifyRegisterReq, registerReq.Email)
	res, err := v.restfulServer.Test(request)
	assert.NoError(v.T(), err)

	assert.Equal(v.T(), 200, res.StatusCode)
}

func (v *VerifyRegisterTestSuite) MockUserGrpcDelivery_FindByEmail(email string) {
	v.userGrpcDelivery.Mock.On("FindByEmail", mock.Anything, email).Return(
		&user.FindUserRes{
			Data: nil,
		}, nil,
	)
}

func (v *VerifyRegisterTestSuite) MockUserGrpcDelivery_Create(data *dto.RegisterReq) {

	v.userGrpcDelivery.Mock.On("Create", mock.Anything, mock.MatchedBy(func(req *user.RegisterReq) bool {
		err := bcrypt.CompareHashAndPassword([]byte(req.Password), []byte("rahasia"))
		return req.Email == data.Email && req.FullName == data.FullName && err == nil
	})).Return(nil)
}

func (v *VerifyRegisterTestSuite) MockHelper_GenerateOtp(otp string) {
	v.util.Mock.On("GenerateOtp").Return(otp, nil)
}

func (v *VerifyRegisterTestSuite) CreateRegisterRequest(body *dto.RegisterReq) *http.Request {
	reqBody := utiltest.MarshalRequestBody(body)

	request := httptest.NewRequest("POST", "/api/auth/register", reqBody)
	request.Header.Set("Content-Type", "application/json")
	return request
}

func (v *VerifyRegisterTestSuite) CreateVerifyRegisterRequest(body *dto.VerifyOtpReq, email string) *http.Request {
	reqBody := utiltest.MarshalRequestBody(body)

	request := httptest.NewRequest("POST", "/api/auth/register/verify", reqBody)

	request.AddCookie(&http.Cookie{
		Name:     "pending_register",
		Value:    base64.StdEncoding.EncodeToString([]byte(email)),
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(30 * time.Minute),
	})

	request.Header.Set("Content-Type", "application/json")
	return request
}

func TestIntegration_VerifyRegister(t *testing.T) {
	suite.Run(t, new(VerifyRegisterTestSuite))
}
