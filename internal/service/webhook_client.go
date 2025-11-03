package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/UmutcanKalkan/auto-message-dispatcher/internal/domain"
)

type WebhookClient interface {
	SendMessage(ctx context.Context, phoneNumber, content string) (*domain.WebhookResponse, error)
}

type webhookClient struct {
	url        string
	authKey    string
	client     *http.Client
	maxRetries int
	retryDelay time.Duration
}

func NewWebhookClient(url, authKey string, timeout time.Duration, maxRetries int, retryDelay time.Duration) WebhookClient {
	return &webhookClient{
		url:        url,
		authKey:    authKey,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (w *webhookClient) SendMessage(ctx context.Context, phoneNumber, content string) (*domain.WebhookResponse, error) {
	var lastErr error

	for attempt := 0; attempt <= w.maxRetries; attempt++ {
		if attempt > 0 {
			delay := w.retryDelay * time.Duration(attempt)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		resp, err := w.doSendMessage(ctx, phoneNumber, content)
		if err == nil {
			return resp, nil
		}

		lastErr = err
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

func (w *webhookClient) doSendMessage(ctx context.Context, phoneNumber, content string) (*domain.WebhookResponse, error) {
	payload := domain.WebhookRequest{
		To:      phoneNumber,
		Content: content,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-ins-auth-key", w.authKey)

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var webhookResp domain.WebhookResponse
	if err := json.Unmarshal(body, &webhookResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &webhookResp, nil
}
