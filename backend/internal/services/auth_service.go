package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/neathhatake/the_Scan/internal/config"
	"github.com/neathhatake/the_Scan/internal/models"
	"github.com/neathhatake/the_Scan/internal/repositories"
	"github.com/neathhatake/the_Scan/pkg/apperrors"
	"golang.org/x/crypto/bcrypt"
)



type AuthService struct {
     userRepo repositories.UserRepository
	 tokenRepo repositories.TokenRepository
	 cfg *config.Config
}



func NewAuthService(userRepo repositories.UserRepository, tokenRepo repositories.TokenRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		tokenRepo: tokenRepo,
		cfg: cfg,
	}
}



//Password 

func (s *AuthService) HashPassword(plain string) (string, error) {
	b , err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", apperrors.Internal("failed hash password", err)
	}
	return string(b),nil

}


func (s *AuthService) CheckPassword(plain , hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil 
}


//JWT 
type Claims struct {
	UserID uint `json:"user_id"`
	Type string `json:"type"`
	jwt.RegisteredClaims
}
func(s *AuthService) CreateAccessToken(userID uint) (string,error){
    claims := Claims{
		UserID: userID,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.cfg.AccessTTLMin) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", apperrors.Internal("failed to sign token", err)
	}
	return token, nil
}

func (s *AuthService) ParseAccessToken(tokenStr string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return 0, apperrors.Unauthorized("invalid or expired token")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || claims.Type != "access" {
		return 0, apperrors.Unauthorized("invalid token type")
	}
	return claims.UserID, nil
}



//Refresh Token 

// ── Refresh tokens ────────────────────────────────────────────

func (s *AuthService) CreateRefreshToken(userID uint) (string, error) {
	b := make([]byte, 48)
	if _, err := rand.Read(b); err != nil {
		return "", apperrors.Internal("failed to generate token", err)
	}
	tokenStr := base64.URLEncoding.EncodeToString(b)

	rt := &models.RefreshToken{
		UserID:    userID,
		Token:     tokenStr,
		ExpiresAt: time.Now().AddDate(0, 0, s.cfg.RefreshTTLDay),
	}
	if err := s.tokenRepo.CreateToken(rt); err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (s *AuthService) RotateRefreshToken(oldToken string) (uint, string, error) {
	rt, err := s.tokenRepo.FindActive(oldToken)
	if err != nil {
		return 0, "", err
	}
	if time.Now().After(rt.ExpiresAt) {
		return 0, "", apperrors.Unauthorized("refresh token expired")
	}
	if err := s.tokenRepo.Revoke(oldToken); err != nil {
		return 0, "", apperrors.Internal("failed to revoke token", err)
	}
	newToken, err := s.CreateRefreshToken(rt.UserID)
	if err != nil {
		return 0, "", err
	}
	return rt.UserID, newToken, nil
}

func (s *AuthService) RevokeRefreshToken(token string) {
	_ = s.tokenRepo.Revoke(token)
}

// ── Business logic ────────────────────────────────────────────

type RegisterInput struct {
	Email    string
	Username string
	Password string
	FullName string
}

func (s *AuthService) Register(input RegisterInput) (*models.User, error) {
	if _, err := s.userRepo.FindByEmail(input.Email); err == nil {
		return nil, apperrors.Conflict("email already registered")
	}
	if _, err := s.userRepo.FindByUsername(input.Username); err == nil {
		return nil, apperrors.Conflict("username already taken")
	}
	hash, err := s.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Email:    input.Email,
		Username: input.Username,
		Password: hash,
		FullName: input.FullName,
	}
	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}
	if err := s.userRepo.CreateStorage(&models.UserStorage{UserID: user.ID}); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, apperrors.Unauthorized("invalid email or password")
	}
	if !s.CheckPassword(password, user.Password) {
		return nil, apperrors.Unauthorized("invalid email or password")
	}
	if !user.IsActive {
		return nil, apperrors.Forbidden("account deactivated")
	}
	return user, nil
}