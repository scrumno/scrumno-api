package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/scrumno/scrumno-api/config"
	iikoConfig "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/config"
	iikoMenuService "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/menu/service"
	menuInterfaces "github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/interfaces"
	v1 "github.com/scrumno/scrumno-api/internal/api/v1"
	codes "github.com/scrumno/scrumno-api/internal/authorize/entity/codes"
	authorizeTokens "github.com/scrumno/scrumno-api/internal/authorize/entity/tokens"
	staffRole "github.com/scrumno/scrumno-api/internal/users/entity/staff-role"
	"github.com/scrumno/scrumno-api/internal/users/entity/user"
)

func main() {
	_ = godotenv.Overload(".env.local")
	_ = godotenv.Overload(".env")
	_ = godotenv.Overload("backend/.env.local")
	_ = godotenv.Overload("backend/.env")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	cfg := config.Load()
	if err := config.Connect(cfg); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	config.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	if err := config.Migrate(
		&user.User{},
		&staffRole.StaffRole{},
		&codes.AuthorizeCode{},
		&authorizeTokens.AuthorizeToken{},
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

	// Стартуем EventManager и регистрируем listeners один раз при запуске основного приложения.
	em := config.GetEventManager()
	iikoCfg := iikoConfig.Load()
	var menuProvider menuInterfaces.MenuProvider = iikoMenuService.NewMenuProvider(iikoCfg)
	config.InitEventManager(em, menuProvider)

	actions := config.DI()

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
