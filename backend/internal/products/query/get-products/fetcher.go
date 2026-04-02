package get_products

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

func (f *Fetcher) Fetch(ctx context.Context, offset int, limit int) ([]product.Product, error) {
	return f.productRepo.GetAll(ctx, offset, limit)
}
