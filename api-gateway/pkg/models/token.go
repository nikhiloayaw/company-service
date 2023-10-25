package models

import "time"

type TokenPayload struct {
	TokenID  string
	UserID   string
	Email    string
	Role     string
	ExpireAt time.Time
}
