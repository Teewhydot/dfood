package repository

import (
	"dfood/internal/database"

	"gorm.io/gorm"
)

type emailRepository struct {
	db *gorm.DB
}

// SendEmail implements EmailRepository.
func (e *emailRepository) SendEmail(to string, subject string, body string) error {
	panic("unimplemented")
}

func NewEmailRepository() EmailRepository {
	return &emailRepository{
		db: database.DB,
	}
}
