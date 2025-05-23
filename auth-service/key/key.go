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

var ErrNotFound = errors.New("token not found")
