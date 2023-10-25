package response

import "time"

type SignUp struct {
	UserID string `json:"user_id"`
}

type SignIn struct {
	AccessToken string    `json:"access_token"`
	ExpireAt    time.Time `json:"expire_at"`
}
