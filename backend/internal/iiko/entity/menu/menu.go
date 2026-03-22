package menu

// Menu соответствует полной структуре ответа iiko /api/1/nomenclature.
type Menu struct {
	CorrelationID     string                `json:"correlationId"`
	Groups            []ProductsGroupInfo   `json:"groups"`
	ProductCategories []ProductCategoryInfo `json:"productCategories"`
	Products          []ProductInfo         `json:"products"`
	Sizes             []Size                `json:"sizes"`
	Revision          int64                 `json:"revision"`
}

type ProductCategoryInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsDeleted bool   `json:"isDeleted"`
}

type Price struct {
	CurrentPrice       float64  `json:"currentPrice"`
	IsIncludedInMenu   bool     `json:"isIncludedInMenu"`
	NextPrice          *float64 `json:"nextPrice"`
	NextIncludedInMenu bool     `json:"nextIncludedInMenu"`
	NextDatePrice      *string  `json:"nextDatePrice"`
}

type SizePrice struct {
	SizeID *string `json:"sizeId"`
	Price  Price   `json:"price"`
}

type SimpleModifierInfo struct {
	ID                  string `json:"id"`
	DefaultAmount       *int32 `json:"defaultAmount"`
	MinAmount           int32  `json:"minAmount"`
	MaxAmount           int32  `json:"maxAmount"`
	Required            *bool  `json:"required"`
	HideIfDefaultAmount *bool  `json:"hideIfDefaultAmount"`
	Splittable          *bool  `json:"splittable"`
	FreeOfChargeAmount  *int32 `json:"freeOfChargeAmount"`
}

type ChildModifierInfo struct {
	ID                  string `json:"id"`
	DefaultAmount       *int32 `json:"defaultAmount"`
	MinAmount           int32  `json:"minAmount"`
	MaxAmount           int32  `json:"maxAmount"`
	Required            *bool  `json:"required"`
	HideIfDefaultAmount *bool  `json:"hideIfDefaultAmount"`
	Splittable          *bool  `json:"splittable"`
	FreeOfChargeAmount  *int32 `json:"freeOfChargeAmount"`
}

type GroupModifierInfo struct {
	ID                                   string              `json:"id"`
	MinAmount                            int32               `json:"minAmount"`
	MaxAmount                            int32               `json:"maxAmount"`
	Required                             *bool               `json:"required"`
	ChildModifiersHaveMinMaxRestrictions *bool               `json:"childModifiersHaveMinMaxRestrictions"`
	ChildModifiers                       []ChildModifierInfo `json:"childModifiers"`
	HideIfDefaultAmount                  *bool               `json:"hideIfDefaultAmount"`
	DefaultAmount                        *int32              `json:"defaultAmount"`
	Splittable                           *bool               `json:"splittable"`
	FreeOfChargeAmount                   *int32              `json:"freeOfChargeAmount"`
}

type Size struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Priority  *int32 `json:"priority"`
	IsDefault *bool  `json:"isDefault"`
}

type ProductInfo struct {
	FatAmount               *float64             `json:"fatAmount"`
	ProteinsAmount          *float64             `json:"proteinsAmount"`
	CarbohydratesAmount     *float64             `json:"carbohydratesAmount"`
	EnergyAmount            *float64             `json:"energyAmount"`
	FatFullAmount           *float64             `json:"fatFullAmount"`
	ProteinsFullAmount      *float64             `json:"proteinsFullAmount"`
	CarbohydratesFullAmount *float64             `json:"carbohydratesFullAmount"`
	EnergyFullAmount        *float64             `json:"energyFullAmount"`
	Weight                  *float64             `json:"weight"`
	GroupID                 *string              `json:"groupId"`
	ProductCategoryID       *string              `json:"productCategoryId"`
	Type                    *string              `json:"type"`
	OrderItemType           string               `json:"orderItemType"`
	ModifierSchemaID        *string              `json:"modifierSchemaId"`
	ModifierSchemaName      *string              `json:"modifierSchemaName"`
	Splittable              bool                 `json:"splittable"`
	MeasureUnit             string               `json:"measureUnit"`
	SizePrices              []SizePrice          `json:"sizePrices"`
	Modifiers               []SimpleModifierInfo `json:"modifiers"`
	GroupModifiers          []GroupModifierInfo  `json:"groupModifiers"`
	ImageLinks              []string             `json:"imageLinks"`
	DoNotPrintInCheque      bool                 `json:"doNotPrintInCheque"`
	ParentGroup             *string              `json:"parentGroup"`
	Order                   int32                `json:"order"`
	FullNameEnglish         *string              `json:"fullNameEnglish"`
	UseBalanceForSell       bool                 `json:"useBalanceForSell"`
	CanSetOpenPrice         bool                 `json:"canSetOpenPrice"`
	PaymentSubject          *string              `json:"paymentSubject"`
	ID                      string               `json:"id"`
	Code                    *string              `json:"code"`
	Name                    string               `json:"name"`
	Description             *string              `json:"description"`
	AdditionalInfo          *string              `json:"additionalInfo"`
	Tags                    []string             `json:"tags"`
	IsDeleted               bool                 `json:"isDeleted"`
	SeoDescription          *string              `json:"seoDescription"`
	SeoText                 *string              `json:"seoText"`
	SeoKeywords             *string              `json:"seoKeywords"`
	SeoTitle                *string              `json:"seoTitle"`
}

type ProductsGroupInfo struct {
	ImageLinks       []string `json:"imageLinks"`
	ParentGroup      *string  `json:"parentGroup"`
	Order            int32    `json:"order"`
	IsIncludedInMenu bool     `json:"isIncludedInMenu"`
	IsGroupModifier  bool     `json:"isGroupModifier"`
	ID               string   `json:"id"`
	Code             *string  `json:"code"`
	Name             string   `json:"name"`
	Description      *string  `json:"description"`
	AdditionalInfo   *string  `json:"additionalInfo"`
	Tags             []string `json:"tags"`
	IsDeleted        bool     `json:"isDeleted"`
	SeoDescription   *string  `json:"seoDescription"`
	SeoText          *string  `json:"seoText"`
	SeoKeywords      *string  `json:"seoKeywords"`
	SeoTitle         *string  `json:"seoTitle"`
}
