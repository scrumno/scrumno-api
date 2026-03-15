package config

import "github.com/scrumno/scrumno-api/shared/utils"

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Sms      SmsConfig
}

type JWTConfig struct {
	SecretKey []byte
}

type SmsConfig struct {
	ApiKey string
	ApiPhoneNumber string
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
		Sms: SmsConfig{
			ApiKey: utils.GetEnv("SMS_API_KEY", ""),
			ApiPhoneNumber: utils.GetEnv("SMS_API_PHONE_NUMBER", ""),
		},
	}
}
