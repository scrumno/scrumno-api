package config

import "github.com/scrumno/scrumno-api/shared/utils"

// Config — настройки домена интеграции с iiko (загружаются из окружения рядом с кодом интеграции).
type Config struct {
	BaseURL        string
	Login          string
	Password       string
	OrganizationID string
	TerminalID     string
}

func Load() Config {
	return Config{
		BaseURL:        utils.GetEnv("IIKO_BASE_URL", "https://990-418-833.iiko.it/resto/"),
		Login:          utils.GetEnv("IIKO_LOGIN", ""),
		Password:       utils.GetEnv("IIKO_PASSWORD", ""),
		OrganizationID: utils.GetEnv("IIKO_ORGANIZATION_ID", ""),
		TerminalID:     utils.GetEnv("IIKO_TERMINAL_ID", ""),
	}
}
