package main

import (
	"fmt"
	"dfood/internal/api/routes"
	"dfood/internal/config"
	"dfood/internal/database"
	"dfood/internal/repository"
	"dfood/internal/service"
	"dfood/pkg/logger"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		logger.Error("Failed to initialize config", "error", err)
		log.Fatal("Failed to initialize config:", err)
	}

	logger.Init(cfg.Env)
	logger.Info("Starting API server", "env", cfg.Env, "port", cfg.Port)

	if err := database.InitDatabase(cfg); err != nil {
		logger.Error("Failed to initialize database", "error", err)
		log.Fatal("Failed to initialize database:", err)
	}
	logger.Info("Database initialized successfully")

	defer func() {
		if err := database.CloseDB(); err != nil {
			logger.Error("Error closing database", "error", err)
		}
	}()

	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo)

	deps := &routes.Dependencies{
		AuthService: authService,
	}

	router := routes.SetupRoutes(deps)

	logger.Info("Server listening", "port", cfg.Port)
	if err := router.Run(":" + fmt.Sprint(cfg.Port)); err != nil {
		logger.Error("Failed to start server", "error", err)
		log.Fatal("Failed to start server:", err)
	}
}
