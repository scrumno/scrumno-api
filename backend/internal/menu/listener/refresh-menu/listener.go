package refresh_menu

import refreshMenu "github.com/scrumno/scrumno-api/internal/menu/command/refresh-menu"

type Listener struct {
	handler *refreshMenu.Handler
}

func NewListener(handler *refreshMenu.Handler) *Listener {
	return &Listener{handler: handler}
}

func (l *Listener) Listen(payload any) {

}
