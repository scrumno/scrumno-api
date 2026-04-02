package menu

import (
	"net/http"
	"reflect"

	"github.com/scrumno/scrumno-api/internal/menu/entity/category"
	"github.com/scrumno/scrumno-api/internal/menu/entity/section"
	getCategories "github.com/scrumno/scrumno-api/internal/menu/query/get-categories"
	getSections "github.com/scrumno/scrumno-api/internal/menu/query/get-sections"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	"github.com/scrumno/scrumno-api/internal/products/entity/product"
	getProducts "github.com/scrumno/scrumno-api/internal/products/query/get-products"
)

type GetMenuAction struct {
	GetCategoriesFetcher *getCategories.Fetcher
	GetSectionsFetcher   *getSections.Fetcher
	GetProductsFetcher   *getProducts.Fetcher
}

func NewGetMenuAction(
	getCategoriesFetcher *getCategories.Fetcher,
	getSectionsFetcher *getSections.Fetcher,
	getProductsFetcher *getProducts.Fetcher,
) *GetMenuAction {
	return &GetMenuAction{
		GetCategoriesFetcher: getCategoriesFetcher,
		GetSectionsFetcher:   getSectionsFetcher,
		GetProductsFetcher:   getProductsFetcher,
	}
}

type GetMenuRequest struct {
	Offset int `json:"offset" example:"0"`
	Limit  int `json:"limit" example:"10"`
}

type GetMenuResponse struct {
	Categories []category.Category `json:"categories"`
	Sections   []section.Section   `json:"sections"`
	Products   []product.Product   `json:"products"`
}

func (a *GetMenuAction) Action(w http.ResponseWriter, r *http.Request) {
	var req GetMenuRequest
	err := utils.DecodeJSONBody(r, &req)
	if err != nil {
		utils.JSONResponse(w, map[string]any{
			"isSuccess": false,
			"error":     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	categories, err := a.GetCategoriesFetcher.Fetch(r.Context(), req.Offset, req.Limit)
	if err != nil {
		utils.JSONResponse(w, map[string]any{
			"isSuccess": false,
			"error":     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	sections, err := a.GetSectionsFetcher.Fetch(r.Context(), req.Offset, req.Limit)
	if err != nil {
		utils.JSONResponse(w, map[string]any{
			"isSuccess": false,
			"error":     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	products, err := a.GetProductsFetcher.Fetch(r.Context(), req.Offset, req.Limit)
	if err != nil {
		utils.JSONResponse(w, map[string]any{
			"isSuccess": false,
			"error":     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	utils.JSONResponse(w, GetMenuResponse{
		Categories: categories,
		Sections:   sections,
		Products:   products,
	}, http.StatusOK)
}

func (a *GetMenuAction) GetInputType() reflect.Type {
	return reflect.TypeOf(GetMenuRequest{})
}
