package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/cache/sms"
	"github.com/avast/retry-go"
)

const (
	sendSMSURL     = "https://smsc.kz/rest/send/"
	checkStatusURL = "https://smsc.kz/sys/status.php"
)

type Client interface {
	Send(ctx context.Context, phone, message string) error
}

type client struct {
	login        string
	password     string
	smsSpamCache sms.Cache
	httpc        *http.Client
}

func NewClient(login, password string, cache sms.Cache) Client {
	return &client{
		login:    login,
		password: password,
		httpc: &http.Client{
			Timeout: 10 * time.Second,
		},
		smsSpamCache: cache,
	}
}

func (c *client) Send(ctx context.Context, phone, message string) error {
	if message == "" {
		return errEmptyMessage
	}
	if len(phone) == 0 {
		return errEmptyPhones
	}

	spam, err := c.smsSpamCache.IsSpam(ctx, phone, message)
	if err != nil {
		return fmt.Errorf("sms spam cache: %w", err)
	}
	if spam {
		return errSpam
	}

	reqBody := createRequestBody(c.login, c.password, message, phone)
	resp, err := c.sendWithRetry(ctx, reqBody)
	if err != nil {
		return err
	}
	if err = resp.getError(); err != nil {
		return err
	}
	status, err := c.getStatus(ctx, phone, resp.MessageID)
	if err != nil {
		return err
	}
	if status.OneOf(errSmsStatuses...) {
		return fmt.Errorf("%w: status: %d", errCheckStatus, status)
	}
	return c.smsSpamCache.Set(ctx, phone, message)
}

func (c *client) sendWithRetry(ctx context.Context, reqBody request) (*response, error) {
	bb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	var resp *http.Response
	err = retry.Do(func() error {
		var req *http.Request
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, sendSMSURL, bytes.NewReader(bb))
		if err != nil {
			return fmt.Errorf("build request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err = c.httpc.Do(req)
		if err != nil {
			return fmt.Errorf("send request: %w", err)
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			b, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			return fmt.Errorf("bad status %d: %s", resp.StatusCode, string(b))
		}
		return nil
	}, retry.Attempts(3), retry.Delay(500*time.Millisecond))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	var sendResp response
	if err = json.Unmarshal(respBody, &sendResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w; body=%q", err, string(respBody))
	}
	return &sendResp, nil
}

func (c *client) getStatus(ctx context.Context, phone string, messageID int) (smsStatus, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, checkStatusURL, nil)
	if err != nil {
		return smsStatusNotFound, fmt.Errorf("build request: %w", err)
	}

	q := req.URL.Query()
	q.Add("login", c.login)
	q.Add("psw", c.password)
	q.Add("phone", phone)
	q.Add("charset", "utf-8")
	q.Add("fmt", strconv.Itoa(int(responseFormatJSON)))
	q.Add("id", strconv.Itoa(messageID))

	req.URL.RawQuery = q.Encode()

	resp, err := c.httpc.Do(req)
	if err != nil {
		return smsStatusNotFound, fmt.Errorf("send request: %w", err)
	}
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()
	respBody, _ := io.ReadAll(resp.Body)

	var statusResp response
	if err = json.Unmarshal(respBody, &statusResp); err != nil {
		return smsStatusNotFound, fmt.Errorf("unmarshal response: %w", err)
	}
	return statusResp.Status, nil
}
