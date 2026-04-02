package save_menu

import (
	"github.com/scrumno/scrumno-api/internal/menu/entity/category"
	"github.com/scrumno/scrumno-api/internal/menu/entity/section"
)

type Command struct {
	Sections   []section.Section
	Categories []category.Category
}
