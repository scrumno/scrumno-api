package product

type Product struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	FatAmount               *float64    `json:"fatAmount" gorm:"type:numeric(10,3)"`
	ProteinsAmount          *float64    `json:"proteinsAmount" gorm:"type:numeric(10,3)"`
	CarbohydratesAmount     *float64    `json:"carbohydratesAmount" gorm:"type:numeric(10,3)"`
	EnergyAmount            *float64    `json:"energyAmount" gorm:"type:numeric(10,3)"`
	FatFullAmount           *float64    `json:"fatFullAmount" gorm:"type:numeric(10,3)"`
	ProteinsFullAmount      *float64    `json:"proteinsFullAmount" gorm:"type:numeric(10,3)"`
	CarbohydratesFullAmount *float64    `json:"carbohydratesFullAmount" gorm:"type:numeric(10,3)"`
	EnergyFullAmount        *float64    `json:"energyFullAmount" gorm:"type:numeric(10,3)"`
	Weight                  *float64    `json:"weight" gorm:"type:numeric(10,3)"`
	GroupID                 *string     `json:"groupId" gorm:"size:128;index"`
	ProductCategoryID       *string     `json:"productCategoryId" gorm:"size:128;index"`
	Type                    *string     `json:"type" gorm:"size:128;index"`
	OrderItemType           string      `json:"orderItemType" gorm:"size:128;index"`
	ModifierSchemaID        *string     `json:"modifierSchemaId" gorm:"size:128;index"`
	ModifierSchemaName      *string     `json:"modifierSchemaName" gorm:"size:255"`
	Splittable              bool        `json:"splittable" gorm:"default:false"`
	MeasureUnit             string      `json:"measureUnit" gorm:"size:32"`
	SizePrices              []SizePrice `json:"sizePrices" gorm:"type:jsonb;serializer:json"`
}

type SizePrice struct {
	SizeID *string `json:"sizeId" gorm:"size:128;index"`
	Price  Price   `json:"price" gorm:"embedded;embeddedPrefix:price_"`
}

type Price struct {
	CurrentPrice       float64 `json:"currentPrice" gorm:"type:numeric(12,2);default:0"`
	IsIncludedInMenu   bool    `json:"isIncludedInMenu" gorm:"default:true"`
	NextPrice          float64 `json:"nextPrice" gorm:"type:numeric(12,2);default:0"`
	NextIncludedInMenu bool    `json:"nextIncludedInMenu" gorm:"default:true"`
	NextDatePrice      string  `json:"nextDatePrice" gorm:"size:128;index"`
}
