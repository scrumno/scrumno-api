package exceptions

import (
	"errors"
)

var (
	ErrCartCreate    	= errors.New("Не удалось создать корзину для пользователя")
	ErrCartFind			= errors.New("Не удалось найти корзину пользователя")
	ErrCartUpdated      = errors.New("Не удалось обновить количество товара")
	ErrCartAddProduct	= errors.New("Не удалось добавить товар в корзину")
)

