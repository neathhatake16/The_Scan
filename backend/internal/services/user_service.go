package services

import (
	"github.com/neathhatake/the_Scan/internal/models"
	"github.com/neathhatake/the_Scan/internal/repositories"
	"github.com/neathhatake/the_Scan/pkg/apperrors"
)



type UserService struct {
	repo repositories.UserRepository
	auth *AuthService
}


func NewUserService(repo repositories.UserRepository, auth *AuthService) *UserService {
	return &UserService{
		repo: repo,
		auth: auth,
	}
}



func (s *UserService) GetProfile(userID uint) (*models.UserResponse,error){
	user , err := s.repo.FindByID(userID)
	if err != nil {
		return nil , err
	}

	res := models.ToUserResponse(user)
	return &res , nil

}


type UpdateProfile struct{
	FullName string
	AvatarURL string
}



func (s *UserService) UpdateProfile(userID uint , input UpdateProfile) (*models.UserResponse, error){
	user , err := s.repo.FindByID(userID)
	if err != nil {
		return nil , err
	}
	field := map[string]any{}

	if input.FullName != "" {
		field["full_name"] = input.FullName
	}
	if input.AvatarURL != "" {
		field["avatar_url"] = input.AvatarURL
	}
	if len(field) > 0 {
		if err := s.repo.UpdateUser(user, field); err != nil {
			return nil, err
		}
	}
	updated, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	res := models.ToUserResponse(updated)
	return &res, nil

} 


func (s *UserService) ChangePassword(userID uint , currentPW , NewPW string) error {
	user , err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}
	if !s.auth.CheckPassword(currentPW, user.Password){
		return apperrors.BadRequest("curent password is incorrect!")
	};
	hash ,err := s.auth.HashPassword(NewPW)
	if err != nil {
		return err
	}

	return s.repo.UpdateUser(user, map[string]any{"password": hash})
}



func (s *UserService) GetStorage(userID uint) (*models.StorageResponse, error){
	st , err := s.repo.FindStorage(userID)
	if err != nil {
		return nil , err
	}
	return &models.StorageResponse{
		TotalBytes:    st.TotalBytes,
		DocumentCount: st.DocumentCount,
		TotalMB:      float64(st.TotalBytes) / 1024 / 1024,
	}, nil
}
	






