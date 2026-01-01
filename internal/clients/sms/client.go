package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/avast/retry-go"
)

const (
	sendSMSURL     = "https://smsc.kz/rest/send/"
	checkStatusURL = "https://smsc.kz/sys/status.php"
	defaultCharset = "utf-8"
)

type Client struct {
	login      string
	password   string
	httpClient *http.Client
}

// NewClient создает новый экземпляр клиента SMSC.KZ
func NewClient(login, password string) *Client {
	return &Client{
		login:    login,
		password: password,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send отправляет SMS на указанный номер телефона
func (c *Client) Send(ctx context.Context, phone, message string) error {
	if message == "" {
		return domainErr.ErrSMSEmptyMessage
	}
	if phone == "" {
		return domainErr.ErrSMSEmptyPhone
	}

	req := sendRequest{
		Login:    c.login,
		Password: c.password,
		Phones:   phone,
		Message:  message,
		Format:   ResponseFormatJSON,
		Charset:  defaultCharset,
	}

	resp, err := c.sendWithRetry(ctx, req)
	if err != nil {
		return err
	}

	if resp.HasError() {
		return resp.GetError()
	}

	return nil
}

// SendWithOptions отправляет SMS с дополнительными параметрами
func (c *Client) sendWithRetry(ctx context.Context, req sendRequest) (*sendResponse, error) {
	var resp *sendResponse
	err := retry.Do(func() error {
		r, err := c.sendRequest(ctx, req)
		if err != nil {
			return err
		}
		resp = r
		return nil
	},
		retry.Attempts(3),
		retry.Delay(500*time.Millisecond),
		retry.DelayType(retry.BackOffDelay),
		retry.Context(ctx),
	)

	if err != nil {
		return nil, domainErr.ErrSMSSendFailed.WithError(err)
	}

	return resp, nil
}

// sendRequest выполняет HTTP запрос на отправку SMS
func (c *Client) sendRequest(ctx context.Context, reqBody sendRequest) (*sendResponse, error) {
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, domainErr.NewInternalError("failed to marshal SMS request", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sendSMSURL, bytes.NewReader(payload))
	if err != nil {
		return nil, domainErr.NewInternalError("failed to create SMS request", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, domainErr.NewInternalError("failed to send SMS request", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		body, _ := io.ReadAll(httpResp.Body)
		return nil, domainErr.NewInternalError("unexpected SMS API status code", nil).
			WithDetail("status_code", httpResp.StatusCode).
			WithDetail("body", string(body))
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, domainErr.NewInternalError("failed to read SMS response body", err)
	}

	var resp sendResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, domainErr.NewInternalError("failed to unmarshal SMS response", err).
			WithDetail("body", string(body))
	}

	return &resp, nil
}

// GetStatus проверяет статус отправленного SMS
func (c *Client) GetStatus(ctx context.Context, phone string, messageID int) (*statusResponse, error) {
	if phone == "" {
		return nil, domainErr.ErrSMSEmptyPhone
	}
	if messageID == 0 {
		return nil, domainErr.NewInvalidInputError("invalid message ID", nil)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, checkStatusURL, nil)
	if err != nil {
		return nil, domainErr.NewInternalError("failed to create status request", err)
	}

	q := req.URL.Query()
	q.Add("login", c.login)
	q.Add("psw", c.password)
	q.Add("phone", phone)
	q.Add("id", fmt.Sprintf("%d", messageID))
	q.Add("fmt", "3") // JSON format
	q.Add("charset", defaultCharset)
	req.URL.RawQuery = q.Encode()

	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, domainErr.NewInternalError("failed to send status request", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		body, _ := io.ReadAll(httpResp.Body)
		return nil, domainErr.NewInternalError("unexpected status API response code", nil).
			WithDetail("status_code", httpResp.StatusCode).
			WithDetail("body", string(body))
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, domainErr.NewInternalError("failed to read status response body", err)
	}

	var resp statusResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, domainErr.NewInternalError("failed to unmarshal status response", err).
			WithDetail("body", string(body))
	}

	if resp.HasError() {
		return nil, resp.GetError()
	}

	return &resp, nil
}
