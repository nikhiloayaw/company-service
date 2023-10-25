package response

import "time"

type Token struct {
	AccessToken         string    `json:"accessToken"`
	AccessTokenExpireAt time.Time `json:"accessTokenExpireAt"`
}
