package config

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
	SSLMode      string
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=%s password=%s client_encoding=UTF8",
		c.Host, c.Port, c.Username, c.DatabaseName, c.SSLMode, c.Password,
	)
}

var DB *gorm.DB

func Connect(cfg *Config) error {
	var err error
	dsn := cfg.Database.DSN()
	slog.Info("Подключение к БД", "dsn", dsn)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	slog.Info("Подключение к БД установлено")
	return nil
}

// migrateUsersPhone adds NOT NULL column "phone" to "users" when the table
// already has rows. Existing rows get a unique placeholder so the constraint holds.
func migrateUsersPhone() error {
	var exists bool
	err := DB.Raw(
		`SELECT EXISTS (
			SELECT 1 FROM information_schema.columns
			WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'phone'
		)`).Scan(&exists).Error
	if err != nil || exists {
		return err
	}
	// Add column with unique default so existing rows satisfy NOT NULL and unique index
	if err := DB.Exec(`ALTER TABLE users ADD COLUMN phone text NOT NULL DEFAULT gen_random_uuid()::text`).Error; err != nil {
		return err
	}
	if err := DB.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_phone ON users(phone)`).Error; err != nil {
		return err
	}
	if err := DB.Exec(`ALTER TABLE users ALTER COLUMN phone DROP DEFAULT`).Error; err != nil {
		return err
	}
	slog.Info("Миграция users.phone выполнена (существующие строки получили placeholder)")
	return nil
}

func Migrate(models ...interface{}) error {
	slog.Info("Миграция базы данных")

	if err := migrateUsersPhone(); err != nil {
		return err
	}

	if err := DB.AutoMigrate(models...); err != nil {
		return err
	}

	slog.Info("Миграция завершена")
	return nil
}

func Close() error {
	slog.Info("Соединение с БД закрывается")

	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	slog.Info("Соединение с БД закрыто")
	return nil
}
