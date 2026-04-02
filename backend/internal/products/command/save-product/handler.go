package save_product

import (
	"context"

	"github.com/scrumno/scrumno-api/internal/products/entity/product"
	"gorm.io/gorm"
)

type Handler struct {
	productRepo product.ProductRepository
}

func NewHandler(productRepo product.ProductRepository) *Handler {
	return &Handler{
		productRepo: productRepo,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) error {
	seen := make(map[string]struct{}, len(cmd.Products))
	for _, incomingProduct := range cmd.Products {
		if incomingProduct.ID == "" {
			continue
		}
		if _, ok := seen[incomingProduct.ID]; ok {
			continue
		}
		seen[incomingProduct.ID] = struct{}{}

		entity := product.Product{
			ExternalID:              incomingProduct.ID,
			FatAmount:               incomingProduct.FatAmount,
			ProteinsAmount:          incomingProduct.ProteinsAmount,
			CarbohydratesAmount:     incomingProduct.CarbohydratesAmount,
			EnergyAmount:            incomingProduct.EnergyAmount,
			FatFullAmount:           incomingProduct.FatFullAmount,
			ProteinsFullAmount:      incomingProduct.ProteinsFullAmount,
			CarbohydratesFullAmount: incomingProduct.CarbohydratesFullAmount,
			EnergyFullAmount:        incomingProduct.EnergyFullAmount,
			Weight:                  incomingProduct.Weight,
			GroupID:                 incomingProduct.GroupID,
			ProductCategoryID:       incomingProduct.ProductCategoryID,
			Type:                    incomingProduct.Type,
			OrderItemType:           incomingProduct.OrderItemType,
			ModifierSchemaID:        incomingProduct.ModifierSchemaID,
			ModifierSchemaName:      incomingProduct.ModifierSchemaName,
			Splittable:              incomingProduct.Splittable,
			MeasureUnit:             incomingProduct.MeasureUnit,
			SizePrices:              make([]product.SizePrice, 0, len(incomingProduct.SizePrices)),
		}

		for _, incomingSizePrice := range incomingProduct.SizePrices {
			nextPrice := 0.0
			if incomingSizePrice.Price.NextPrice != nil {
				nextPrice = *incomingSizePrice.Price.NextPrice
			}

			nextDatePrice := ""
			if incomingSizePrice.Price.NextDatePrice != nil {
				nextDatePrice = *incomingSizePrice.Price.NextDatePrice
			}

			entity.SizePrices = append(entity.SizePrices, product.SizePrice{
				SizeID: incomingSizePrice.SizeID,
				Price: product.Price{
					CurrentPrice:       incomingSizePrice.Price.CurrentPrice,
					IsIncludedInMenu:   incomingSizePrice.Price.IsIncludedInMenu,
					NextPrice:          nextPrice,
					NextIncludedInMenu: incomingSizePrice.Price.NextIncludedInMenu,
					NextDatePrice:      nextDatePrice,
				},
			})
		}

		if existing, err := h.productRepo.FindByExternalID(ctx, incomingProduct.ID); err == nil && existing != nil {
			entity.ID = existing.ID
		} else if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if _, err := h.productRepo.Save(ctx, &entity); err != nil {
			return err
		}
	}

	return nil
}
