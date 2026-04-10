package config

import (
	"os"

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
		OrganizationID:   uuid.MustParse(os.Getenv("IIKO_ORGANIZATION_ID")),
		TerminalGroupID:  uuid.MustParse(os.Getenv("IIKO_TERMINAL_GROUP_ID")),
		SnapshotFilePath: os.Getenv("IIKO_SNAPSHOT_FILE_PATH"),
	}
}
