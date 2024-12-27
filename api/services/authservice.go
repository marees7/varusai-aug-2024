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
	SignUpService(models.Users) *dto.ErrorResponse
	LoginService(dto.LoginRequest) (*models.Users, *dto.ErrorResponse)
}

type authService struct {
	repositories.IAuthRepository
}

func CommenceAuthService(auth repositories.IAuthRepository) IAuthService {
	return &authService{auth}
}

func (authRepo *authService) SignUpService(user models.Users) *dto.ErrorResponse {
	hashedPin, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		loggers.ErrorLog.Println("Password generation error")
		return &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: "Password generation error"}
	}

	user.Password = string(hashedPin)
	if err := authRepo.IAuthRepository.SignUpUser(user); err != nil {
		return err
	}

	return nil
}

func (authRepo *authService) LoginService(loginRequest dto.LoginRequest) (*models.Users, *dto.ErrorResponse) {
	user, err := authRepo.IAuthRepository.LoginUser(loginRequest)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return nil, &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: "invalid password"}
	}

	return user, nil
}
