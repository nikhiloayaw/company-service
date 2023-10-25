package token

import (
	"time"
)

type TokenAuth interface {
	GenerateToken(req Payload) (string, error)
	VerifyToken(tokenString string) (Payload, error)
}

type Payload struct {
	TokenID  string
	UserID   string
	Email    string
	Role     string
	ExpireAt time.Time
}

type TokenResponse struct {
	TokenID     string
	TokenString string
}
