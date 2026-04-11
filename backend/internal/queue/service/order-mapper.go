package service

import (
	"context"
	"encoding/json"
	"math"
	"strings"

	"github.com/google/uuid"
	cartEntity "github.com/scrumno/scrumno-api/internal/cart/entity"
	"github.com/scrumno/scrumno-api/internal/products/entity/modifier"
	"github.com/scrumno/scrumno-api/internal/products/entity/product"
	"github.com/scrumno/scrumno-api/internal/queue/entity"
	"gorm.io/datatypes"
)

type QueueOrderMapper interface {
	MapCartToOrder(ctx context.Context, queueID uuid.UUID, orderID uuid.UUID, cart *cartEntity.Cart) (entity.OrdersQueueOrder, error)
}

type queueOrderMapper struct {
	productRepo  product.ProductRepository
	modifierRepo modifier.ModifierRepository
}

func NewQueueOrderMapper(productRepo product.ProductRepository, modifierRepo modifier.ModifierRepository) QueueOrderMapper {
	return &queueOrderMapper{
		productRepo:  productRepo,
		modifierRepo: modifierRepo,
	}
}

func (m *queueOrderMapper) MapCartToOrder(ctx context.Context, queueID uuid.UUID, orderID uuid.UUID, cart *cartEntity.Cart) (entity.OrdersQueueOrder, error) {
	result := entity.OrdersQueueOrder{
		ID:         orderID,
		ExternalID: orderID.String(),
		QueueID:    queueID,
	}
	if cart == nil {
		return result, nil
	}

	for _, cartItem := range cart.Items {
		baseCookMinutes := 0
		productEntity, err := m.productRepo.FindByExternalID(ctx, cartItem.ProductID.String())
		if err == nil && productEntity != nil {
			cook, err := m.productRepo.FindCookingTimeProductTableByProductID(ctx, productEntity.ID)
			if err == nil && cook != nil {
				baseCookMinutes += cook.CookingTime
			}
		}

		modifierIDs := collectModifierIDs(cartItem.Modifiers, cartItem.CommonModifiers)
		if len(modifierIDs) > 0 {
			cookingTimes, err := m.modifierRepo.GetCookingTimeModifierTableByModifierIDs(ctx, modifierIDs)
			if err == nil {
				for _, mod := range cookingTimes {
					baseCookMinutes += mod.CookingTime
				}
			}
		}

		quantity := int(math.Ceil(cartItem.Quantity))
		if quantity <= 0 {
			quantity = 1
		}

		result.Items = append(result.Items, entity.OrderItem{
			ID:               uuid.New(),
			OrderID:          orderID,
			ProductID:        cartItem.ProductID.String(),
			Quantity:         quantity,
			BaseCookMinutes:  baseCookMinutes,
			GrowthFactor:     0,
			ComplexityFactor: 1,
		})
	}

	return result, nil
}

func collectModifierIDs(modifiers *datatypes.JSON, commonModifiers *datatypes.JSON) []string {
	resultSet := map[string]struct{}{}
	collectModifierIDsFromJSON(resultSet, modifiers)
	collectModifierIDsFromJSON(resultSet, commonModifiers)

	result := make([]string, 0, len(resultSet))
	for id := range resultSet {
		result = append(result, id)
	}
	return result
}

func collectModifierIDsFromJSON(target map[string]struct{}, raw *datatypes.JSON) {
	if raw == nil || len(*raw) == 0 {
		return
	}

	var payload any
	if err := json.Unmarshal(*raw, &payload); err != nil {
		return
	}

	walkModifierPayload(payload, target)
}

func walkModifierPayload(value any, target map[string]struct{}) {
	switch typed := value.(type) {
	case []any:
		for _, item := range typed {
			walkModifierPayload(item, target)
		}
	case map[string]any:
		for key, nested := range typed {
			if isModifierIDKey(key) {
				id, ok := nested.(string)
				if ok {
					id = strings.TrimSpace(id)
					if id != "" {
						target[id] = struct{}{}
					}
				}
			}
			walkModifierPayload(nested, target)
		}
	}
}

func isModifierIDKey(key string) bool {
	switch strings.ToLower(key) {
	case "modifierid", "modifier_id", "id":
		return true
	default:
		return false
	}
}
