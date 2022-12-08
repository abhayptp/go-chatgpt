package chatgpt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tmaxmax/go-sse"
	"golang.org/x/net/proxy"
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
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")

	var lastItem []byte
	var validator sse.ResponseValidator = func(r *http.Response) error {
		if r.StatusCode != http.StatusOK {
			return fmt.Errorf("%d", r.StatusCode)
		}
		return nil
	}

	client := sse.Client{
		HTTPClient:              http.DefaultClient,
		DefaultReconnectionTime: 5 * time.Second,
		ResponseValidator:       validator,
	}
	proxyHost := os.Getenv("PROXY")
	if proxyHost != "" {
		u, err := url.Parse(proxyHost)
		if err != nil {
			return nil, err
		}
		dialer, err := proxy.FromURL(u, proxy.Direct)
		if err != nil {
			return nil, err
		}
		client.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
					return dialer.Dial(network, address)
				},
			},
		}
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
