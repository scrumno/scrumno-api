package interfaces

type MenuProvider interface {
	GetMenu() (any, error)
}

type MenuBuilder interface {
	BuildBody(data any) (any, error)
}

type GetMenuHandler interface {
	Handle() any
}
