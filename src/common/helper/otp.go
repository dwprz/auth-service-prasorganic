package helper

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func (h *HelperImpl) GenerateOtp() (string, error) {
	max := big.NewInt(1000000)
	otp, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", otp), nil
}


