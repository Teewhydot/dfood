package database

import (
	"dfood/internal/config"
	"dfood/internal/models"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase(cfg *config.Config) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.DB.Datasource), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("could not get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	// Auto migrate the schema
	if err = DB.AutoMigrate(
		&models.User{},
		&models.Address{},
		&models.Permission{},
		&models.Restaurant{},
		&models.RestaurantFoodCategory{},
		&models.Food{},
		&models.Order{},
		&models.PaymentMethod{},
		&models.Card{},
		&models.PaymentTransaction{},
		&models.Chat{},
		&models.Message{},
		&models.Notification{},
		&models.FavoriteFood{},
		&models.FavoriteRestaurant{},
		&models.RecentKeyword{},
	); err != nil {
		return fmt.Errorf("could not migrate database: %w", err)
	}

	return nil
}

func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
