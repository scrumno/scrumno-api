package helpers

import (
	"encoding/json"
	"fmt"
)

func CreateBody[T any](body any) ([]byte, error) {
	_, ok := body.(T)
	if !ok {
		var expected T
		return nil, fmt.Errorf("невалидный формат body, должен быть: %T", expected)
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("не удалось закодировать JSON: %w", err)
	}

	return bodyJson, nil
}
