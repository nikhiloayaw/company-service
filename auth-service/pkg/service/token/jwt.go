package token

import (
	"auth-service/pkg/config"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type jwtAuth struct {
	key string
}

func NewJwtTokenAuth(cfg config.Config) TokenAuth {

	return &jwtAuth{
		key: cfg.JwtKey,
	}
}

var (
	ErrInvalidExpireTime  = errors.New("payload expire time is already elapsed")
	ErrInvalidToken       = errors.New("invalid token")
	ErrFailedToParseToken = errors.New("failed to parse token to claims")
	ErrExpiredToken       = errors.New("token expired")
)

type jwtClaims struct {
	UserID    string
	Email     string
	TokenID   string
	ExpiresAt time.Time
	Role      string
}

// Generate JWT token for the given payload
func (c *jwtAuth) GenerateToken(payload Payload) (string, error) {

	// check the expire time is vlaid
	if time.Since(payload.ExpireAt) > 0 {
		return "", ErrInvalidExpireTime
	}

	claims := &jwtClaims{
		TokenID:   payload.TokenID,
		UserID:    payload.UserID,
		Email:     payload.Email,
		Role:      payload.Role,
		ExpiresAt: payload.ExpireAt,
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := tkn.SignedString([]byte(c.key))

	if err != nil {
		return "", fmt.Errorf("failed to sign the token: %w", err)
	}

	return tokenString, nil
}

// Verify JWT token string and return payload
func (c *jwtAuth) VerifyToken(tokenString string) (Payload, error) {

	// parse the token string to jwt token
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return []byte(c.key), nil
	})

	if err != nil {
		// check error is token expired
		if errors.Is(err, ErrExpiredToken) {
			return Payload{}, ErrExpiredToken
		}
		return Payload{}, ErrInvalidToken
	}

	// parse the token into claims
	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return Payload{}, ErrFailedToParseToken
	}

	response := Payload{
		TokenID:  claims.TokenID,
		UserID:   claims.UserID,
		Email:    claims.Email,
		Role:     claims.Role,
		ExpireAt: claims.ExpiresAt,
	}
	return response, nil
}

// Validate claims
func (c *jwtClaims) Valid() error {
	if time.Since(c.ExpiresAt) > 0 {
		return ErrExpiredToken
	}
	return nil
}
