package config

import "github.com/scrumno/scrumno-api/shared/utils"

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Sms      SmsConfig
	Iiko     IikoConfig
}

type JWTConfig struct {
	SecretKey []byte
	AccessTokenTtl time.Duration
	RefreshTokenTtl time.Duration
	AccessSecret string
	RefreshSecret string
}

type SmsConfig struct {
	ApiKey string
	ApiPhoneNumber string
}

type IikoConfig struct {
	BaseURL        string
	Login          string
	Password       string
	OrganizationID string
	TerminalID     string
}

func Load() *Config {
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
			SecretKey: utils.GetEnv("JWT_SECRET", ""),
			AccessSecret:    utils.GetEnv("JWT_ACCESS_SECRET", ""),
			RefreshSecret:   utils.GetEnv("JWT_REFRESH_SECRET", ""),
			AccessTokenTtl:  utils.GetEnv("JWT_ACCESS_TOKEN_TTL", "900"),
			RefreshTokenTtl: utils.GetEnv("JWT_REFRESH_TOKEN_TTL", "604800"),
		},
		Sms: SmsConfig{
			ApiKey: utils.GetEnv("SMS_API_KEY", ""),
			ApiPhoneNumber: utils.GetEnv("SMS_API_PHONE_NUMBER", ""),
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
