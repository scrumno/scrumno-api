package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func (s *Service) CheckAndSaveWithUploads(key string, payload any) (bool, error) {
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

	safeKey := sanitizeKeySegment(key)

	// 1) Сохраняем snapshot json.
	if err := saveSnapshotFile(safeKey, payload); err != nil {
		return false, fmt.Errorf("failed to save snapshot: %w", err)
	}

	// 2) Загружаем фото (best-effort).
	if err := downloadMenuImages(safeKey, payload); err != nil {
		// Ошибка не должна ломать refresh-menu: слепок всё равно уже сохранён.
		slog.Error("refresh-menu: failed to download some images", "error", err)
	}

	// 3) Только после успешной записи snapshot обновляем hash.
	if err := s.store.Set(key, newHash); err != nil {
		return true, fmt.Errorf("данные изменились, но не удалось сохранить новый спепок: %w", err)
	}

	return true, nil
}

func sanitizeKeySegment(key string) string {
	// На всякий случай защищаемся от слэшей/переходов директорий.
	key = strings.ReplaceAll(key, "/", "_")
	key = strings.ReplaceAll(key, "\\", "_")
	return key
}

func saveSnapshotFile(key string, payload any) error {
	dateSegment := time.Now().UTC().Format("20060102")
	snapDir := filepath.Join("upload", "snapshots", key, dateSegment)
	if err := os.MkdirAll(snapDir, 0o755); err != nil {
		return err
	}

	snapPath := filepath.Join(snapDir, "snapshot.json")
	snapBytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(snapPath, snapBytes, 0o644)
}

func downloadMenuImages(key string, payload any) error {
	links, err := extractImageLinksFromMenuPayload(payload)
	if err != nil {
		return err
	}
	if len(links) == 0 {
		return nil
	}

	photoDir := filepath.Join("upload", "photos", key)
	if err := os.MkdirAll(photoDir, 0o755); err != nil {
		return err
	}

	client := &http.Client{Timeout: 20 * time.Second}
	seen := make(map[string]struct{}, len(links))

	var downloadErr error
	for _, link := range links {
		if link == "" {
			continue
		}
		if _, ok := seen[link]; ok {
			continue
		}
		seen[link] = struct{}{}

		if err := downloadFileIfNotExists(client, link, photoDir); err != nil {
			// продолжаем; в конце вернем первую ошибку, чтобы handler мог залогировать при желании
			if downloadErr == nil {
				downloadErr = err
			}
			slog.Error("refresh-menu: photo download failed", "url", link, "error", err)
		}
	}

	return downloadErr
}

func extractImageLinksFromMenuPayload(payload any) ([]string, error) {
	var envelope struct {
		Groups []struct {
			ImageLinks []string `json:"imageLinks"`
		} `json:"groups"`
		Products []struct {
			ImageLinks []string `json:"imageLinks"`
		} `json:"products"`
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &envelope); err != nil {
		return nil, err
	}

	links := make([]string, 0, len(envelope.Groups)+len(envelope.Products))
	for _, g := range envelope.Groups {
		links = append(links, g.ImageLinks...)
	}
	for _, p := range envelope.Products {
		links = append(links, p.ImageLinks...)
	}
	return links, nil
}

func downloadFileIfNotExists(client *http.Client, fileURL, dir string) error {
	u, err := url.Parse(fileURL)
	if err != nil {
		return err
	}

	ext := path.Ext(u.Path)
	if ext == "" || len(ext) > 5 {
		ext = ".img"
	}

	hash := sha256.Sum256([]byte(fileURL))
	filename := hex.EncodeToString(hash[:]) + ext
	dstPath := filepath.Join(dir, filename)

	if _, err := os.Stat(dstPath); err == nil {
		return nil // уже есть
	}

	req, err := http.NewRequest(http.MethodGet, fileURL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("http %d: %s", resp.StatusCode, string(b))
	}

	f, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

