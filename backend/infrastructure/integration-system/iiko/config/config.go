package config

import (
	"log/slog"
	"os"
	"strings"

	"github.com/google/uuid"
)

type Config struct {
	BaseURL          string    `env:"IIKO_BASE_URL"`
	Login            string    `env:"IIKO_LOGIN"`
	AccessToken      string    `env:"IIKO_ACCESS_TOKEN"`
	OrganizationID   uuid.UUID `env:"IIKO_ORGANIZATION_ID"`
	TerminalGroupID  uuid.UUID `env:"IIKO_TERMINAL_GROUP_ID"`
	SnapshotFilePath string    `env:"IIKO_SNAPSHOT_FILE_PATH"`
}

func Load() *Config {
	return &Config{
		BaseURL:          os.Getenv("IIKO_BASE_URL"),
		Login:            os.Getenv("IIKO_LOGIN"),
		AccessToken:      os.Getenv("IIKO_ACCESS_TOKEN"),
		OrganizationID:   parseUUIDEnvOrNil("IIKO_ORGANIZATION_ID"),
		TerminalGroupID:  parseUUIDEnvOrNil("IIKO_TERMINAL_GROUP_ID"),
		SnapshotFilePath: os.Getenv("IIKO_SNAPSHOT_FILE_PATH"),
	}
}

func parseUUIDEnvOrNil(key string) uuid.UUID {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return uuid.Nil
	}
	parsed, err := uuid.Parse(raw)
	if err != nil {
		slog.Warn("Некорректный UUID в env", "key", key, "value", raw, "error", err)
		return uuid.Nil
	}
	return parsed
}
