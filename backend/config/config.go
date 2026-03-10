package config

import "github.com/scrumno/scrumno-api/shared/utils"

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Iiko     IikoConfig
}

type JWTConfig struct {
	SecretKey []byte
}

type IikoConfig struct {
	BaseURL        string
	Login          string
	Password       string
	OrganizationID string
	TerminalID     string
}

func Load() *Config {
	secret := utils.GetEnv("JWT_SECRET", "default-dev-secret-change-in-production")
	secretKey := []byte(secret)

	return &Config{
		Server: ServerConfig{
			Port: utils.GetEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:         utils.GetEnv("DATABASE_HOST", "localhost"),
			Port:         utils.GetEnv("DATABASE_PORT", "5432"),
			Username:     utils.GetEnv("DATABASE_USERNAME", ""),
			Password:     utils.GetEnv("DATABASE_PASSWORD", ""),
			DatabaseName: utils.GetEnv("DATABASE_NAME", ""),
			SSLMode:      utils.GetEnv("DATABASE_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey: secretKey,
		},
		Iiko: IikoConfig{
			BaseURL:        utils.GetEnv("IIKO_BASE_URL", "https://990-418-833.iiko.it/resto/"),
			Login:          utils.GetEnv("IIKO_LOGIN", ""),
			Password:       utils.GetEnv("IIKO_PASSWORD", ""),
			OrganizationID: utils.GetEnv("IIKO_ORGANIZATION_ID", ""),
			TerminalID:     utils.GetEnv("IIKO_TERMINAL_ID", ""),
		},
	}
}
