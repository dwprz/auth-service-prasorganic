package helper

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateOauthState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)
	return state, nil
}
