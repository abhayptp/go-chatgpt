package chatgpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tmaxmax/go-sse"
)

const (
	USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"
)

type Client interface {
	Send(*request) (chan *response, error)
}

type client struct {
	credentials *Credentials
}

func NewClient(credentials *Credentials) *client {
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
		return nil, fmt.Errorf("failed to create request, error: %v", err)
	}

	err = c.refreshAccessTokenIfExpired()
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authority", "chat.openai.com")
	req.Header.Set("Authorization", c.credentials.BearerToken)
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", USER_AGENT)

	var lastItem []byte
	var validator sse.ResponseValidator = func(r *http.Response) error {
		if r.StatusCode != http.StatusOK {
			return fmt.Errorf("%d", r.StatusCode)
		}
		return nil
	}

	client := sse.Client{
		HTTPClient:              http.DefaultClient,
		DefaultReconnectionTime: 8 * time.Second,
		ResponseValidator:       validator,
	}

	conn := client.NewConnection(req)
	conn.SubscribeMessages(func(event sse.Event) {
		if event.String() != "[DONE]" {
			lastItem = event.Data
		}
	})

	err = conn.Connect()
	if err != nil {
		return nil, err
	}

	if len(lastItem) > 0 {
		var respStruct response
		err = json.Unmarshal(lastItem, &respStruct)
		if err != nil {
			return nil, err
		} else {
			return &respStruct, nil
		}
	}
	return nil, errors.New("result is empty")
}

func (c *client) refreshAccessTokenIfExpired() error {
	if !c.credentials.tokenExpiryTime.IsZero() &&
		c.credentials.tokenExpiryTime.Before(time.Now()) &&
		c.credentials.BearerToken != "" {
		return nil
	}

	req, err := http.NewRequest("GET", "https://chat.openai.com/api/auth/session", nil)
	if err != nil {
		return fmt.Errorf("failed to create request, error: %v", err)
	}

	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Cookie", fmt.Sprintf("__Secure-next-auth.session-token=%s", c.credentials.SessionToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	var result sessionResponse
	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if result.Error != "" {
		if result.Error == "RefreshAccessTokenError" {
			return errors.New("session token has expired")
		}

		return errors.New(result.Error)
	}

	bearerToken := result.AccessToken
	if bearerToken == "" {
		return errors.New("unauthorized")
	}

	expiryTime, err := time.Parse(time.RFC3339, result.Expires)
	if err != nil {
		return fmt.Errorf("failed to parse expiry time: %v", err)
	}
	c.credentials.BearerToken = "Bearer " + bearerToken
	c.credentials.tokenExpiryTime = expiryTime
	return nil
}
