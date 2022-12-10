package chatgpt

import "time"

type Credentials struct {
	BearerToken     string
	SessionToken    string
	tokenExpiryTime time.Time
}

func NewCredentials(bearerToken string) *Credentials {
	return &Credentials{
		BearerToken: bearerToken,
	}
}
