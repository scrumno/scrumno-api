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
	// CheckAndSaveWithUploads:
	// 1) compares payload hash with stored hash
	// 2) if changed -> saves snapshot json into upload/snapshots/<key>/snapshot
	// 3) downloads images from payload (imageLinks) into upload/photos/<key>/ (best-effort)
	CheckAndSaveWithUploads(key string, payload any) (bool, error)
}

type SnapshotStore interface {
	Get(key string) (string, error)
	Set(key string, hash string) error
}
