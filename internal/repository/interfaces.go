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

type EventRepository interface {
	Create(event *models.Event) error
	GetAll() ([]models.Event, error)
	GetByID(id string) (*models.Event, error)
	Update(id string, event *models.Event) (*models.Event, error)
	Delete(id string) error
}
