package entity

type OrderItem struct {
	Type             string            `json:"type"`
	Amount           float64           `json:"amount"`
	ProductSizeID    string            `json:"productSizeId,omitempty"`
	ComboInformation *ComboInformation `json:"comboInformation,omitempty"`
	Comment          string            `json:"comment,omitempty"`
}

type ComboInformation struct {
	ComboID        string `json:"comboId"`
	ComboSourceID  string `json:"comboSourceId"`
	ComboGroupID   string `json:"comboGroupId"`
	ComboGroupName string `json:"comboGroupName,omitempty"`
}

type OrderCombo struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Amount    float64 `json:"amount"`
	Price     float64 `json:"price"`
	SourceID  string  `json:"sourceId"`
	ProgramID string  `json:"programId,omitempty"`
	SizeID    string  `json:"sizeId,omitempty"`
}
