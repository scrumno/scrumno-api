package create_user

import (
	"context"

	interfaces "github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/interfaces"
	user "github.com/scrumno/scrumno-api/internal/authorize/entity"
)

type Handler struct {
	userRepo     user.RegistrationRepository
	customerSync interfaces.CustomerSyncService
}

func NewHandler(userRepo user.RegistrationRepository, customerSync interfaces.CustomerSyncService) *Handler {
	return &Handler{
		userRepo:     userRepo,
		customerSync: customerSync,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (*user.User, error) {
	user := user.NewUser(cmd.Phone)

	if h.customerSync != nil {
		outerUser := h.customerSync.SyncGet(ctx, user)

		if outerUser == nil {
			if err := h.customerSync.Sync(ctx, user); err != nil {
				return nil, err
			}
		}

	}

	if err := h.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
