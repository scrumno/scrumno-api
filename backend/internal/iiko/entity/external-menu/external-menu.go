package external_menu

import "encoding/json"

// MenusDataResponse соответствует ответу iiko /api/2/menu.
type MenusDataResponse struct {
	CorrelationID   string          `json:"correlationId"`
	ExternalMenus   []ExternalMenu  `json:"externalMenus"`
	PriceCategories []PriceCategory `json:"priceCategories"`
}

type ExternalMenu struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PriceCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// MenuRequest соответствует контракту iiko /api/2/menu/by_id.
type MenuRequest struct {
	ExternalMenuID  string   `json:"externalMenuId"`
	OrganizationIDs []string `json:"organizationIds"`
	PriceCategoryID *string  `json:"priceCategoryId,omitempty"`
	Version         *int32   `json:"version,omitempty"`
	Language        *string  `json:"language,omitempty"`
	AsyncMode       *bool    `json:"asyncMode,omitempty"`
	StartRevision   *int64   `json:"startRevision,omitempty"`
}

// ByIDResponse хранит полное тело ответа /api/2/menu/by_id.
// В iiko есть несколько форматов (V2/V3/V4), поэтому сохраняем payload целиком.
type ByIDResponse struct {
	FormatVersion *int            `json:"formatVersion,omitempty"`
	Raw           json.RawMessage `json:"raw"`
	Payload       map[string]any  `json:"payload"`
}
