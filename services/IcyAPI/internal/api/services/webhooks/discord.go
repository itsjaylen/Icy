package webhooks

import (
	"bytes"
	"encoding/json"
	"errors"
	logger "itsjaylen/IcyLogger"
	"net/http"
	"time"
)

type DiscordWebhook struct {
	Content string `json:"content"`
}

const (
	maxRetries     = 3
	retryDelay     = 5 * time.Second
	requestTimeout = 10 * time.Second
)

func SendDiscordWebhook(webhookURL, message string) error {
	payload := DiscordWebhook{
		Content: message,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: requestTimeout}
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(data))
		if err != nil {
			lastErr = err
			logger.Error.Printf("[Webhook] Attempt %d failed: %v", attempt, err)
			time.Sleep(retryDelay)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			logger.Error.Printf("[Webhook] Attempt %d failed: %v", attempt, err)
			time.Sleep(retryDelay)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			logger.Info.Println("[Webhook] Discord webhook sent successfully!")
			return nil
		}

		lastErr = errors.New(resp.Status)
		logger.Error.Printf("[Webhook] Attempt %d failed: %s", attempt, resp.Status)
		time.Sleep(retryDelay)
	}

	logger.Error.Printf("[Webhook] Last error: %v", lastErr)
	return lastErr
}
