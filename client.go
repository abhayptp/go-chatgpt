package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client interface {
	Send(*request) (chan *response, error)
}

type client struct {
	credentials *credentials
}

func NewClient(credentials *credentials) *client {
	return &client{
		credentials: credentials,
	}
}

func (c *client) Send(r *request) (res *response, err error) {
	reqBytes, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("ill request, error: %v", err)
	}

	req, err := http.NewRequest("POST", "https://chat.openai.com/backend-api/conversation", bytes.NewReader(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("couldn't create request, error: %v", err)
	}

	req.Header.Set("Authority", "chat.openai.com")
	req.Header.Set("Authorization", c.credentials.BearerToken)
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't send request, error: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non 200 response, status code: %v", resp.StatusCode)
	}
	defer resp.Body.Close()

	buf := make([]byte, 102400)
	fmt.Println(resp.StatusCode)

	// TODO: Remove aritifical timeout. Keep listening to events until "DONE" event is received.
	time.Sleep(5 * time.Second)

	_, err = resp.Body.Read(buf)
	if err != nil {
		return nil, err
	}

	lines := bytes.Split(buf, []byte("\n"))
	for _, line := range lines[:len(lines)-4] {
		if len(line) < 6 {
			continue
		}
		line = line[6:]
		var respStruct response
		err = json.Unmarshal(line, &respStruct)
		if err != nil {
			return res, fmt.Errorf("couldn't unmarshal response, error: %v", err)
		}
		res = &respStruct
		fmt.Println(res)
	}

	return res, nil
}
