package model

import (
	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/cart/entity"
)

type BaseSuccessResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error,omitempty"`
}

type SuccessResponse struct {
	IsSuccess bool      `json:"isSuccess"`
	Error     string    `json:"error,omitempty"`
	CartID    uuid.UUID `json:"cartId"`
}

type CartSuccessResponse struct {
	IsSuccess bool         `json:"isSuccess"`
	Cart      *entity.Cart `json:"cart"`
}

type BaseErrorResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error,omitempty"`
	ErrorCode string `json:"errorCode,omitempty"`
}
