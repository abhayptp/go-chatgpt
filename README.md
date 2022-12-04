# ChatGPT client (unofficial)

Unofficial golang client for ChatGPT. Reverse Engineered from [chat.openai.com](https://chat.openai.com)

## Usage

1. Install the package.
```bash
go get github.com/abhayptp/go-chatgpt.git
```

2. Get bearer token from the browser.

<img src="https://user-images.githubusercontent.com/22256898/205469104-d99b6a6a-18d2-4fea-9a58-6936d3be6479.png" width=80%>


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
