package models

import "time"

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"size:500;not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	Revoked   bool      `gorm:"default:false;not null"`
	CreatedAt time.Time
}

// ── Auth DTOs ─────────────────────────────────────────────────

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
