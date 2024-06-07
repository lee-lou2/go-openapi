package utils

import (
	"encoding/base64"
)

// Base64Encode 인코딩 함수
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode 디코딩 함수
func Base64Decode(encoded string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return data, nil
}
