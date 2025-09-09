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
}
