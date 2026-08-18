package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type setWebhookResponse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

func SetupWebhook(botToken, appDomain, secret string) error {
	webhookURL := fmt.Sprintf("%s/api/v1/telegram/webhook", appDomain)
	telegramAPIURL := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", botToken)

	data := url.Values{}
	data.Set("url", webhookURL)
	if secret != "" {
		data.Set("secret_token", secret)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.PostForm(telegramAPIURL, data)
	if err != nil {
		return fmt.Errorf("failed to send setWebhook request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read setWebhook response body: %w", err)
	}

	var tgResp setWebhookResponse
	if err := json.Unmarshal(body, &tgResp); err != nil {
		return fmt.Errorf("failed to parse setWebhook response: %w", err)
	}

	if !tgResp.Ok {
		return fmt.Errorf("telegram setWebhook failed: %s", tgResp.Description)
	}

	log.Printf("[TELEGRAM] Webhook successfully registered: %s", webhookURL)
	return nil
}
