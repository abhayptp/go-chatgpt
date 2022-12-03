package chatgpt

type content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type requestMessage struct {
	ID      string  `json:"id"`
	Role    string  `json:"role"`
	Content content `json:"content"`
}

type request struct {
	Action          string           `json:"action"`
	Messages        []requestMessage `json:"messages"`
	ParentMessageID string           `json:"parent_message_id"`
	Model           string           `json:"model"`
}

func newRequest(action string, m []requestMessage, par string, model string) *request {
	return &request{
		Action:          action,
		Messages:        m,
		ParentMessageID: par,
		Model:           model,
	}
}
