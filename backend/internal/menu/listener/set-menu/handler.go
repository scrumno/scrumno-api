package setmenu

import (
	"context"

	"github.com/scrumno/scrumno-api/internal/menu/categories/entity/category"
	"github.com/scrumno/scrumno-api/internal/menu/products/entity/product"
	getMenuFromProvider "github.com/scrumno/scrumno-api/internal/menu/query/get-menu-from-provider"
)

type Handler struct {
	GetMenuFromProvider *getMenuFromProvider.Fetcher
	ProductRepository   *product.ProductRepository
	CategoryRepository  *category.CategoryRepository
}

func NewHandler(
	getMenuFromProvider *getMenuFromProvider.Fetcher,
	productRepository *product.ProductRepository,
	categoryRepository *category.CategoryRepository,
) *Handler {
	return &Handler{
		GetMenuFromProvider: getMenuFromProvider,
		ProductRepository:   productRepository,
		CategoryRepository:  categoryRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, query getMenuFromProvider.Query) error {
	return nil
}
