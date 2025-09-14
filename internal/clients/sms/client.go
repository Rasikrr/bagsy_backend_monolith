package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func NewClient(login, host, password string) Client {
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
	reqBody := createRequestBody(c.login, c.password, message, phones)
	return c.doRequest(ctx, http.MethodPost, "/rest/send/ ", reqBody)
}

func (c *client) doRequest(ctx context.Context, method, path string, body any) error {
	bb, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.host+path, bytes.NewReader(bb))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpc.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bad status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	return nil
}
