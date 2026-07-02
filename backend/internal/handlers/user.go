package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neathhatake/the_Scan/internal/middleware"
	"github.com/neathhatake/the_Scan/internal/services"
	"github.com/neathhatake/the_Scan/pkg/respone"
)


type UserHandler struct {
	usrHandler *services.UserService
}

func NewUserHandler(usrHandler *services.UserService) *UserHandler {
	return &UserHandler{
		usrHandler: usrHandler,
	}
}



// GET/users/me
// @Summary     Get current user profile
// @Description Returns authenticated user profile information
// @Tags        User
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} respone.Envelop{data=models.UserResponse} "User profile"
// @Failure     401 {object} respone.Envelop{error=string} "Unauthorized"
// @Failure     500 {object} respone.Envelop{error=string} "Internal server error"
// @Router      /users/me [get]
func (h *UserHandler) GetProfile(c *gin.Context){
	profile , err := h.usrHandler.GetProfile(middleware.UserID(c))
	if err != nil {
		respone.Error(c, err)
		return
	}

	respone.OK(c, profile)
		
}


// PATCH/users/me


func (h *UserHandler) UpdateProfile(c *gin.Context){
	var req struct {
		FullName string `json:"full_name"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {

		respone.Error(c,respone.ValidationError(err))
		return
	}

	profile , err := h.usrHandler.UpdateProfile(middleware.UserID(c), services.UpdateProfile{
		FullName: req.FullName,
		AvatarURL: req.AvatarURL,
	})
	if err != nil {
		respone.Error(c, err)
		return
	}
	respone.OK(c, profile)

}

//POST /users/me/change-password


func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req struct {
		 CurrentPassword string `json:"current_password" binding:"required"`
		 NewPassword     string `json:"new_password" binding:"required,min=8,max=128"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respone.Error(c,respone.ValidationError(err))
		return
	}
	if err := h.usrHandler.ChangePassword(middleware.UserID(c), req.CurrentPassword, req.NewPassword); err != nil {
		respone.Error(c, err)
		return
	}
	respone.OK(c, "password update successfully")
}

// GET/users/me/storage

func (h *UserHandler) GetStorage(c *gin.Context){
	storage , err := h.usrHandler.GetStorage(middleware.UserID(c))

	if err != nil {
		respone.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, respone.Envelop{
		Success: true,
		Data: storage,
	})
}



