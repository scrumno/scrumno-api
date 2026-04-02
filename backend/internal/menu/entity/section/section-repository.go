package section

import (
	factory "github.com/scrumno/scrumno-api/shared/factories/gorm"
	"github.com/scrumno/scrumno-api/shared/interfaces/base"
	"gorm.io/gorm"
)

type SectionRepository interface {
	base.BaseRepository[Section]
}

type sectionRepository struct {
	base.BaseRepository[Section]
}

func NewSectionRepository(db *gorm.DB) SectionRepository {
	return &sectionRepository{
		BaseRepository: factory.NewGormRepository[Section](db),
	}
}
