package util

import (
	"encoding/json"
	"log"
	"strings"
)

func MarshalRequestBody(data any) *strings.Reader {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	reqBody := strings.NewReader(string(jsonData))

	return reqBody
}
