package chatgpt

type credentials struct {
	BearerToken string
}

func NewCredentials(bearerToken string) *credentials {
	return &credentials{
		BearerToken: bearerToken,
	}
}
