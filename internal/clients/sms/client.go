package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

	reqBody := sendRequest{
		Login:    c.login,
		Password: c.password,
		Phones:   phone,
		Message:  message,
		Format:   ResponseFormatJSON,
		Charset:  defaultCharset,
	}

	var resp sendResponse
	err := c.doWithRetry(ctx, http.MethodPost, sendSMSURL, reqBody, &resp)
	if err != nil {
		return err
	}

	if resp.HasError() {
		return resp.GetError()
	}
	return nil
}

// GetStatus проверяет статус отправленного SMS
func (c *Client) GetStatus(ctx context.Context, phone string, messageID int) (*StatusResponse, error) {
	if phone == "" {
		return nil, domainErr.ErrSMSEmptyPhone
	}
	if messageID == 0 {
		return nil, domainErr.NewInvalidInputError("invalid message ID", nil)
	}

	// Формируем URL с query параметрами
	params := url.Values{}
	params.Add("login", c.login)
	params.Add("psw", c.password)
	params.Add("phone", phone)
	params.Add("id", strconv.Itoa(messageID))
	params.Add("fmt", strconv.Itoa(ResponseFormatJSON.Int()))

	urlWithQuery := checkStatusURL + "?" + params.Encode()

	var resp StatusResponse
	if err := c.doWithRetry(ctx, http.MethodGet, urlWithQuery, nil, &resp); err != nil {
		return nil, err
	}

	if resp.HasError() {
		return nil, resp.GetError()
	}
	return &resp, nil
}

// nolint: gocognit
func (c *Client) doWithRetry(ctx context.Context, method, urlStr string, body, out any) error {
	err := retry.Do(func() error {
		var (
			reqBody []byte
			err     error
		)
		if body != nil {
			reqBody, err = json.Marshal(body)
			if err != nil {
				return domainErr.NewInternalError("failed to marshal request body", err)
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, urlStr, bytes.NewReader(reqBody))
		if err != nil {
			return domainErr.NewInternalError("failed to create request", err)
		}

		if body != nil {
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
		}

		httpResp, err := c.httpClient.Do(req)
		if err != nil {
			return domainErr.NewInternalError("failed to execute http request", err)
		}
		defer httpResp.Body.Close()

		if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
			respBody, _ := io.ReadAll(httpResp.Body)
			return domainErr.NewInternalError("unexpected api status code", nil).
				WithDetail("status_code", httpResp.StatusCode).
				WithDetail("body", string(respBody))
		}

		if out != nil {
			respBody, readErr := io.ReadAll(httpResp.Body)
			if readErr != nil {
				return domainErr.NewInternalError("failed to read response body", readErr)
			}

			if unmarshallErr := json.Unmarshal(respBody, out); unmarshallErr != nil {
				return domainErr.NewInternalError("failed to unmarshal response", unmarshallErr).
					WithDetail("body", string(respBody))
			}
		}
		return nil
	},
		retry.Attempts(3),
		retry.Delay(500*time.Millisecond),
		retry.DelayType(retry.BackOffDelay),
		retry.Context(ctx),
	)

	if err != nil {
		return domainErr.ErrSMSRequestFailed.WithError(err)
	}
	return nil
}
