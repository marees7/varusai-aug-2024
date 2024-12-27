package repositories

import (
	"errors"
	"net/http"
	"shopping-site/pkg/models"
	"shopping-site/utils/dto"

	"gorm.io/gorm"
)

type IAuthRepository interface {
	LoginUser(dto.LoginRequest) (*models.Users, *dto.ErrorResponse)
	SignUpUser(models.Users) *dto.ErrorResponse
}

type authRepository struct {
	*gorm.DB
}

func CommenceAuthRepository(db *gorm.DB) IAuthRepository {
	return &authRepository{db}
}

func (db *authRepository) SignUpUser(user models.Users) *dto.ErrorResponse {
	record := db.Where("email=?", user.Email).First(&user)
	if record.RowsAffected == 0 {
		record = db.Create(&user)
		if record.Error != nil {
			return &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: record.Error.Error()}
		}

		return nil
	} else if record.RowsAffected == 1 {
		return &dto.ErrorResponse{Status: http.StatusConflict, Error: "user already exists"}
	} else {
		return &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: record.Error.Error()}
	}
}

func (db *authRepository) LoginUser(loginRequest dto.LoginRequest) (*models.Users, *dto.ErrorResponse) {
	var user models.Users

	record := db.Where("email=?", loginRequest.Email).First(&user)
	if errors.Is(record.Error, gorm.ErrRecordNotFound) {
		return nil, &dto.ErrorResponse{Status: http.StatusNotFound, Error: "user not found"}
	} else if record.Error != nil {
		return nil, &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: record.Error.Error()}
	}

	return &user, nil
}
