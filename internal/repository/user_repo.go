package repository

import (
	"database/sql"
	"net/http"
	"strconv"

	"dfood/internal/database"
	"dfood/internal/models"
	"dfood/pkg/errors"
)

type userRepository struct {
	db *sql.DB
}

// UpdatePassword implements UserRepository.
func (r *userRepository) UpdatePassword(email string, hashedPassword string) error {
	query := `
			UPDATE users
    		SET password = ?
  			WHERE email = ?
	`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "Failed to prepare update user statement", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(hashedPassword, email)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "Failed to update user password", err)
	}
	return nil
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: database.DB,
	}
}

func (r *userRepository) Create(user *models.User) error {
	query := `
	INSERT INTO users (first_name, last_name, email, password, id)
	VALUES (?, ?, ?, ?, ?);
	`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "Failed to prepare save_user statement", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Password, user.Id)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "Failed to register user", err)
	}
	return nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, first_name, last_name, email, password FROM users WHERE email = ?`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewHTTPError(http.StatusNotFound, "User not found", err)
		}
		return nil, errors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user by email", err)
	}
	return &user, nil
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	userIdInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid user ID format", err)
	}

	query := `SELECT id, first_name, last_name, email FROM users WHERE id = ?`

	var user models.User
	err = r.db.QueryRow(query, userIdInt).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewHTTPError(http.StatusNotFound, "User not found", err)
		}
		return nil, errors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user details", err)
	}
	return &user, nil
}

func (r *userRepository) EmailExists(email string) (bool, error) {
	query := `SELECT id FROM users WHERE email = ?`

	var id int64
	err := r.db.QueryRow(query, email).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.NewHTTPError(http.StatusInternalServerError, "Failed to check if user exists", err)
	}
	return true, nil
}
