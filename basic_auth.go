package addon

import (
	"encoding/base64"
	"errors"
	"strings"
)

var (
	ErrNoBasicAuth = errors.New("not a basic authentication")
)

type BasicAuth struct {
	Username string
	Password string
}

func parseBasicAuth(data string) (BasicAuth, error) {
	auth := BasicAuth{}

	// Check if auth string looks like basic auth data
	if !strings.HasPrefix(data, "Basic ") {
		return auth, ErrNoBasicAuth
	}

	// Decode the username and password contents
	decoded, err := base64.StdEncoding.DecodeString(strings.Replace(data, "Basic ", "", 1))
	if err != nil {
		return auth, ErrNoBasicAuth
	}

	// Extract username and password
	chunks := strings.Split(string(decoded), ":")
	auth.Username = chunks[0]
	auth.Password = chunks[1]

	return auth, nil
}
