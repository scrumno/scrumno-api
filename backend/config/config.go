package config

import (
	"strconv"
	"time"

	"github.com/scrumno/scrumno-api/shared/utils"
)

type Config struct {
	Server            ServerConfig
	Database          DatabaseConfig
	JWT               JWTConfig
	Sms               SmsConfig
	IntegrationSystem IntegrationSystemConfig
}

type JWTConfig struct {
	SecretKey       string
	AccessTokenTtl  time.Duration
	RefreshTokenTtl time.Duration
	AccessSecret    string
	RefreshSecret   string
}

type SmsConfig struct {
	ApiKey         string
	ApiPhoneNumber string
}

type IntegrationSystemConfig struct {
	IntegrationSystem string
}

func Load() *Config {
	accessTokenTtl, _ := strconv.Atoi(utils.GetEnv("JWT_ACCESS_TOKEN_TTL", "15"))
	refreshTokenTtl, _ := strconv.Atoi(utils.GetEnv("JWT_REFRESH_TOKEN_TTL", "10080"))

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
			SecretKey:       utils.GetEnv("JWT_SECRET", ""),
			AccessSecret:    utils.GetEnv("JWT_ACCESS_SECRET", ""),
			RefreshSecret:   utils.GetEnv("JWT_REFRESH_SECRET", ""),
			AccessTokenTtl:  time.Duration(accessTokenTtl) * time.Second,
			RefreshTokenTtl: time.Duration(refreshTokenTtl) * time.Second,
		},
		Sms: SmsConfig{
			ApiKey:         utils.GetEnv("SMS_API_KEY", ""),
			ApiPhoneNumber: utils.GetEnv("SMS_API_PHONE_NUMBER", ""),
		},
		IntegrationSystem: IntegrationSystemConfig{
			IntegrationSystem: utils.GetEnv("INTEGRATION_SYSTEM", ""),
		},
	}
}
