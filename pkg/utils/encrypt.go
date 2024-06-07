package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"go-openapi/config"
	"io"
	"sync"
)

type AES256 struct {
	Key string
}

var (
	aesInstance *AES256
	aesOnce     sync.Once
)

// Encrypt AES 암호화
func (c *AES256) Encrypt(stringToEncrypt string) (encryptedString string) {
	key, _ := hex.DecodeString(c.Key)
	plaintext := []byte(stringToEncrypt)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

// Decrypt AES 복호화
func (c *AES256) Decrypt(encryptedString string) (decryptedString string) {
	key, _ := hex.DecodeString(c.Key)
	enc, _ := hex.DecodeString(encryptedString)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("%s", plaintext)
}

// NewAES256 AES256 암호화
func NewAES256() *AES256 {
	aesOnce.Do(func() {
		key := config.GetEnv("AES256_KEY")
		aesInstance = &AES256{
			Key: key,
		}
	})
	return aesInstance
}
