package get_categories

import (
	"context"

	"github.com/scrumno/scrumno-api/internal/menu/entity/category"
)

type Fetcher struct {
	categoryRepo category.CategoryRepository
}

func NewFetcher(categoryRepo category.CategoryRepository) *Fetcher {
	return &Fetcher{categoryRepo: categoryRepo}
}

func (f *Fetcher) Fetch(ctx context.Context, offset int, limit int) ([]category.Category, error) {
	return f.categoryRepo.GetAll(ctx, offset, limit)
}
