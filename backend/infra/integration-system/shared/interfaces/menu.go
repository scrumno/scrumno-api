package interfaces

type MenuProvider interface {
	GetMenu(OrgId string) any
}
