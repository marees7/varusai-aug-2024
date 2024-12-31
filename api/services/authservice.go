package services

import (
	"net/http"
	"shopping-site/api/repositories"
	"shopping-site/pkg/loggers"
	"shopping-site/pkg/models"
	"shopping-site/utils/dto"

	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	SignUp(models.Users) *dto.ErrorResponse
	Login(dto.LoginRequest) (*models.Users, *dto.ErrorResponse)
}

type authService struct {
	repositories.IAuthRepository
}

func CommenceAuthService(auth repositories.IAuthRepository) IAuthService {
	return &authService{auth}
}

func (authRepo *authService) SignUp(user models.Users) *dto.ErrorResponse {
	hashedPin, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		loggers.ErrorLog.Println("Password generation error")
		return &dto.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Password generation error"}
	}

	user.Password = string(hashedPin)
	if err := authRepo.IAuthRepository.SignUp(user); err != nil {
		return err
	}

	return nil
}

func (authRepo *authService) Login(loginRequest dto.LoginRequest) (*models.Users, *dto.ErrorResponse) {
	user, err := authRepo.IAuthRepository.Login(loginRequest)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "invalid password"}
	}

	return user, nil
}
