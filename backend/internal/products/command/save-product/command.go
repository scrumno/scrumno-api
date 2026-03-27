package save_product

import payloadMenuModel "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/menu/model"

type Command struct {
	Products []payloadMenuModel.MenuProduct
}
