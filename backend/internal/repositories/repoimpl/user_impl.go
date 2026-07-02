package repoimpl

import (
	"github.com/neathhatake/the_Scan/internal/models"
	"github.com/neathhatake/the_Scan/internal/repositories"
	"github.com/neathhatake/the_Scan/pkg/apperrors"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepo{
		db: db,
	}
}


func (r *userRepo)CreateUser(user *models.User) error {
	if err := r.db.Create(user).Error ; err != nil {
		return apperrors.Internal("failed to create user...",err)
	}
	return nil
}

func (r *userRepo) FindByID(id uint) (*models.User, error){
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil , apperrors.NotFound("failed to find user...")
	}
	return &user , nil
}


func (r *userRepo) FindByEmail(email string) (*models.User, error){
	var user models.User
	if err := r.db.Where("email = ?",email).First(&user).Error; err != nil {
		return nil , apperrors.NotFound("failed to find user...")
	}
	return &user , nil
}

func (r *userRepo) FindByUsername(username string)(*models.User,error){
	var user models.User
	if err := r.db.Where("username = ?",username).First(&user).Error; err != nil {
		return nil , apperrors.NotFound("failed to find user...")
	}
	return &user , nil
}


func (r *userRepo) UpdateUser(user *models.User, feilds map[string]any) error {
	if err := r.db.Model(user).Updates(feilds).Error; err != nil {
		return apperrors.Internal("failed to update user...",err)
	}
	return nil

}

func (r *userRepo) CreateStorage(stor *models.UserStorage) error {
	if err := r.db.Create(stor).Error; err != nil {
		return apperrors.Internal("failed to create storage...",err)
	}
	return nil
}	


func (r *userRepo) FindStorage(userID uint) (*models.UserStorage, error) {
	var s models.UserStorage
	if err := r.db.Where("user_id = ?", userID).First(&s).Error; err != nil {
		return nil, apperrors.NotFound("storage info not found")
	}
	return &s, nil
}

func (r *userRepo) IncrementStorage(userID uint, bytes int64) error {
	if err := r.db.Model(&models.UserStorage{}).Where("user_id = ?", userID).Updates(map[string]any{
		"total_bytes":    gorm.Expr("total_bytes + ?", bytes),
		"document_count": gorm.Expr("document_count + 1"),
	}).Error; err != nil {
		return apperrors.Internal("failed to increment storage", err)
	}
	return nil
}

func (r *userRepo) DecrementStorage(userID uint, bytes int64) error {
	if err := r.db.Model(&models.UserStorage{}).Where("user_id = ?", userID).Updates(map[string]any{
		"total_bytes":    gorm.Expr("GREATEST(0, total_bytes - ?)", bytes),
		"document_count": gorm.Expr("GREATEST(0, document_count - 1)"),
	}).Error; err != nil {
		return apperrors.Internal("failed to decrement storage", err)
	}
	return nil
}

