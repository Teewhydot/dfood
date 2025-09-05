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
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		logger.Error("Failed to initialize config", "error", err)

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

	server := &http.Server{
		Addr:         ":" + fmt.Sprint(cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Server listening", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", "error", err)
			log.Fatal("Failed to start server:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	if err := server.Close(); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited")
}
