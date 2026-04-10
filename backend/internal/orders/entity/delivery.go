package entity

type DeliveryPoint struct {
	Coordinates           *Coordinates `json:"coordinates,omitempty"`
	Address               any          `json:"address,omitempty"`
	ExternalCartographyID string       `json:"externalCartographyId,omitempty"`
	Comment               string       `json:"comment,omitempty"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
