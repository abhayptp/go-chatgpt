package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client interface {
	Send(*request) (chan *response, error)
}

type client struct {
	credentials *credentials
}

func NewClient(credentials *credentials) Client {
	return &client{
		credentials: credentials,
	}
}

func (c *client) Send(r *request) (chan *response, error) {
	reqBytes, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("ill request, error: %v", err)
	}

	req, err := http.NewRequest("POST", "https://chat.openai.com/backend-api/conversation", bytes.NewReader(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("couldn't create request, error: %v", err)
	}

	req.Header.Set("Authority", "chat.openai.com")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't send request, error: %v", err)
	}
	defer resp.Body.Close()

	res := make(chan *response)

	buf := make([]byte, 1024)
	go func() {
		for {
			_, err := resp.Body.Read(buf)
			if err != nil {
				break
			}
			var resp response
			err = json.Unmarshal(buf, &res)
			if err != nil {
				res <- &response{
					Error: "couldn't unmarshal the response",
				}
			}
			res <- &resp
		}
		close(res)
	}()

	return res, nil
}
