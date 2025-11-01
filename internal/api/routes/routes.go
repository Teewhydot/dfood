package routes

import (
	"time"

	"dfood/internal/api/handlers"
	"dfood/internal/api/middleware"
	"dfood/internal/repository"
	"dfood/internal/service"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	AuthService       service.AuthService
	UserService       service.UserService
	RestaurantService service.RestaurantService
	FoodService       service.FoodService
	OrderService      service.OrderService
	AddressService    service.AddressService
	UserRepository    repository.UserRepository
}

func SetupRoutes(deps *Dependencies) *gin.Engine {
	router := gin.New()

	// Global Middleware
	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())
	router.Use(middleware.RateLimitMiddleware(10, time.Minute)) // 10 requests per minute per IP

	// Initialize WebSocket Service
	wsService := service.NewWebSocketService()
	deps.UserService.SetWebSocketService(wsService)

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(deps.AuthService)
	userHandler := handlers.NewUserHandler(deps.UserService)
	restaurantHandler := handlers.NewRestaurantHandler(deps.RestaurantService)
	foodHandler := handlers.NewFoodHandler(deps.FoodService)
	orderHandler := handlers.NewOrderHandler(deps.OrderService)
	addressHandler := handlers.NewAddressHandler(deps.AddressService)
	wsHandler := handlers.NewWebSocketHandler(wsService, deps.UserService)

	// API v1 Routes
	v1 := router.Group("/api/v1")
	{
		// 1. Authentication Endpoints
		auth := v1.Group("/auth")
		{
			// Public Authentication
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh-token", authHandler.RefreshToken)
			auth.POST("/forgot-password", authHandler.SendPasswordReset)

			// Protected Authentication (require valid token)
			authProtected := auth.Group("")
			authProtected.Use(middleware.AuthMiddleware(deps.UserRepository))
			{
				authProtected.POST("/logout", authHandler.Logout)
				authProtected.DELETE("/account", authHandler.DeleteAccount)
				authProtected.PUT("/password", authHandler.UpdatePassword)
				authProtected.GET("/me", authHandler.GetCurrentUser)
			}
		}

		// 2. User Profile Endpoints (Protected)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware(deps.UserRepository))
		{
			// Profile Management
			users.GET("/:userId", userHandler.GetProfile)
			users.PUT("/:userId", userHandler.UpdateProfile)
			users.PATCH("/:userId/:field", userHandler.UpdateProfileField)
			users.DELETE("/:userId/profile-image", userHandler.DeleteProfileImage)

			// User Addresses
			users.GET("/:userId/addresses", addressHandler.GetUserAddresses)
			users.POST("/:userId/addresses", addressHandler.SaveAddress)
		}

		// 3. Restaurant Endpoints
		restaurants := v1.Group("/restaurants")
		{
			// Restaurant Data
			restaurants.GET("", restaurantHandler.GetAllRestaurants)
			restaurants.GET("/:id", restaurantHandler.GetRestaurantByID)
			restaurants.GET("/popular", restaurantHandler.GetPopularRestaurants)
			restaurants.GET("/nearby", restaurantHandler.GetNearbyRestaurants)
			restaurants.GET("/search", restaurantHandler.SearchRestaurants)
			restaurants.GET("/category/:category", restaurantHandler.GetRestaurantsByCategory)
			restaurants.GET("/:id/menu", restaurantHandler.GetRestaurantMenu)
		}

		// 4. Food/Menu Endpoints
		foods := v1.Group("/foods")
		{
			// Food Data
			foods.GET("", foodHandler.GetAllFoods)
			foods.GET("/:id", foodHandler.GetFoodByID)
			foods.GET("/popular", foodHandler.GetPopularFoods)
			foods.GET("/recommended", foodHandler.GetRecommendedFoods)
			foods.GET("/category/:category", foodHandler.GetFoodsByCategory)
			foods.GET("/restaurant/:restaurantId", foodHandler.GetFoodsByRestaurant)
			foods.GET("/search", foodHandler.SearchFoods)
		}

		// 5. Order Endpoints (Protected)
		orders := v1.Group("/orders")
		orders.Use(middleware.AuthMiddleware(deps.UserRepository))
		{
			// Order Management
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/user/:userId", orderHandler.GetUserOrders)
			orders.GET("/:orderId", orderHandler.GetOrderByID)
			orders.PUT("/:orderId/status", orderHandler.UpdateOrderStatus)
			orders.DELETE("/:orderId", orderHandler.CancelOrder)
			orders.GET("/:orderId/track", orderHandler.TrackOrder)
		}

		// 6. WebSocket Endpoints
		// User profile updates (real-time)
		v1.GET("/users/:userId/watch", wsHandler.WatchUserProfile)

		// WebSocket statistics
		v1.GET("/websocket/stats", wsHandler.GetConnectionStats)

	}

	return router
}
