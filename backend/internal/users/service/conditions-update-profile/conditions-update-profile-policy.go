package conditions_update_profile_policy

import (
	"context"
	"time"

	user "github.com/scrumno/scrumno-api/internal/authorize/entity"
)

type Handler struct {}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(
	ctx context.Context, 
	fullName *string, 
	birthDate *time.Time, 
	isActive *bool, 
	email *string, 
	user *user.User,
) (map[string]any, error) {
	
	fields := map[string]any{}
	
	if fullName != nil {
		fields["full_name"] = *fullName
	}
	if birthDate != nil && user.BirthDate == nil {
		fields["birth_date"] = *birthDate
	}
	if isActive != nil && *isActive != user.IsActive {
		fields["is_active"] = *isActive
	}
	if email != nil {
		fields["email"] = *email
	}

	if len(fields) == 0 {
		return nil, nil
	}
	
	return fields, nil
}