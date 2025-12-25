package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

type EmailService struct {
	webhookURL string
}

type EmailRequest struct {
	ToEmail          string `json:"toEmail"`
	VerificationCode string `json:"verificationCode"`
}

func NewEmailService(webhookURL string) *EmailService {
	return &EmailService{
		webhookURL: webhookURL,
	}
}

func (es *EmailService) SendVerificationEmail(email string) (string, error) {
	// Generate 6-digit code
	code := fmt.Sprintf("%06d", rand.Intn(900000)+100000)

	// Prepare request payload
	payload := EmailRequest{
		ToEmail:          email,
		VerificationCode: code,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal email request: %v", err)
	}

	// Send POST request to Google Apps Script webhook
	resp, err := http.Post(es.webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to send email: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("email service returned status: %d", resp.StatusCode)
	}

	return code, nil
}
