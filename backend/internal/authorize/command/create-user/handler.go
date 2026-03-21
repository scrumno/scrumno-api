package create_user

import (
	"context"

	user "github.com/scrumno/scrumno-api/internal/authorize/entity"
)

type Handler struct {
	userRepo user.RegistrationRepository
}

func NewHandler(userRepo user.RegistrationRepository) *Handler {
	return &Handler{
		userRepo: userRepo,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (*user.User, error) {
	user := user.NewUser(cmd.Phone)

	if err := h.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
