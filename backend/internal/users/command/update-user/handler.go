package update_user

import "github.com/scrumno/scrumno-api/internal/users/entity/user"

type Handler struct {
	repository user.UserRepositoryInterface
}

func NewHandler(repository user.UserRepositoryInterface) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) Handle(cmd Command) (UserDTO, error) {
	updates := map[string]interface{}{}

	if cmd.Phone != nil {
		updates["phone"] = *cmd.Phone
	}
	if cmd.FullName != nil {
		updates["full_name"] = *cmd.FullName
	}
	if cmd.BirthDate != nil {
		updates["birth_date"] = *cmd.BirthDate
	}

	updated, err := h.repository.Update(cmd.ID, updates)
	if err != nil {
		return UserDTO{}, err
	}

	return UserDTO{
		ID:        updated.ID,
		Phone:     updated.Phone,
		FullName:  updated.FullName,
		BirthDate: updated.BirthDate,
		IsActive:  updated.IsActive,
		CreatedAt: updated.CreatedAt,
	}, nil
}
