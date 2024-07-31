package util

import (
	"encoding/json"
	"io"
	"log"
)

func UnmarshalResponseBody(data io.ReadCloser) map[string]any {
	dataByte, err := io.ReadAll(data)
	if err != nil {
		log.Fatal(err)
	}

	res := make(map[string]any)

	err = json.Unmarshal(dataByte, &res)
	if err != nil {
		log.Fatal(err)
	}

	return res
}
