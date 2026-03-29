package model

type RefreshMenuSuccessPayload struct {
	CorrelationID     string             `json:"correlationId"`
	Groups            []StockListGroup   `json:"groups"`
	ProductCategories []MenuItemCategory `json:"productCategories"`
	Products          []MenuProduct      `json:"products"`
	// Sizes             []MenuSize         `json:"sizes"`
}

type RefreshMenuErrorPayload struct {
	CorrelationID    string  `json:"correlationId"`
	ErrorDescription string  `json:"errorDescription"`
	Error            *string `json:"error"`
}

type StockListGroup struct {
	ImageLinks       []string `json:"imageLinks"`
	ParentGroup      *string  `json:"parentGroup"`
	Order            int32    `json:"order"`
	IsIncludedInMenu bool     `json:"isIncludedInMenu"`
	IsGroupModifier  bool     `json:"isGroupModifier"`

	ID             string   `json:"id"`
	Code           *string  `json:"code"`
	Name           string   `json:"name"`
	Description    *string  `json:"description"`
	AdditionalInfo *string  `json:"additionalInfo"`
	Tags           []string `json:"tags"`
	IsDeleted      bool     `json:"isDeleted"`

	SeoDescription *string `json:"seoDescription"`
	SeoText        *string `json:"seoText"`
	SeoKeywords    *string `json:"seoKeywords"`
	SeoTitle       *string `json:"seoTitle"`
}

type MenuItemCategory struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsDeleted bool   `json:"isDeleted"`
}

type MenuProduct struct {
	FatAmount               *float64 `json:"fatAmount"`
	ProteinsAmount          *float64 `json:"proteinsAmount"`
	CarbohydratesAmount     *float64 `json:"carbohydratesAmount"`
	EnergyAmount            *float64 `json:"energyAmount"`
	FatFullAmount           *float64 `json:"fatFullAmount"`
	ProteinsFullAmount      *float64 `json:"proteinsFullAmount"`
	CarbohydratesFullAmount *float64 `json:"carbohydratesFullAmount"`
	EnergyFullAmount        *float64 `json:"energyFullAmount"`
	Weight                  *float64 `json:"weight"`

	GroupID           *string `json:"groupId"`
	ProductCategoryID *string `json:"productCategoryId"`
	Type              *string `json:"type"`
	OrderItemType     string  `json:"orderItemType"`

	ModifierSchemaID   *string `json:"modifierSchemaId"`
	ModifierSchemaName *string `json:"modifierSchemaName"`

	Splittable  bool   `json:"splittable"`
	MeasureUnit string `json:"measureUnit"`

	SizePrices     []ProductSizePrice     `json:"sizePrices"`
	Modifiers      []ProductModifier      `json:"modifiers"`
	GroupModifiers []ProductModifierGroup `json:"groupModifiers"`

	ImageLinks         []string `json:"imageLinks"`
	DoNotPrintInCheque bool     `json:"doNotPrintInCheque"`
	ParentGroup        *string  `json:"parentGroup"`
	Order              int32    `json:"order"`
	FullNameEnglish    *string  `json:"fullNameEnglish"`

	UseBalanceForSell bool    `json:"useBalanceForSell"`
	CanSetOpenPrice   bool    `json:"canSetOpenPrice"`
	PaymentSubject    *string `json:"paymentSubject"`

	ID             string   `json:"id"`
	Code           *string  `json:"code"`
	Name           string   `json:"name"`
	Description    *string  `json:"description"`
	AdditionalInfo *string  `json:"additionalInfo"`
	Tags           []string `json:"tags"`
	IsDeleted      bool     `json:"isDeleted"`

	SeoDescription *string `json:"seoDescription"`
	SeoText        *string `json:"seoText"`
	SeoKeywords    *string `json:"seoKeywords"`
	SeoTitle       *string `json:"seoTitle"`
}

type ProductSizePrice struct {
	SizeID *string      `json:"sizeId"`
	Price  ProductPrice `json:"price"`
}

type ProductPrice struct {
	CurrentPrice       float64  `json:"currentPrice"`
	IsIncludedInMenu   bool     `json:"isIncludedInMenu"`
	NextPrice          *float64 `json:"nextPrice"`
	NextIncludedInMenu bool     `json:"nextIncludedInMenu"`
	NextDatePrice      *string  `json:"nextDatePrice"`
}

type ProductModifier struct {
	ID                  string `json:"id"`
	DefaultAmount       *int32 `json:"defaultAmount"`
	MinAmount           int32  `json:"minAmount"`
	MaxAmount           int32  `json:"maxAmount"`
	Required            *bool  `json:"required"`
	HideIfDefaultAmount *bool  `json:"hideIfDefaultAmount"`
	Splittable          *bool  `json:"splittable"`
	FreeOfChargeAmount  *int32 `json:"freeOfChargeAmount"`
}

type ProductModifierGroup struct {
	ID                                   string                 `json:"id"`
	MinAmount                            int32                  `json:"minAmount"`
	MaxAmount                            int32                  `json:"maxAmount"`
	Required                             *bool                  `json:"required"`
	ChildModifiersHaveMinMaxRestrictions *bool                  `json:"childModifiersHaveMinMaxRestrictions"`
	ChildModifiers                       []ProductChildModifier `json:"childModifiers"`
}

type ProductChildModifier struct {
	ID                  string `json:"id"`
	DefaultAmount       *int32 `json:"defaultAmount"`
	MinAmount           int32  `json:"minAmount"`
	MaxAmount           int32  `json:"maxAmount"`
	Required            *bool  `json:"required"`
	HideIfDefaultAmount *bool  `json:"hideIfDefaultAmount"`
	Splittable          *bool  `json:"splittable"`
	FreeOfChargeAmount  *int32 `json:"freeOfChargeAmount"`
}

// type MenuSize struct {
// 	ID        string `json:"id"`
// 	Name      string `json:"name"`
// 	Priority  *int32 `json:"priority"`
// 	IsDefault *bool  `json:"isDefault"`
// }
