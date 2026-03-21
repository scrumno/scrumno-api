package update_user_profile

import "time"

type Command struct {
	Phone     string
	FullName  *string
	BirthDate *time.Time
	IsActive  *bool
	Email     *string
}
