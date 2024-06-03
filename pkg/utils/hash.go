package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA256 SHA-256 해시 생성
func SHA256(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	hashedBytes := hash.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)
	return hashedString
}
