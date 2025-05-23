package key

import (
	"errors"
)

type JWK struct {
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type Token struct {
	Access  string
	Refresh string
}

var (
	ErrUserNotFound       = errors.New("user not found or credentials missing")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Authorization errors (token required)
var (
	ErrInvalidToken = errors.New("unauthorized: invalid token")
	ErrExpiredToken = errors.New("unauthorized: expired token")
)
