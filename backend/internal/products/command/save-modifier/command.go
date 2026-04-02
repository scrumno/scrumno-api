package save_modifier

import payloadMenuModel "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/menu/model"

type Command struct {
	ChildModifiers []payloadMenuModel.ProductChildModifier
	Modifiers      []payloadMenuModel.ProductModifier
	Groups         []payloadMenuModel.ProductModifierGroup
}
