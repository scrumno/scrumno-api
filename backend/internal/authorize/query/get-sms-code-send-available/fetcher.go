package get_sms_code_send_available

import (
	"context"

	codes "github.com/scrumno/scrumno-api/internal/authorize/entity/codes"
)

type Fetcher struct {
	codesRepository codes.SmsCodesRepository
}

func NewFetcher(
	codesRepository codes.SmsCodesRepository,
) *Fetcher {
	return &Fetcher{
		codesRepository: codesRepository,
	}
}

func (h *Fetcher) Fetch(ctx context.Context, phone string) (bool, error) {
	_, err := h.codesRepository.ValidateCodeByCreatedAt(ctx, phone)
	if err != nil {
		return false, err
	}

	return true, nil
}
