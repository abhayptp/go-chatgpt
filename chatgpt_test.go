package chatgpt

import (
	"fmt"
	"testing"
)

// TestSendMessage tests the send method of the client
func TestSendMessage(t *testing.T) {
	// Prepare test data
	chatgpt := NewChatGpt(NewClient(&Credentials{
		BearerToken:  "Bearer <Bearer-Token>",
		SessionToken: "<Session-Token>",
	}))
	mockRequest := "hello"

	// Run test
	res, err := chatgpt.SendMessage(mockRequest)
	if err != nil {
		t.Errorf("error sending request, %v", err)
	}

	// Handle response
	fmt.Println(res)
}
