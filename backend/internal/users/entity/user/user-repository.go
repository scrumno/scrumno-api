package user

import "github.com/scrumno/scrumno-api/shared/interfaces/base"

type userRepository interface {
	base.BaseRepository[User]
}
