package update_user_profile

import (
	"context"
	"errors"

	interfaces "github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/interfaces"
	userRepo "github.com/scrumno/scrumno-api/internal/authorize/entity"
	conditionsUpdateProfilePolicy "github.com/scrumno/scrumno-api/internal/users/service/conditions-update-profile"
)

type Handler struct {
	userRepo                      userRepo.RegistrationRepository
	conditionsUpdateProfilePolicy *conditionsUpdateProfilePolicy.Handler
	customerSync                  interfaces.CustomerSyncService
}

func NewHandler(userRepo userRepo.RegistrationRepository, conditionsUpdateProfilePolicy *conditionsUpdateProfilePolicy.Handler, customerSync interfaces.CustomerSyncService) *Handler {
	return &Handler{
		userRepo:                      userRepo,
		conditionsUpdateProfilePolicy: conditionsUpdateProfilePolicy,
		customerSync:                  customerSync,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) error {
	user, err := h.userRepo.FindByPhone(ctx, cmd.Phone)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("Пользователь не найден")
	}

	fields, err := h.conditionsUpdateProfilePolicy.Handle(
		ctx,
		cmd.FullName,
		cmd.BirthDate,
		cmd.IsActive,
		cmd.Email,
		user,
	)
	if err != nil {
		return err
	}
	if fields == nil {
		return nil
	}

	if h.customerSync != nil {
		updatedUser := *user
		if cmd.FullName != nil {
			updatedUser.FullName = cmd.FullName
		}
		if cmd.BirthDate != nil {
			updatedUser.BirthDate = cmd.BirthDate
		}
		if cmd.IsActive != nil {
			updatedUser.IsActive = *cmd.IsActive
		}
		if cmd.Email != nil {
			updatedUser.Email = cmd.Email
		}
		if err := h.customerSync.Sync(ctx, &updatedUser); err != nil {
			return err
		}
	}

	err = h.userRepo.UpdateFieldsByPhone(ctx, user.Phone, fields)
	if err != nil {
		return err
	}

	return nil
}
