package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts the API key from the request headers.
// Examle:
// Authorization: ApiKey <key>
func GetAPIKey(headers http.Header) (string, error) {

	val := headers.Get("Authorization")

	if val == "" {
		return "", errors.New("no authentication info found")
	}

	vals := strings.Split(val, " ")

	if len(vals) != 2 {
		return "", errors.New("malformed auth bearer")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("malformed auth bearer")
	}

	return vals[1], nil
}
