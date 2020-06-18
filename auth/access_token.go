package auth

import "time"

type AccessToken struct {
	Token string `json:"token"`
	Expires int64 `json:"expires"`
}

func NewAccessToken(tkn Token, exp time.Time) *AccessToken {
	return &AccessToken{
		Token: tkn.String(),
		Expires: exp.Unix(),
	}
}
