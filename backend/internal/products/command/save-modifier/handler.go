package save_modifier

import (
	"context"

	"github.com/scrumno/scrumno-api/internal/products/entity/modifier"
)

type Handler struct {
	modifierRepo modifier.ModifierRepository
}

func NewHandler(modifierRepo modifier.ModifierRepository) *Handler {
	return &Handler{
		modifierRepo: modifierRepo,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) error {
	for _, group := range cmd.Groups {
		if err := h.modifierRepo.SaveProductModifierGroup(ctx, &modifier.ProductModifierGroup{
			ID:        group.ID,
			MinAmount: group.MinAmount,
			MaxAmount: group.MaxAmount,
		}); err != nil {
			return err
		}
	}

	for _, childModifier := range cmd.ChildModifiers {
		if err := h.modifierRepo.SaveProductChildModifier(ctx, &modifier.ProductChildModifier{
			ID:            childModifier.ID,
			DefaultAmount: childModifier.DefaultAmount,
			MinAmount:     childModifier.MinAmount,
			MaxAmount:     childModifier.MaxAmount,
		}); err != nil {
			return err
		}
	}

	for _, m := range cmd.Modifiers {
		if err := h.modifierRepo.SaveProductModifier(ctx, &modifier.ProductModifier{
			ID:                  m.ID,
			DefaultAmount:       m.DefaultAmount,
			MinAmount:           m.MinAmount,
			MaxAmount:           m.MaxAmount,
			Required:            m.Required,
			HideIfDefaultAmount: m.HideIfDefaultAmount,
			Splittable:          m.Splittable,
			FreeOfChargeAmount:  m.FreeOfChargeAmount,
		}); err != nil {
			return err
		}
	}

	return nil
}
