package section

type Section struct {
	ImageLinks       []string `json:"imageLinks" gorm:"type:jsonb;serializer:json"`
	ParentGroup      *string  `json:"parentGroup"`
	Order            int32    `json:"order"`
	IsIncludedInMenu bool     `json:"isIncludedInMenu"`
	IsGroupModifier  bool     `json:"isGroupModifier"`

	ID             string   `json:"id"`
	Code           *string  `json:"code"`
	Name           string   `json:"name"`
	Description    *string  `json:"description"`
	AdditionalInfo *string  `json:"additionalInfo"`
	Tags           []string `json:"tags" gorm:"type:jsonb;serializer:json"`
	IsDeleted      bool     `json:"isDeleted"`

	SeoDescription *string `json:"seoDescription"`
	SeoText        *string `json:"seoText"`
	SeoKeywords    *string `json:"seoKeywords"`
	SeoTitle       *string `json:"seoTitle"`
}
