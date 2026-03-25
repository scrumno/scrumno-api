package update_user_profile

import (
	"context"
	"errors"
	
	userRepo "github.com/scrumno/scrumno-api/internal/authorize/entity"
	conditionsUpdateProfilePolicy "github.com/scrumno/scrumno-api/internal/users/service/conditions-update-profile"
)

type Handler struct {
	userRepo userRepo.RegistrationRepository
	conditionsUpdateProfilePolicy *conditionsUpdateProfilePolicy.Handler
}

func NewHandler(userRepo userRepo.RegistrationRepository, conditionsUpdateProfilePolicy *conditionsUpdateProfilePolicy.Handler) *Handler {
	return &Handler{
		userRepo: userRepo,
		conditionsUpdateProfilePolicy: conditionsUpdateProfilePolicy,
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

	err = h.userRepo.UpdateFieldsByPhone(ctx, user.Phone, fields)
	if err != nil {
		return err
	}

	return nil
}