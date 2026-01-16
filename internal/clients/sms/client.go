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

	"github.com/avast/retry-go"
	"github.com/cockroachdb/errors"
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

// SendMessage отправляет SMS на указанный номер телефона
func (c *Client) SendMessage(ctx context.Context, phone, message string) error {
	if message == "" {
		return ErrEmptyMessage
	}
	if phone == "" {
		return ErrEmptyPhone
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
		return nil, ErrEmptyPhone
	}
	if messageID == 0 {
		return nil, ErrInvalidMsgID
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
				return errors.Wrap(ErrMarshalFailed, err.Error())
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, urlStr, bytes.NewReader(reqBody))
		if err != nil {
			return errors.Wrap(ErrCreateRequestFailed, err.Error())
		}

		if body != nil {
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
		}

		httpResp, err := c.httpClient.Do(req)
		if err != nil {
			return errors.Wrap(ErrHTTPRequestFailed, err.Error())
		}
		defer httpResp.Body.Close()

		if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
			respBody, _ := io.ReadAll(httpResp.Body)
			return errors.Wrapf(ErrUnexpectedStatus, "status_code=%d, body=%s", httpResp.StatusCode, string(respBody))
		}

		if out != nil {
			respBody, readErr := io.ReadAll(httpResp.Body)
			if readErr != nil {
				return errors.Wrap(ErrReadBodyFailed, readErr.Error())
			}

			if unmarshallErr := json.Unmarshal(respBody, out); unmarshallErr != nil {
				return errors.Wrapf(ErrUnmarshalFailed, "%s, body=%s", unmarshallErr.Error(), string(respBody))
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
		return errors.Wrap(ErrRequestFailed, err.Error())
	}
	return nil
}
