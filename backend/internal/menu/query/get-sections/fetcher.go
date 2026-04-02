package get_sections

import (
	"context"

	"github.com/scrumno/scrumno-api/internal/menu/entity/section"
)

type Fetcher struct {
	sectionRepo section.SectionRepository
}

func NewFetcher(sectionRepo section.SectionRepository) *Fetcher {
	return &Fetcher{sectionRepo: sectionRepo}
}

func (f *Fetcher) Fetch(ctx context.Context, offset int, limit int) ([]section.Section, error) {
	return f.sectionRepo.GetAll(ctx, offset, limit)
}
