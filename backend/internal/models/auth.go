package models

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email      string    `gorm:"size:255;not null;uniqueIndex" json:"email"`
	Username   string    `gorm:"size:100;not null;uniqueIndex" json:"username"`
	Password   string    `gorm:"size:255;not null" json:"-"`
	FullName   string    `gorm:"size:255" json:"full_name"`
	AvatarURL  string    `gorm:"size:500" json:"avatar_url"`
	IsActive   bool      `gorm:"default:true;not null" json:"is_active"`
	IsVerified bool      `gorm:"default:false;not null" json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Documents     []ScannedDocument `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	// RefreshTokens []RefreshToken    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Storage *UserStorage `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

type UserStorage struct {
	UserID        uint      `gorm:"primaryKey" json:"user_id"`
	TotalBytes    int64     `gorm:"default:0;not null" json:"total_bytes"`
	DocumentCount int       `gorm:"default:0;not null" json:"document_count"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Request /Respone DTO

type RegisterRequest struct {
	Email    string `json:"email"  binding:"required"`
	Username string `json:"username"  binding:"required,min=8,max=32,alphanum"`
	Password string `json:"password"  binding:"required,min=8,max=32"`
	FullName string `json:"full_name" `
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	FullName  string `json:"full_name"`
	AvatarURL string `json:"avatar_url"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=128"`
}

type UserResponse struct {
	ID         uint      `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	FullName   string    `json:"full_name"`
	AvatarURL  string    `json:"avatar_url"`
	IsActive   bool      `json:"is_active"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

func ToUserResponse(u *User) UserResponse {
	return UserResponse{
		ID:         u.ID,
		Email:      u.Email,
		Username:   u.Username,
		FullName:   u.FullName,
		AvatarURL:  u.AvatarURL,
		IsActive:   u.IsActive,
		IsVerified: u.IsVerified,
		CreatedAt:  u.CreatedAt,
	}
}

type StorageResponse struct {
	TotalBytes    int64   `json:"total_bytes"`
	DocumentCount int     `json:"document_count"`
	TotalMB       float64 `json:"total_mb"`
}
