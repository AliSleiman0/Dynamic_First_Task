package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenService manages JWT creation and related helpers.
type TokenService struct {
	secret []byte
	ttl    time.Duration
	issuer string
}

// NewTokenService creates a TokenService instance.
func NewTokenService(secret string, ttl time.Duration, issuer string) *TokenService {
	return &TokenService{
		secret: []byte(secret),
		ttl:    ttl,
		issuer: issuer,
	}
}

// CreateAccessToken creates a signed HS256 token with common claims.
// userID can be any identifier (int used here for example).
// extra may include fields like "role", "name", etc.
func (t *TokenService) CreateAccessToken(userID int, extra map[string]any) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"iss": t.issuer,
		"iat": now.Unix(),
		"exp": now.Add(t.ttl).Unix(),
	}

	// add optional extra claims in a safe way
	for k, v := range extra {
		if k == "sub" || k == "iss" || k == "iat" || k == "exp" {
			continue
		}
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.secret)
}

// ExpiresInSeconds returns the TTL in seconds as an exported helper for other packages.
func (t *TokenService) ExpiresInSeconds() int {
	return int(t.ttl.Seconds())
}
