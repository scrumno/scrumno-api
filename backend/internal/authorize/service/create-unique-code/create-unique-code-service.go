package service_code_generator

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
)

type CreateUniqueCodeService struct {}

func NewCreateUniqueCodeService() *CreateUniqueCodeService {
	return &CreateUniqueCodeService{}
}

func (s *CreateUniqueCodeService) CreateUniqueCode(ctx context.Context) (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%04d", n.Int64()), nil
}