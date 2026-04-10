package save_modifier

import (
	"context"

	"github.com/scrumno/scrumno-api/internal/products/entity/modifier"
	"gorm.io/gorm"
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
	seenGroups := make(map[string]struct{}, len(cmd.Groups))
	for _, group := range cmd.Groups {
		if group.ID == "" {
			continue
		}
		if _, ok := seenGroups[group.ID]; ok {
			continue
		}
		seenGroups[group.ID] = struct{}{}
		if err := h.modifierRepo.SaveProductModifierGroup(ctx, &modifier.ProductModifierGroup{
			ID:        group.ID,
			MinAmount: group.MinAmount,
			MaxAmount: group.MaxAmount,
		}); err != nil {
			return err
		}
	}

	seenChild := make(map[string]struct{}, len(cmd.ChildModifiers))
	for _, childModifier := range cmd.ChildModifiers {
		if childModifier.ID == "" {
			continue
		}
		if _, ok := seenChild[childModifier.ID]; ok {
			continue
		}
		seenChild[childModifier.ID] = struct{}{}
		if err := h.modifierRepo.SaveProductChildModifier(ctx, &modifier.ProductChildModifier{
			ID:            childModifier.ID,
			DefaultAmount: childModifier.DefaultAmount,
			MinAmount:     childModifier.MinAmount,
			MaxAmount:     childModifier.MaxAmount,
		}); err != nil {
			return err
		}
	}

	seenModifiers := make(map[string]struct{}, len(cmd.Modifiers))
	for _, m := range cmd.Modifiers {
		if m.ID == "" {
			continue
		}
		if _, ok := seenModifiers[m.ID]; ok {
			continue
		}
		seenModifiers[m.ID] = struct{}{}
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
		if _, err := h.modifierRepo.FindCookingTimeModifierTableByExternalID(ctx, m.ID); err == nil {
		} else if err == gorm.ErrRecordNotFound {
			if err := h.modifierRepo.UpdateCookingTimeModifierTable(ctx, &modifier.CookingTimeModifierTable{
				ModifierID:  m.ID,
				CookingTime: 0,
			}); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
