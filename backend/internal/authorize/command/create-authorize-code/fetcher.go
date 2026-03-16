package create_authorize_code

import (
	"context"
	"fmt"

	codes "github.com/scrumno/scrumno-api/internal/authorize/entity/codes"
	createUniqueCode "github.com/scrumno/scrumno-api/internal/authorize/service/create-unique-code"
)

type Handler struct {
	codesRepo codes.SmsCodesRepository
	createUniqueCodeSvc *createUniqueCode.CreateUniqueCodeService
}

func NewHandler(
	codesRepo codes.SmsCodesRepository,
	createUniqueCodeSvc *createUniqueCode.CreateUniqueCodeService,
) *Handler {
	return &Handler{
		codesRepo: codesRepo,
		createUniqueCodeSvc: createUniqueCodeSvc,
	}
}

func (h *Handler) Handle(ctx context.Context, phone string, codeType codes.CodesType) (*codes.AuthorizeCode, error) {
	uniqueCode, err := h.createUniqueCodeSvc.CreateUniqueCode(ctx)
	if err != nil {
		return nil, err
	}

	authorizeCode := codes.NewAuthorizeCode(phone, uniqueCode, codeType)
	fmt.Printf("authorizeCode: %+v\n", authorizeCode)
	if err = h.codesRepo.Create(ctx, authorizeCode); err != nil {
		return nil, err
	}

	return authorizeCode, nil
}