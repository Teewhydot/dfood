package repository

import (
	"dfood/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id string) (*models.User, error)
	EmailExists(email string) (bool, error)
	UpdatePassword(email, hashedPassword string) error
	Update(id string, updates map[string]interface{}) error
	UpdateField(id, field string, value interface{}) error
	UpdateFCMToken(id, token string) error
}

type RestaurantRepository interface {
	GetAll(limit, offset int) ([]models.Restaurant, error)
	GetByID(id string) (*models.Restaurant, error)
	GetPopular(limit int) ([]models.Restaurant, error)
	GetNearby(latitude, longitude, radius float64, limit int) ([]models.Restaurant, error)
	Search(query string, limit, offset int) ([]models.Restaurant, error)
	GetByCategory(category string, limit, offset int) ([]models.Restaurant, error)
}

type FoodRepository interface {
	GetAll(limit, offset int) ([]models.Food, error)
	GetByID(id string) (*models.Food, error)
	GetPopular(limit int) ([]models.Food, error)
	GetByCategory(category string, limit, offset int) ([]models.Food, error)
	GetByRestaurant(restaurantID string, limit, offset int) ([]models.Food, error)
	Search(query string, limit, offset int) ([]models.Food, error)
}

type OrderRepository interface {
	Create(order *models.Order) error
	GetByID(id string) (*models.Order, error)
	GetByUserID(userID string, limit, offset int) ([]models.Order, error)
	UpdateStatus(id string, status models.OrderStatus) error
	Delete(id string) error
}

type PaymentRepository interface {
	GetPaymentMethods() ([]models.PaymentMethod, error)
	GetUserCards(userID string) ([]models.Card, error)
	CreateCard(card *models.Card) error
	DeleteCard(id string) error
	CreateTransaction(transaction *models.PaymentTransaction) error
	GetTransactionByID(id string) (*models.PaymentTransaction, error)
	UpdateTransaction(id string, updates map[string]interface{}) error
}

type AddressRepository interface {
	GetByUserID(userID string) ([]models.Address, error)
	Create(address *models.Address) error
	Update(id string, updates map[string]interface{}) error
	Delete(id string) error
	GetByID(id string) (*models.Address, error)
	GetDefaultByUserID(userID string) (*models.Address, error)
	SetDefault(userID, addressID string) error
}

type FavoritesRepository interface {
	GetFavoriteFoods(userID string) ([]models.Food, error)
	GetFavoriteRestaurants(userID string) ([]models.Restaurant, error)
	AddFavoriteFood(userID, foodID string) error
	RemoveFavoriteFood(userID, foodID string) error
	AddFavoriteRestaurant(userID, restaurantID string) error
	RemoveFavoriteRestaurant(userID, restaurantID string) error
	IsFoodFavorite(userID, foodID string) (bool, error)
	IsRestaurantFavorite(userID, restaurantID string) (bool, error)
	ClearAllFavorites(userID string) error
	GetFavoritesStats(userID string) (map[string]int, error)
}

type ChatRepository interface {
	GetByUserID(userID string) ([]models.Chat, error)
	GetByID(id string) (*models.Chat, error)
	Create(chat *models.Chat) error
	UpdateLastMessage(id, message string) error
	GetMessages(chatID string, limit, offset int) ([]models.Message, error)
	CreateMessage(message *models.Message) error
	MarkMessageAsRead(messageID string) error
	DeleteMessage(messageID string) error
	GetOrCreateChat(senderID, receiverID string, orderID *string) (*models.Chat, error)
}

type NotificationRepository interface {
	GetByUserID(userID string, limit, offset int) ([]models.Notification, error)
	Create(notification *models.Notification) error
	MarkAsRead(id string) error
	Delete(id string) error
}
