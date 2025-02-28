package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSecurePassword() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	password := base64.URLEncoding.EncodeToString(bytes)[:12]
	return password, nil
}
