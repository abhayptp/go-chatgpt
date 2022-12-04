# ChatGPT API

Unofficial golang client for ChatGPT. Reverse Engineered from [chat.openai.com](chat.openai.com)

## Usage

1. Get the bearer token from the browser.


2. Install the package.
```bash
go get github.com/abhayptp/go-chatgpt.git
```

3. Pass the bearer token while initializing client.

```go
package main

import (
	"fmt"

	"github.com/abhayptp/go-chatgpt"
)

func main() {

	// Initialize. Copy bearer-token from browser developer tools.
	c := chatgpt.NewChatGpt(chatgpt.NewClient(chatgpt.NewCredentials("Bearer <Bearer-Token>")))

	// Send message
	res, err := c.SendMessage("hello")
	if err != nil {
		// Handle err
	}

	// Handle response
	fmt.Println(res)
}
```