package check_ontime_code

import (
	"context"
	"errors"
	
	codes "github.com/scrumno/scrumno-api/internal/authorize/entity/codes"
)

type Handler struct {
	codesRepo codes.SmsCodesRepository
}

func NewHandler(codesRepo codes.SmsCodesRepository) *Handler {
	return &Handler{
		codesRepo: codesRepo,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) error {
	ok, err := h.codesRepo.ValidateCode(ctx, cmd.Phone, cmd.Code, cmd.CodeType)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Код не прошел проверку")
	}

	return nil
}