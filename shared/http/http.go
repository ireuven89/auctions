package http

import (
	"crypto/rsa"
	"encoding/json"
	"log"
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
	jwksURL   string
	keys      map[string]interface{}
	mu        sync.RWMutex
	lastFetch time.Time
}

func NewJWKSProvider(jwksURL string) *JWKSProvider {
	p := &JWKSProvider{
		jwksURL: jwksURL,
		keys:    make(map[string]interface{}),
	}
	p.refreshKeys()    // initial load
	go p.autoRefresh() // periodic refresh in background
	return p
}

func (p *JWKSProvider) refreshKeys() {
	resp, err := http.Get(p.jwksURL)
	if err != nil {
		log.Printf("failed to fetch JWKS: %v", err)
		return
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		log.Printf("failed to decode JWKS: %v", err)
		return
	}

	newKeys := make(map[string]interface{})
	for _, raw := range jwks.Keys {
		key, err := jwt.ParseRSAPublicKeyFromPEM(raw)
		if err == nil {
			// In a real setup, you'd match by key ID (kid)
			newKeys["default"] = key
		}
	}

	p.mu.Lock()
	p.keys = newKeys
	p.lastFetch = time.Now()
	p.mu.Unlock()

	log.Println("JWKS keys refreshed")
}

func (p *JWKSProvider) autoRefresh() {
	for {
		time.Sleep(5 * time.Minute)
		p.refreshKeys()
	}
}

func (p *JWKSProvider) GetPublicKey() (interface{}, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	key, ok := p.keys["default"]
	return key, ok
}

// JWTMiddleware applies JWT validation to all routes except those in publicPaths.
func JWTMiddleware(publicKey *rsa.PublicKey, publicPaths []string) func(http.Handler) http.Handler {
	// Build a set for O(1) path lookups
	public := make(map[string]struct{}, len(publicPaths))
	for _, p := range publicPaths {
		public[p] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := public[r.URL.Path]; ok {
				next.ServeHTTP(w, r)
				return
			}
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return publicKey, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			// Optionally: put claims in context here
			next.ServeHTTP(w, r)
		})
	}
}
