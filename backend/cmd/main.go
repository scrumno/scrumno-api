package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/scrumno/scrumno-api/config"
	v1 "github.com/scrumno/scrumno-api/internal/api/v1"
	codes "github.com/scrumno/scrumno-api/internal/authorize/entity/codes"
	authorizeTokens "github.com/scrumno/scrumno-api/internal/authorize/entity/tokens"
	cartEntity "github.com/scrumno/scrumno-api/internal/cart/entity"
	category "github.com/scrumno/scrumno-api/internal/menu/entity/category"
	section "github.com/scrumno/scrumno-api/internal/menu/entity/section"
	modifier "github.com/scrumno/scrumno-api/internal/products/entity/modifier"
	"github.com/scrumno/scrumno-api/internal/products/entity/product"
	staffRole "github.com/scrumno/scrumno-api/internal/users/entity/staff-role"
	"github.com/scrumno/scrumno-api/internal/users/entity/user"
)

func main() {
	_ = godotenv.Overload(".env")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := config.Load()
	if err := config.Connect(cfg); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	config.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	if err := config.Migrate(
		&cartEntity.Cart{},
		&cartEntity.CartItem{},
		&user.User{},
		&staffRole.StaffRole{},
		&codes.AuthorizeCode{},
		&authorizeTokens.AuthorizeToken{},
		&product.Product{},
		&modifier.ProductModifier{},
		&modifier.ProductChildModifier{},
		&modifier.ProductModifierGroup{},
		&section.Section{},
		&category.Category{},
	); err != nil {
		logger.Error("миграция БД", "error", err)
		os.Exit(1)
	}

	defer func() {
		err := config.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	actions, listeners := config.DI()

	// Стартуем EventManager и регистрируем listeners один раз при запуске основного приложения.
	em := config.GetEventManager()
	config.InitEventManager(em, listeners)

	router := v1.SetupRouter(cfg, actions)
	addr := ":" + cfg.Server.Port
	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("Сервер запущен", "address", addr)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Сервер не запустился", "error", err)
		os.Exit(1)
	}
}
