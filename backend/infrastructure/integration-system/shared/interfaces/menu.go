package interfaces

type MenuProvider interface {
	GetMenu(OrgId string) any
}

type MenuBuilder interface {
	BuildBody(data any) (any, error)
}
