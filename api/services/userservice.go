package services

import (
	"net/http"
	"shopping-site/api/repositories"
	"shopping-site/api/validation"
	"shopping-site/pkg/loggers"
	"shopping-site/pkg/models"
	"shopping-site/utils/dto"
	"shopping-site/utils/helper"

	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	UpdateUser(uuid.UUID, *models.Users) *dto.ErrorResponse
	PlaceOrder(uuid.UUID, models.Orders) (*models.Orders, *dto.ErrorResponse)
	UpdateOrder(uuid.UUID, string) *dto.ErrorResponse
	GetOrders(uuid.UUID, map[string]interface{}) (*[]models.Orders, int64, *dto.ErrorResponse)
	GetOrder(uuid.UUID, string) (*models.Orders, *dto.ErrorResponse)
	GetProducts(map[string]interface{}) (*[]models.Products, int64, *dto.ErrorResponse)
	GetProduct(string) (*models.Products, *dto.ErrorResponse)
}

type userService struct {
	repositories.IUserRepository
}

func CommenceUserService(user repositories.IUserRepository) IUserService {
	return &userService{user}
}

func (repo *userService) PlaceOrder(userId uuid.UUID, order models.Orders) (*models.Orders, *dto.ErrorResponse) {
	return repo.IUserRepository.PlaceOrder(userId, order)
}

func (repo *userService) GetOrders(userId uuid.UUID, filters map[string]interface{}) (*[]models.Orders, int64, *dto.ErrorResponse) {
	return repo.IUserRepository.GetOrders(userId, filters)
}

func (repo *userService) GetProducts(filters map[string]interface{}) (*[]models.Products, int64, *dto.ErrorResponse) {
	return repo.IUserRepository.GetProducts(filters)
}

func (repo *userService) GetOrder(userId uuid.UUID, id string) (*models.Orders, *dto.ErrorResponse) {
	orderId, err := helper.PasreUuid(id)
	if err != nil {
		return nil, &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: err.Error()}
	}

	return repo.IUserRepository.GetOrder(userId, orderId)
}

func (repo *userService) GetProduct(id string) (*models.Products, *dto.ErrorResponse) {
	productId, err := uuid.Parse(id)
	if err != nil {
		return nil, &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: err.Error()}
	}
	return repo.IUserRepository.GetProduct(productId)
}

func (repo *userService) UpdateOrder(userId uuid.UUID, id string) *dto.ErrorResponse {
	orderId, err := uuid.Parse(id)
	if err != nil {
		return &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: err.Error()}
	}
	return repo.IUserRepository.UpdateOrder(userId, orderId)
}

func (repo *userService) UpdateUser(userId uuid.UUID, user *models.Users) *dto.ErrorResponse {
	user.UserID = userId

	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Phone == "" || user.Password == "" {
		loggers.WarnLog.Println("Required fields should not be empty")
		return &dto.ErrorResponse{
			Status: fiber.StatusBadRequest,
			Error:  "Required fields should not be empty"}
	}

	if err := validation.ValidateUser(*user); err != nil {
		loggers.WarnLog.Println(err.Error())
		return &dto.ErrorResponse{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		}
	}

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		loggers.ErrorLog.Println("Password hasing error")
		return &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: "Password hasing error"}
	}

	user.Password = string(hashedPin)

	if len(user.Address) != 0 {

		for _, data := range user.Address {
			if data.AddressID == uuid.Nil || data.DoorNo == "" || data.Street == "" || data.City == "" || data.State == "" || data.ZipCode == 0 {
				loggers.WarnLog.Println("Required Address fields should not be empty")
				return &dto.ErrorResponse{
					Status: fiber.StatusBadRequest,
					Error:  "Required Address fields should not be empty"}
			}
		}
	}
	return repo.IUserRepository.UpdateUser(user)
}
