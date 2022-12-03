package chatgpt

type responseMessage struct {
	ID         string      `json:"id"`
	Role       string      `json:"role"`
	User       interface{} `json:"user"`
	CreateTime interface{} `json:"create_time"`
	UpdateTime interface{} `json:"update_time"`
	Content    struct {
		ContentType string   `json:"content_type"`
		Parts       []string `json:"parts"`
	} `json:"content"`
	EndTurn   interface{}            `json:"end_turn"`
	Weight    float64                `json:"weight"`
	Metadata  map[string]interface{} `json:"metadata"`
	Recipient string                 `json:"recipient"`
}

type response struct {
	Message        responseMessage `json:"message"`
	ConversationID string          `json:"conversation_id"`
	Error          interface{}     `json:"error"`
}
