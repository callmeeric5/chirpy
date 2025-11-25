package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("no authorization header")
	}

	const prefix = "ApiKey "
	if !strings.HasPrefix(auth, prefix) {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return strings.TrimPrefix(auth, prefix), nil
}
