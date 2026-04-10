package create_user

import (
	"context"

	"github.com/google/uuid"
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
	phone := cmd.Phone
	u := user.NewUser(phone)

	if h.customerSync != nil {
		get, err := h.customerSync.SyncGet(ctx, &phone, nil)
		if err != nil {
			return nil, err
		}

		if get == nil {
			setResp, err := h.customerSync.Sync(ctx, u)
			if err != nil {
				return nil, err
			}

			if setResp != nil {
				get, err = h.customerSync.SyncGet(ctx, nil, &setResp.ID)
			}

			if err != nil {
				return nil, err
			}
		}

		if get != nil {
			if len(get.WalletBalances) > 0 {
				out := make([]user.WalletBalance, 0, len(get.WalletBalances))
				for _, w := range get.WalletBalances {
					row := user.WalletBalance{
						ID:      uuid.New(),
						UserID:  u.ID,
						Name:    w.Name,
						Type:    user.WalletType(w.Type),
						Balance: w.Balance,
					}

					if w.ID != uuid.Nil {
						id := w.ID
						row.IntegrationWalletId = &id
					}

					out = append(out, row)
				}
				u.WalletBalance = out
			}
		}
	}

	if err := h.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}
