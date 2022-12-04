package chatgpt

import (
	"github.com/google/uuid"
)

type chatGpt struct {
	client *client

	conversationId string
	lastMessageId  string
}

func NewChatGpt(c *client) *chatGpt {
	return &chatGpt{
		client:        c,
		lastMessageId: uuid.Must(uuid.NewRandom()).String(),
	}
}

func (c *chatGpt) SendMessage(m string) (string, error) {
	req := newRequest("next",
		[]requestMessage{
			{
				ID:   uuid.Must(uuid.NewRandom()).String(),
				Role: "user",
				Content: content{
					ContentType: "text",
					Parts:       []string{m},
				},
			},
		},
		c.lastMessageId,
		"text-davinci-002-render",
	)
	resp, err := c.client.Send(req)
	if err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", nil
	}

	if c.conversationId == "" {
		c.conversationId = resp.ConversationID
	}

	return resp.Message.Content.Parts[0], nil
}
