package repositories

import "github.com/neathhatake/the_Scan/internal/models"

type UserRepository interface {
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User, fields map[string]any) error

	CreateStorage(s *models.UserStorage) error
	FindStorage(userID uint) (*models.UserStorage, error)
	IncrementStorage(userID uint, bytes int64) error
	DecrementStorage(userID uint, bytes int64) error
}
