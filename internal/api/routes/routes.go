package routes

import (
	"time"

	"dfood/internal/api/handlers"
	"dfood/internal/api/middleware"
	"dfood/internal/service"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	AuthService         service.AuthService
	UserService         service.UserService
	RestaurantService   service.RestaurantService
	FoodService         service.FoodService
	OrderService        service.OrderService
	PaymentService      service.PaymentService
	AddressService      service.AddressService
	FavoritesService    service.FavoritesService
	ChatService         service.ChatService
	NotificationService service.NotificationService
	UploadService       service.UploadService
}

func SetupRoutes(deps *Dependencies) *gin.Engine {
	router := gin.New()

	// Global Middleware
	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())
	router.Use(middleware.RateLimitMiddleware(10, time.Minute)) // 10 requests per minute per IP

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(deps.AuthService)
	userHandler := handlers.NewUserHandler(deps.UserService)
	restaurantHandler := handlers.NewRestaurantHandler(deps.RestaurantService)
	foodHandler := handlers.NewFoodHandler(deps.FoodService)
	orderHandler := handlers.NewOrderHandler(deps.OrderService)
	paymentHandler := handlers.NewPaymentHandler(deps.PaymentService)
	addressHandler := handlers.NewAddressHandler(deps.AddressService)
	favoritesHandler := handlers.NewFavoritesHandler(deps.FavoritesService)
	chatHandler := handlers.NewChatHandler(deps.ChatService)
	notificationHandler := handlers.NewNotificationHandler(deps.NotificationService)
	uploadHandler := handlers.NewUploadHandler(deps.UploadService)

	// API v1 Routes
	v1 := router.Group("/api/v1")
	{
		// 1. Authentication Endpoints
		auth := v1.Group("/auth")
		{
			// User Authentication
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.DELETE("/delete-account", authHandler.DeleteAccount)

			// Email Management
			auth.POST("/send-password-reset", authHandler.SendPasswordReset)
			auth.POST("/send-email-verification", authHandler.SendEmailVerification)
			auth.GET("/verify-email-status", authHandler.VerifyEmailStatus)
			auth.GET("/current-user", authHandler.GetCurrentUser)
			auth.POST("/password/update", authHandler.UpdatePassword)
		}

		// 2. User Profile Endpoints
		users := v1.Group("/users")
		{
			// Profile Management
			users.GET("/:userId", userHandler.GetProfile)
			users.PUT("/:userId", userHandler.UpdateProfile)
			users.PATCH("/:userId/:field", userHandler.UpdateProfileField)
			users.POST("/:userId/upload-image", userHandler.UploadProfileImage)
			users.DELETE("/:userId/profile-image", userHandler.DeleteProfileImage)
			users.GET("/:userId/stream", userHandler.GetProfileStream)
			users.POST("/:userId/sync-profile", userHandler.SyncProfile)

			// User Addresses
			users.GET("/:userId/addresses", addressHandler.GetUserAddresses)
			users.POST("/:userId/addresses", addressHandler.SaveAddress)
			users.PUT("/:userId/addresses/:addressId", addressHandler.UpdateAddress)
			users.DELETE("/:userId/addresses/:addressId", addressHandler.DeleteAddress)
			users.GET("/:userId/addresses/default", addressHandler.GetDefaultAddress)
			users.PUT("/:userId/addresses/:addressId/set-default", addressHandler.SetDefaultAddress)
			users.GET("/:userId/addresses/stream", addressHandler.GetAddressStream)

			// Favorites Management
			users.GET("/:userId/favorites/foods", favoritesHandler.GetFavoriteFoods)
			users.GET("/:userId/favorites/restaurants", favoritesHandler.GetFavoriteRestaurants)
			users.POST("/:userId/favorites/foods/:foodId", favoritesHandler.AddFavoriteFood)
			users.DELETE("/:userId/favorites/foods/:foodId", favoritesHandler.RemoveFavoriteFood)
			users.POST("/:userId/favorites/restaurants/:restaurantId", favoritesHandler.AddFavoriteRestaurant)
			users.DELETE("/:userId/favorites/restaurants/:restaurantId", favoritesHandler.RemoveFavoriteRestaurant)
			users.GET("/:userId/favorites/foods/:foodId/status", favoritesHandler.CheckFoodFavoriteStatus)
			users.GET("/:userId/favorites/restaurants/:restaurantId/status", favoritesHandler.CheckRestaurantFavoriteStatus)
			users.POST("/:userId/favorites/foods/:foodId/toggle", favoritesHandler.ToggleFoodFavorite)
			users.POST("/:userId/favorites/restaurants/:restaurantId/toggle", favoritesHandler.ToggleRestaurantFavorite)
			users.DELETE("/:userId/favorites", favoritesHandler.ClearAllFavorites)
			users.GET("/:userId/favorites/stats", favoritesHandler.GetFavoritesStats)
			users.GET("/:userId/favorites/foods/stream", favoritesHandler.GetFavoriteFoodsStream)
			users.GET("/:userId/favorites/restaurants/stream", favoritesHandler.GetFavoriteRestaurantsStream)

			// User Chats
			users.GET("/:userId/chats", chatHandler.GetUserChats)
			users.GET("/:userId/chats/stream", chatHandler.GetChatsStream)

			// User Notifications
			users.GET("/:userId/notifications", notificationHandler.GetUserNotifications)
			users.GET("/:userId/notifications/stream", notificationHandler.GetNotificationsStream)
			users.POST("/:userId/fcm-token", notificationHandler.UpdateFCMToken)
			users.GET("/:userId/fcm-token", notificationHandler.GetFCMToken)
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

		// 5. Order Endpoints
		orders := v1.Group("/orders")
		{
			// Order Management
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/user/:userId", orderHandler.GetUserOrders)
			orders.GET("/:orderId", orderHandler.GetOrderByID)
			orders.PUT("/:orderId/status", orderHandler.UpdateOrderStatus)
			orders.DELETE("/:orderId", orderHandler.CancelOrder)
			orders.GET("/:orderId/track", orderHandler.TrackOrder)
		}

		// 6. Payment Endpoints
		payments := v1.Group("/payments")
		{
			// Payment Methods
			payments.GET("/methods", paymentHandler.GetPaymentMethods)
			payments.GET("/cards/:userId", paymentHandler.GetUserCards)
			payments.POST("/cards", paymentHandler.SaveCard)
			payments.DELETE("/cards/:cardId", paymentHandler.DeleteCard)

			// Payment Processing
			payments.POST("/process", paymentHandler.ProcessPayment)
			payments.GET("/transaction/:transactionId", paymentHandler.GetTransactionDetails)
			payments.POST("/refund", paymentHandler.ProcessRefund)
		}

		// 7. Chat/Messaging Endpoints
		chats := v1.Group("/chats")
		{
			// Chat Management
			chats.GET("/:chatId", chatHandler.GetChatDetails)
			chats.POST("", chatHandler.CreateOrGetChat)
			chats.PUT("/:chatId/last-message", chatHandler.UpdateLastMessage)
			chats.GET("/:chatId/messages", chatHandler.GetChatMessages)
			chats.POST("/:chatId/messages", chatHandler.SendMessage)
			chats.GET("/:chatId/messages/stream", chatHandler.GetMessagesStream)
			chats.GET("/:chatId/new-messages/stream", chatHandler.GetNewMessagesStream)
		}

		// Message Management
		messages := v1.Group("/messages")
		{
			messages.PUT("/:messageId/read", chatHandler.MarkMessageAsRead)
			messages.DELETE("/:messageId", chatHandler.DeleteMessage)
		}

		// 8. Notification Endpoints
		notifications := v1.Group("/notifications")
		{
			// Notification Management
			notifications.POST("", notificationHandler.SendNotification)
			notifications.PUT("/:notificationId/read", notificationHandler.MarkNotificationAsRead)
			notifications.DELETE("/:notificationId", notificationHandler.DeleteNotification)
		}

		// Push Notifications
		pushNotifications := v1.Group("/push-notifications")
		{
			pushNotifications.POST("/send", notificationHandler.SendPushNotification)
		}

		// 9. File Upload Endpoints
		upload := v1.Group("/upload")
		{
			// Image Management
			upload.POST("/profile-image", uploadHandler.UploadProfileImage)
			upload.POST("/food-image", uploadHandler.UploadFoodImage)
			upload.POST("/restaurant-image", uploadHandler.UploadRestaurantImage)
			upload.DELETE("/:imageId", uploadHandler.DeleteImage)
		}
	}

	return router
}
