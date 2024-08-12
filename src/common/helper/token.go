package helper

import (
	"fmt"
	"time"

	"github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/golang-jwt/jwt/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"google.golang.org/grpc/codes"
)

func GenerateAccessToken(userId string, email string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":     "prasorganic-auth-service",
		"user_id": userId,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})

	accessToken, err := token.SignedString(config.Conf.Jwt.PrivateKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func GenerateRefreshToken() (string, error) {
	tokenId, err := gonanoid.New()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "prasorganic-auth-service",
		"id":  tokenId,
		"exp": time.Now().Add(24 * 30 * time.Hour).Unix(),
	})

	accessToken, err := token.SignedString(config.Conf.Jwt.PrivateKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func VerifyJwt(token string) (*jwt.MapClaims, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected token method: %v", t.Header["alg"])
		}

		return config.Conf.Jwt.PublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		return &claims, nil
	}

	return nil, &errors.Response{HttpCode: 401, GrpcCode: codes.Unauthenticated, Message: "token is invalid"}
}
