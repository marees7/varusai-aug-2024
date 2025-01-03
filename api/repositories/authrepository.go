package repositories

import (
	"errors"
	"net/http"
	"shopping-site/pkg/models"
	"shopping-site/utils/dto"

	"gorm.io/gorm"
)

type IAuthRepository interface {
	Login(dto.LoginRequest) (*models.Users, *dto.ErrorResponse)
	SignUp(models.Users) *dto.ErrorResponse
}

type authRepository struct {
	*gorm.DB
}

func CommenceAuthRepository(db *gorm.DB) IAuthRepository {
	return &authRepository{db}
}

// signup new member
func (db *authRepository) SignUp(user models.Users) *dto.ErrorResponse {
	var userExist models.Users
	// check user already avilable
	record := db.Where("email=?", user.Email).First(&userExist)
	if record.RowsAffected == 0 {
		// create new user
		record = db.Create(&user)
		if record.Error != nil {
			return &dto.ErrorResponse{
				Status: http.StatusBadRequest,
				Error:  record.Error.Error()}
		}
		return nil
	} else if record.RowsAffected == 1 {
		return &dto.ErrorResponse{
			Status: http.StatusConflict,
			Error:  userExist.Role + " already exists"}
	} else {
		return &dto.ErrorResponse{
			Status: http.StatusNotModified,
			Error:  record.Error.Error()}
	}
}

// login new member
func (db *authRepository) Login(loginRequest dto.LoginRequest) (*models.Users, *dto.ErrorResponse) {
	var user models.Users

	// check the user avilable and get his records
	record := db.Where("email=?", loginRequest.Email).First(&user)
	if errors.Is(record.Error, gorm.ErrRecordNotFound) {
		return nil, &dto.ErrorResponse{
			Status: http.StatusNotFound,
			Error:  "email does not exist"}
	} else if record.Error != nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  record.Error.Error()}
	}

	return &user, nil
}
