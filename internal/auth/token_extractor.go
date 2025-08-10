package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	tokenString, ok := headers["Authorization"]
	if !ok {
		return "", fmt.Errorf("missing authorization header")
	}
	return strings.ReplaceAll(tokenString[0], "ApiKey ", ""), nil
}

func GetBearerToken(headers http.Header) (string, error) {
	tokenString, ok := headers["Authorization"]
	if !ok {
		return "", fmt.Errorf("missing authorization header")
	}
	return strings.ReplaceAll(tokenString[0], "Bearer ", ""), nil
}
