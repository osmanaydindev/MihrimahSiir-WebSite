// Package mail, Resend HTTP API üzerinden bildirim maili gönderir.
// SMTP yerine HTTP API kullanılıyor: VPS'lerde giden 587 portu bazen
// kapalı oluyor ve HTTP tarafında hata mesajları çok daha net.
package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	resendEndpoint = "https://api.resend.com/emails"
	maxBodyBytes   = 64 << 10
)

// resendRequest, Resend gövde şeması. Alanlar json.Marshal'dan geçtiği
// için header enjeksiyonu için string birleştirme asla kullanılmaz.
type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
	ReplyTo string   `json:"reply_to,omitempty"`
}

type resendResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConnsPerHost:   2,
		IdleConnTimeout:       60 * time.Second,
		ResponseHeaderTimeout: 8 * time.Second,
	},
}

// sendViaResend, tek bir maili gönderir ve sağlayıcı id'sini döner.
func sendViaResend(ctx context.Context, to, subject, html string) (string, error) {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("RESEND_API_KEY tanımlı değil")
	}

	payload := resendRequest{
		From: mailFrom(),
		To:   []string{to},
		// Subject downstream'de gerçek bir header'a yazılıyor.
		Subject: sanitizeHeaderValue(subject),
		HTML:    html,
		ReplyTo: os.Getenv("MAIL_REPLY_TO"),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Resend 2 istek/sn sınırlıyor; 429'da bir kez tekrar dene.
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		id, retryable, err := doResendRequest(ctx, apiKey, body)
		if err == nil {
			return id, nil
		}
		lastErr = err
		if !retryable {
			break
		}
		select {
		case <-time.After(600 * time.Millisecond):
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
	return "", lastErr
}

func doResendRequest(ctx context.Context, apiKey string, body []byte) (string, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendEndpoint, bytes.NewReader(body))
	if err != nil {
		return "", false, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", true, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))

	var parsed resendResponse
	_ = json.Unmarshal(respBody, &parsed)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return parsed.ID, false, nil
	}

	message := parsed.Message
	if message == "" {
		message = strings.TrimSpace(string(respBody))
	}
	retryable := resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500
	return "", retryable, fmt.Errorf("resend %d: %s", resp.StatusCode, message)
}

func mailFrom() string {
	if from := os.Getenv("MAIL_FROM"); from != "" {
		return from
	}
	return "MihrimahSiir <noreply@resend.mihrimahsiir.com>"
}

// sanitizeHeaderValue, CR/LF temizler. HTTP API'de ham SMTP header'ı
// yok ama Subject downstream'de header'a yazıldığı için yine de temizlenir.
func sanitizeHeaderValue(v string) string {
	v = strings.ReplaceAll(v, "\r", " ")
	v = strings.ReplaceAll(v, "\n", " ")
	return strings.TrimSpace(v)
}
