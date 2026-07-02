package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/neathhatake/the_Scan/internal/models"
	"github.com/neathhatake/the_Scan/internal/services"
	"github.com/neathhatake/the_Scan/pkg/respone"
)

type AuthHandler struct {
	authService *services.AuthService
}


func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}



// POST/auth/register


func (h *AuthHandler) Register(c *gin.Context){
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respone.Error(c, respone.ValidationError(err))
		return 
	}

	user , err := h.authService.Register(services.RegisterInput{
		Email: req.Email,
		Username: req.Username,
		Password: req.Password,
		FullName: req.FullName,
	})
	if err != nil {
		respone.Error(c, err)
		return
	}
	token , err := h.issueToken(user.ID)
	if err != nil {
		respone.Error(c, err)
		return
	}

	respone.OK(c,token)



}

func (h *AuthHandler) issueToken(userID uint) (*models.TokenResponse, error){
	access , err := h.authService.CreateAccessToken(userID)
	if err != nil {
		return nil, err
	}
	refresh , err := h.authService.CreateRefreshToken(userID)
	if err != nil {
		return nil, err
	}
	return &models.TokenResponse{
		AccessToken: access,
		RefreshToken: refresh,
		TokenType: "bearer",
	}, nil 
}


// POST/auth/login


func (h *AuthHandler) Login(c *gin.Context){
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respone.Error(c, respone.ValidationError(err))
		return 
	}
	user , err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		respone.Error(c, err)
		return
	}

	token , err := h.issueToken(user.ID)
	if err != nil {
		respone.Error(c, err)
		return
	}
	respone.OK(c, token)
}

// POST/auth/logout


func (h *AuthHandler) Logout(c *gin.Context){
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respone.Error(c, respone.ValidationError(err))
		return
	}
	h.authService.RevokeRefreshToken(req.RefreshToken)
	respone.Message(c,"logout successfully")
	
}

// POST/auth/refresh


func (h *AuthHandler) Refresh(c *gin.Context){
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respone.Error(c,respone.ValidationError(err))
		return
	}
	userID , newRe , err := h.authService.RotateRefreshToken(req.RefreshToken)
	if err != nil {
		respone.Error(c, err)
		return
	}
	access , err := h.authService.CreateAccessToken(userID)
	if err != nil {
		respone.Error(c, err)
		return
	}
	respone.OK(c, models.TokenResponse{
		AccessToken: access,
		RefreshToken: newRe,
		TokenType: "bearer",
	})
}