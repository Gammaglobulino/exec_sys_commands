package storing_secure_password

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
)

var secretKey = "---" //deprecate--- must be placed to a secure place..

func GenerateSalt() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}
func HashPassword(plainText string) string {
	hash := hmac.New(sha256.New, []byte(secretKey))
	io.WriteString(hash, plainText)
	hashedValue := hash.Sum(nil)
	return hex.EncodeToString(hashedValue)
}
