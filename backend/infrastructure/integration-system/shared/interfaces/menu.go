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

type SnapshotService interface {
	GenerateHash(payload any) (string, error)
	CheckAndSave(key string, payload any) (bool, error)
}

type SnapshotStore interface {
	Get(key string) (string, error)
	Set(key string, hash string) error
}
