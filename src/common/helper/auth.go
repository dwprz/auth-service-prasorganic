package helper

import (
	"fmt"

	"math/rand"
)

func GenerateOtp() string {
	otp := rand.Intn(1000000)
	return fmt.Sprintf("%06d", otp)
}
