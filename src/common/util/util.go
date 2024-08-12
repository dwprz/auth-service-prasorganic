package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"github.com/dwprz/prasorganic-auth-service/src/interface/util"
)

// *untuk utilitas yang jika ditest memiliki kebutuhan untuk mocking

type UtilImpl struct{}

func New() util.Util {
	return &UtilImpl{}
}

func (u *UtilImpl) GenerateOtp() (string, error) {
	max := big.NewInt(1000000)
	otp, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", otp), nil
}

