package create_authorize_code

import (
	"context"
	"fmt"

	authorizeCodes "github.com/scrumno/scrumno-api/internal/authorize/entity/codes"
	createUniqueCode "github.com/scrumno/scrumno-api/internal/authorize/service/create-unique-code"
)

type Handler struct {
	codesRepo           authorizeCodes.SmsCodesRepository
	createUniqueCodeSvc *createUniqueCode.CreateUniqueCodeService
}

func NewHandler(
	codesRepo authorizeCodes.SmsCodesRepository,
	createUniqueCodeSvc *createUniqueCode.CreateUniqueCodeService,
) *Handler {
	return &Handler{
		codesRepo:           codesRepo,
		createUniqueCodeSvc: createUniqueCodeSvc,
	}
}

func (h *Handler) Handle(ctx context.Context, phone string, codeType authorizeCodes.CodesType) (*authorizeCodes.AuthorizeCode, error) {
	uniqueCode, err := h.createUniqueCodeSvc.CreateUniqueCode(ctx)
	if err != nil {
		return nil, err
	}

	authorizeCode := authorizeCodes.NewAuthorizeCode(phone, uniqueCode, codeType)
	fmt.Printf("authorizeCode: %+v\n", authorizeCode)
	if err = h.codesRepo.Create(ctx, authorizeCode); err != nil {
		return nil, err
	}

	return authorizeCode, nil
}
