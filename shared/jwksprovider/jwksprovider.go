package jwksprovider

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWKS struct {
	Keys []json.RawMessage `json:"keys"`
}

type JWKSProvider struct {
	mu       sync.RWMutex
	keys     map[string]interface{}
	jwksURL  string
	interval time.Duration
}

func NewJWKSProvider(jwksURL string, refreshInterval time.Duration) *JWKSProvider {
	p := &JWKSProvider{
		keys:     make(map[string]interface{}),
		jwksURL:  jwksURL,
		interval: refreshInterval,
	}
	go p.startRefreshing()
	return p
}

func (p *JWKSProvider) startRefreshing() {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		p.refreshKeys()
		<-ticker.C
	}
}

func (p *JWKSProvider) refreshKeys() {
	resp, err := http.Get(p.jwksURL)
	if err != nil {
		return // optionally log error
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err = json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return // optionally log error
	}

	updatedKeys := make(map[string]interface{})
	for _, raw := range jwks.Keys {
		var key jwt.MapClaims // or specific JWK type
		if err := json.Unmarshal(raw, &key); err != nil {
			continue
		}
		if kid, ok := key["kid"].(string); ok {
			updatedKeys[kid] = key
		}
	}

	p.mu.Lock()
	p.keys = updatedKeys
	p.mu.Unlock()
}

func (p *JWKSProvider) GetKey(kid string) (interface{}, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	key, ok := p.keys[kid]
	if !ok {
		return nil, errors.New("key not found")
	}
	return key, nil
}

// ✅ JWTMiddleware example using kid
func JWTMiddleware(jwksProvider *JWKSProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				kid, ok := token.Header["kid"].(string)
				if !ok {
					return nil, errors.New("kid header not found")
				}
				return jwksProvider.GetKey(kid)
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "jwt_claims", token.Claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
