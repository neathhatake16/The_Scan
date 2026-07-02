package repositories

import (
	"github.com/neathhatake/the_Scan/internal/models"
)

type TokenRepository interface {
	CreateToken(rt *models.RefreshToken) error
	FindActive(token string) (*models.RefreshToken, error)
	Revoke(token string) error
	RevokeAllForUser(userID uint) error
}



