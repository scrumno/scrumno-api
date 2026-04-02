package save_menu

import (
	"context"

	"github.com/scrumno/scrumno-api/internal/menu/entity/category"
	"github.com/scrumno/scrumno-api/internal/menu/entity/section"
)

type Handler struct {
	sectionRepo  section.SectionRepository
	categoryRepo category.CategoryRepository
}

func NewHandler(sectionRepo section.SectionRepository, categoryRepo category.CategoryRepository) *Handler {
	return &Handler{
		sectionRepo:  sectionRepo,
		categoryRepo: categoryRepo,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) error {
	for _, section := range cmd.Sections {
		if _, err := h.sectionRepo.Save(ctx, &section); err != nil {
			return err
		}
	}
	for _, category := range cmd.Categories {
		if _, err := h.categoryRepo.Save(ctx, &category); err != nil {
			return err
		}
	}
	return nil
}
