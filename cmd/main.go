package main

import (
	"dfood/internal/api/routes"
	"dfood/internal/config"
	"dfood/internal/database"
	"dfood/internal/repository"
	"dfood/internal/service"
	"dfood/internal/utils"
	"dfood/pkg/logger"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	cfg, err := config.New()
	if err != nil {
		logger.Error("Failed to initialize config", "error", err)
		log.Fatal("Failed to initialize config:", err)
	}

	logger.Init(cfg.Env)
	logger.Info("Starting API server", "env", cfg.Env, "port", cfg.Port)

	// Initialize JWT secret
	if err := utils.InitJWT(cfg.JWT.Secret); err != nil {
		logger.Error("Failed to initialize JWT", "error", err)
		log.Fatal("Failed to initialize JWT:", err)
	}
	logger.Info("JWT initialized successfully")

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

	// Initialize repositories
	userRepo := repository.NewUserRepository()
	restaurantRepo := repository.NewRestaurantRepository()
	foodRepo := repository.NewFoodRepository()
	orderRepo := repository.NewOrderRepository()
	addressRepo := repository.NewAddressRepository()

	// Initialize services
	emailService := service.NewEmailService(service.EmailConfig{
		APIKey:    cfg.SendGrid.APIKey,
		FromEmail: cfg.SendGrid.FromEmail,
		FromName:  cfg.SendGrid.FromName,
	})
	authService := service.NewAuthService(userRepo, emailService)
	userService := service.NewUserService(userRepo)
	restaurantService := service.NewRestaurantService(restaurantRepo, foodRepo)
	foodService := service.NewFoodService(foodRepo)
	orderService := service.NewOrderService(orderRepo, userRepo, restaurantRepo, foodRepo)
	addressService := service.NewAddressService(addressRepo, userRepo)

	deps := &routes.Dependencies{
		AuthService:       authService,
		UserService:       userService,
		RestaurantService: restaurantService,
		FoodService:       foodService,
		OrderService:      orderService,
		AddressService:    addressService,
		UserRepository:    userRepo,
	}

	router := routes.SetupRoutes(deps)

	logger.Info("Server listening", "port", cfg.Port)
	if err := router.Run(":" + fmt.Sprint(cfg.Port)); err != nil {
		logger.Error("Failed to start server", "error", err)
		log.Fatal("Failed to start server:", err)
	}
}
