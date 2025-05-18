package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	token := make([]byte, 256)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	result := hex.EncodeToString(token)
	return result, nil
}
