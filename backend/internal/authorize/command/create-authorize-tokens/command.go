package create_authorize_tokens

import "github.com/google/uuid"

type Command struct {
	Phone               string
	UserID              uuid.UUID
	SessionID           string
	RevokePreviousToken bool
}
