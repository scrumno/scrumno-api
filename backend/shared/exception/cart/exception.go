package exceptions

import (
	"errors"
)

var (
	ErrCartCreate        = errors.New("Не удалось создать корзину для пользователя")
	ErrCartFind          = errors.New("Не удалось найти корзину пользователя")
	ErrCartUpdated       = errors.New("Не удалось обновить количество товара")
	ErrCartAddProduct    = errors.New("Не удалось добавить товар в корзину")
	ErrCartRemoveProduct = errors.New("Не удалось удалить товар из корзины")
	ErrCartClear         = errors.New("Не удалось очистить корзину")
	ErrCartRecalculate   = errors.New("Не удалось пересчитать корзину")
	ErrCartInvalidQuantity = errors.New("Некорректное количество товара")
	ErrCartInvalidPrice    = errors.New("Некорректная цена товара")
)
