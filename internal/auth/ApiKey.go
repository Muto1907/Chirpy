package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	pkHeader := headers.Get("Authorization")
	if pkHeader == "" {
		return "", errors.New("API Key Header missing")
	}
	trimmedHeader := strings.TrimPrefix(pkHeader, "ApiKey ")
	return trimmedHeader, nil
}
