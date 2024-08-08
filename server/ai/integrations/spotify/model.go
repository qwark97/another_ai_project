package spotify

import (
	"time"
)

type Credentials struct {
	ClientID     string
	ClientSecret string
}

type oAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`

	validTo time.Time
}

func (oar oAuthResponse) valid(limit time.Time) bool {
	return !limit.After(oar.validTo)
}
