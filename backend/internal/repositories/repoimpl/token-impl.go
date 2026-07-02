package repoimpl

import (
	"github.com/neathhatake/the_Scan/internal/models"
	"github.com/neathhatake/the_Scan/internal/repositories"
	"github.com/neathhatake/the_Scan/pkg/apperrors"
	"gorm.io/gorm"
)


type tokenRepo struct {
	db *gorm.DB
}


func NewTokenRepository(db *gorm.DB) repositories.TokenRepository {
	return &tokenRepo{
		db: db,
	}
}


func (r *tokenRepo)  CreateToken(rt *models.RefreshToken) error{
	if err := r.db.Create(rt).Error; err != nil {
		return apperrors.Internal("failed to save refresh token!",err)
	}
	return nil
}


func (r *tokenRepo) FindActive(token string) (*models.RefreshToken, error){
	var rt models.RefreshToken
	err := r.db.Where("token = ? AND revoked = false", token).First(&rt).Error
	if err != nil {
		return nil, apperrors.Unauthorized("invalid refresh token")
	}
	return &rt, nil
}


func (r *tokenRepo) Revoke(token string)error{
		return r.db.Model(&models.RefreshToken{}).
		Where("token = ?", token).
		Update("revoked", true).Error
}

func (r *tokenRepo) RevokeAllForUser(userID uint) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error
}
