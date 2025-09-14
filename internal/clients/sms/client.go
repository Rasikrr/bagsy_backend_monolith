package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client interface {
	Send(ctx context.Context, message string, phones ...string) error
}

type client struct {
	host     string
	login    string
	password string
	httpc    *http.Client
}

func NewClient(host, login, password string) Client {
	return &client{
		host:     host,
		login:    login,
		password: password,
		httpc: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *client) Send(ctx context.Context, message string, phones ...string) error {
	if message == "" {
		return errEmptyMessage
	}
	if len(phones) == 0 {
		return errEmptyPhones
	}
	reqBody := createRequestBody(c.login, c.password, message, phones)

	resp, err := c.send(ctx, reqBody)
	if err != nil {
		return err
	}
	if err := resp.getError(); err != nil {
		return err
	}
	return nil
}

func (c *client) send(ctx context.Context, reqBody request) (*sendResponse, error) {
	bb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.host+"/rest/send/", bytes.NewReader(bb))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, _ := io.ReadAll(resp.Body)

	var sendResp sendResponse

	if err := json.Unmarshal(respBody, &sendResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &sendResp, nil
}
