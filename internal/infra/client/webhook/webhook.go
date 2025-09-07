package webhook

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
	url    string
}

func NewWebhookClient(url string, timeout time.Duration) *Client {
	transport := &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		DisableCompression:    false,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	return &Client{
		client: client,
		url:    url,
	}
}

type WebhookRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type WebhookResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

func (c *Client) SendMessage(ctx context.Context, to string, content string) (*WebhookResponse, error) {
	reqBody := WebhookRequest{
		To:      to,
		Content: content,
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		return nil, fmt.Errorf("status_code=%d, body=%s", resp.StatusCode, string(body))
	}

	var webhookResp WebhookResponse
	if err := json.Unmarshal(body, &webhookResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &webhookResp, nil
}
