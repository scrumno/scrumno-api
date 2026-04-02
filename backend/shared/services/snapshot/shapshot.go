package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/interfaces"
)

type Service struct {
	store interfaces.SnapshotStore
}

func NewSnapshotService(store interfaces.SnapshotStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GenerateHash(payload any) (string, error) {
	dataBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("ошибка маршалинга структуры для хеширования: %w", err)
	}

	hash := sha256.Sum256(dataBytes)
	return hex.EncodeToString(hash[:]), nil
}

func (s *Service) CheckAndSave(key string, payload any) (bool, error) {
	newHash, err := s.GenerateHash(payload)
	if err != nil {
		return false, err
	}

	oldHash, err := s.store.Get(key)
	if err != nil {
		return false, fmt.Errorf("не удалось получить старый спепок из хранилища: %w", err)
	}

	if oldHash == newHash {
		return false, nil
	}

	if err := s.store.Set(key, newHash); err != nil {
		return true, fmt.Errorf("данные изменились, но не удалось сохранить новый спепок: %w", err)
	}

	return true, nil
}
