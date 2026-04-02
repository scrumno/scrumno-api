package file_storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Record struct {
	Hash      string    `json:"hash"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FileStore struct {
	mu       sync.RWMutex
	filePath string
}

func NewFileStore(filePath string) *FileStore {
	if filePath == "" {
		filePath = filepath.Join("upload", "snapshots", "_hashes.json")
	}

	return &FileStore{
		filePath: filePath,
	}
}

func (s *FileStore) read() (map[string]Record, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]Record), nil
		}
		return nil, err
	}

	var snapshots map[string]Record
	if err := json.Unmarshal(data, &snapshots); err != nil {
		return make(map[string]Record), err
	}

	return snapshots, nil
}

func (s *FileStore) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshots, err := s.read()
	if err != nil {
		return "", err
	}

	if record, exists := snapshots[key]; exists {
		return record.Hash, nil
	}

	return "", nil
}

func (s *FileStore) Set(key string, hash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	snapshots, err := s.read()
	if err != nil {
		return err
	}

	snapshots[key] = Record{
		Hash:      hash,
		UpdatedAt: time.Now(),
	}

	data, err := json.MarshalIndent(snapshots, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(s.filePath), 0o755); err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}
