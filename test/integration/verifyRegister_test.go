package test

import (
	"encoding/base64"
	"github.com/dwprz/prasorganic-auth-service/src/mock/client"
	"github.com/dwprz/prasorganic-auth-service/src/mock/helper"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/restful"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	"github.com/dwprz/prasorganic-auth-service/test/util"
	"github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// go test -v ./test/integration/... -count=1 -p=1
// go test -run ^TestIntegration_VerifyRegister$  -v ./test/integration -count=1

type VerifyRegisterTestSuite struct {
	suite.Suite
	restfulServer  *restful.Server
	userGrpcClient *client.UserGrpcMock
	redisDB        *redis.ClusterClient
	conf           *config.Config
	logger         *logrus.Logger
	helper         *helper.HelperMock
}

func (v *VerifyRegisterTestSuite) SetupSuite() {
	// mock
	v.userGrpcClient = client.NewUserMock()

	restfulServer, redisDB, conf, logger, helper := util.NewRestfulServer(v.userGrpcClient)
	v.restfulServer = restfulServer
	v.redisDB = redisDB
	v.conf = conf
	v.logger = logger
	v.helper = helper
}

func (v *VerifyRegisterTestSuite) TearDownSuite() {
	v.redisDB.Close()
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

	v.MockUserGrpcClient_FindByEmail(registerReq.Email)
	v.MockHelper_GenerateOtp(otp)

	request := v.CreateRegisterRequest(registerReq)
	_, err := v.restfulServer.Test(request)
	assert.NoError(v.T(), err)

	// verify register
	verifyRegisterReq := &dto.VerifyRegisterReq{
		Otp: otp,
	}

	v.MockUserGrpcClient_Create(registerReq)

	request = v.CreateVerifyRegisterRequest(verifyRegisterReq, registerReq.Email)
	res, err := v.restfulServer.Test(request)
	assert.NoError(v.T(), err)

	assert.Equal(v.T(), 200, res.StatusCode)
}

func (v *VerifyRegisterTestSuite) MockUserGrpcClient_FindByEmail(email string) {

	v.userGrpcClient.Mock.On("FindByEmail", mock.Anything, &user.Email{
		Email: email,
	}).Return(&user.FindUserResponse{Data: nil}, nil)
}

func (v *VerifyRegisterTestSuite) MockUserGrpcClient_Create(data *dto.RegisterReq) {

	v.userGrpcClient.Mock.On("Create", mock.Anything, mock.MatchedBy(func(req *user.RegisterRequest) bool {
		err := bcrypt.CompareHashAndPassword([]byte(req.Password), []byte(data.Password))
		return req.Email == data.Email && req.FullName == data.FullName && err == nil
	})).Return(nil)
}

func (v *VerifyRegisterTestSuite) MockHelper_GenerateOtp(otp string) {
	v.helper.Mock.On("GenerateOtp").Return(otp, nil)
}

func (v *VerifyRegisterTestSuite) CreateRegisterRequest(body *dto.RegisterReq) *http.Request {
	reqBody := util.MarshalRequestBody(body)

	request := httptest.NewRequest("POST", "/api/auth/register", reqBody)
	request.Header.Set("Content-Type", "application/json")
	return request
}

func (v *VerifyRegisterTestSuite) CreateVerifyRegisterRequest(body *dto.VerifyRegisterReq, email string) *http.Request {
	reqBody := util.MarshalRequestBody(body)

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
