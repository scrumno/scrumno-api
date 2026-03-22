package access

// AccessToken доменная сущность результата авторизации iiko.
// Храним только поле, нужное остальному приложению.
type AccessToken struct {
	Token string
}
