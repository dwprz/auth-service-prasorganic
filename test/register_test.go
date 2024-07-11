package test

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/dwprz/prasorganic-auth-service/src/app"
	"github.com/dwprz/prasorganic-auth-service/src/common/config"
	"github.com/dwprz/prasorganic-auth-service/test/util_test"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RegisterTestSuite struct {
	suite.Suite
	App   *fiber.App
	Conf  *config.Config
	Redis *redis.ClusterClient
	Util  util_test.UtilTest
}

func (s *RegisterTestSuite) SetupSuite() {
	appServer, conf, redis := app.NewAppTest()
	s.App = appServer
	s.Conf = conf
	s.Redis = redis

	ctx := context.Background()
	s.Util = util_test.NewUtilTest(ctx, conf)
}

func (s *RegisterTestSuite) TearDownSuite() {
	err := s.Redis.Close()
	if err != nil {
		log.Printf("error redis (close): %+v\n", err.Error())
	}
}

func (s *RegisterTestSuite) TestRegister_Success() {
	registerReq := strings.NewReader(`{
		"email": "johndoe123@gmail.com", 
		"full_name": "John Doe", 
		"password": "rahasia"
		}`)

	request := httptest.NewRequest("POST", "/api/auth/register", registerReq)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", s.Conf.App.AuthSecretKey)
	res, err := s.App.Test(request)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), 200, res.StatusCode)

	data := s.parseResponseBody(res.Body)
	assert.NotNil(s.T(), data["data"])
}

func (s *RegisterTestSuite) TestRegister_InvalidInput() {
	registerReq := strings.NewReader(`{
		"email": 12345, 
		"full_name": "John Doe", 
		"password": "rahasia"
		}`)

	request := httptest.NewRequest("POST", "/api/auth/register", registerReq)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", s.Conf.App.AuthSecretKey)
	res, err := s.App.Test(request)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), 400, res.StatusCode)

	data := s.parseResponseBody(res.Body)
	assert.NotNil(s.T(), data["errors"])
}

func (s *RegisterTestSuite) TestRegister_InvalidAuthKey() {
	registerReq := strings.NewReader(`{
		"email": "johndoe123@gmail.com", 
		"full_name": "John Doe", 
		"password": "rahasia"
		}`)

	request := httptest.NewRequest("POST", "/api/auth/register", registerReq)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "INVALID AUTH KEY")
	res, err := s.App.Test(request)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), 401, res.StatusCode)

	data := s.parseResponseBody(res.Body)
	log.Printf("%+v\n", data)
	assert.NotNil(s.T(), data["errors"])
}

func (s *RegisterTestSuite) parseResponseBody(body io.ReadCloser) map[string]any {
	byte, err := io.ReadAll(body)
	assert.NoError(s.T(), err)

	data := make(map[string]any)

	err = json.Unmarshal(byte, &data)
	assert.NoError(s.T(), err)

	return data
}

func TestRegister(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}
