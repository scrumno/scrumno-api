package get_products_by_category_id

import (
	"context"

	"github.com/scrumno/scrumno-api/internal/products/entity/product"
)

type Fetcher struct {
	productRepo product.ProductRepository
}

func NewFetcher(productRepo product.ProductRepository) *Fetcher {
	return &Fetcher{productRepo: productRepo}
}

func (f *Fetcher) Fetch(ctx context.Context, categoryID string, offset int, limit int) ([]product.Product, error) {
	return f.productRepo.GetByCategoryID(ctx, categoryID, offset, limit)
}
