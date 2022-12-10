# ChatGPT client (unofficial)

Unofficial golang client for ChatGPT. Reverse Engineered from [chat.openai.com](https://chat.openai.com)

## Usage

1. Install the package.
```bash
go get github.com/abhayptp/go-chatgpt
```

2. Get bearer token from the browser.

To avoid needing to refresh bearer token every hour, you can also copy "__Secure-next-auth.session-token" key from cookie and pass it in Credentials in Step 3. 


<img src="https://user-images.githubusercontent.com/22256898/205469104-d99b6a6a-18d2-4fea-9a58-6936d3be6479.png" width=80%>


3. Pass the bearer token while initializing client.

```go
package main

import (
	"fmt"

	"github.com/abhayptp/go-chatgpt"
)

func main() {

	// Initialize. Copy bearer-token and session-token from browser developer tools.
	c := chatgpt.NewChatGpt(chatgpt.NewClient(&chatgpt.Credentials{
		BearerToken: "Bearer <bearer-token>",
		SessionToken: "<session-token>",
		}))

	// Send message
	res, err := c.SendMessage("hello")
	if err != nil {
		// Handle err
	}

	// Handle response
	fmt.Println(res)
}
```
