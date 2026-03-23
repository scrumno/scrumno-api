package get_sms_code

import (
	"context"
	"errors"

	"github.com/scrumno/scrumno-api/shared/services/sms"
)

type Fetcher struct {
	smsService *sms.SmsService
}

func NewFetcher(
	smsService *sms.SmsService,
) *Fetcher {
	return &Fetcher{
		smsService: smsService,
	}
}

func (h *Fetcher) Fetch(ctx context.Context, phone, message string) (bool, error) {
	if h.smsService != nil && h.smsService.ApiKey != "" {
		if err := h.smsService.SendSmsMessage(ctx, phone, message); err != nil {
			return false, errors.New("Не удалось отправить код")
		}
	}

	return true, nil
}
