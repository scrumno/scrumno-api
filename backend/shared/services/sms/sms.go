package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type SMSRequest struct {
	Number      string `json:"number"`
	Destination string `json:"destination"`
	Text        string `json:"text"`
}

type SmsService struct {
	ApiKey         string
	ApiPhoneNumber string
}

type Config struct {
	ApiKey         string
	ApiPhoneNumber string
}

func NewSmsService(cfg Config) *SmsService {
	return &SmsService{
		ApiKey:         cfg.ApiKey,
		ApiPhoneNumber: cfg.ApiPhoneNumber,
	}
}

func (s *SmsService) SendSmsMessage(ctx context.Context, phone string, message string) error {
	requestData := SMSRequest{
		Number:      s.ApiPhoneNumber, // заменяем на наш купленный номер
		Destination: phone,
		Text:        message,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.exolve.ru/messaging/v1/SendSMS", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", s.ApiKey) // заменяем на наш API ключ
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send SMS: received status code %d", resp.StatusCode)
	}
	return nil
}
