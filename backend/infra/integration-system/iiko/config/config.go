package config

import "os"

type Config struct {
	BaseURL         string `env:"IIKO_BASE_URL"`
	OrganizationID  string `env:"IIKO_ORGANIZATION_ID"`
	TerminalGroupID string `env:"IIKO_TERMINAL_GROUP_ID"`
}

func Load() *Config {
	return &Config{
		BaseURL:         os.Getenv("IIKO_BASE_URL"),
		OrganizationID:  os.Getenv("IIKO_ORGANIZATION_ID"),
		TerminalGroupID: os.Getenv("IIKO_TERMINAL_GROUP_ID"),
	}
}
