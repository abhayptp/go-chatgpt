package chatgpt

import (
	"fmt"
	"testing"
)

// TestSend tests the send method of the client
func TestSendMessage(t *testing.T) {
	// Prepare test data
	credentials := &credentials{BearerToken: "Bearer <Bearer-Token>"}
	client := NewChatGpt(NewClient(credentials))
	mockRequest := "hello"

	// Run test
	res, err := client.SendMessage(mockRequest)
	if err != nil {
		t.Errorf("error sending request, %v", err)
	}

	// Handle response
	fmt.Println(res)
}
