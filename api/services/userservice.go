package services

import (
	"net/http"
	"shopping-site/api/repositories"
	"shopping-site/api/validation"
	"shopping-site/pkg/loggers"
	"shopping-site/pkg/models"
	"shopping-site/utils/constants"
	"shopping-site/utils/dto"
	"shopping-site/utils/helper"

	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	CreateOrder(uuid.UUID, models.Orders) (*models.Orders, *dto.ErrorResponse)
	GetOrders(uuid.UUID, map[string]interface{}) (*[]models.Orders, int64, *dto.ErrorResponse)
	GetOrder(uuid.UUID, string) (*models.Orders, *dto.ErrorResponse)
	GetProducts(map[string]interface{}) (*[]models.Products, int64, *dto.ErrorResponse)
	GetProduct(string) (*models.Products, *dto.ErrorResponse)
	UpdateUser(uuid.UUID, *models.Users) *dto.ErrorResponse
	UpdateOrder(uuid.UUID, string, *models.Orders) *dto.ErrorResponse
}

type userService struct {
	repositories.IUserRepository
}

func CommenceUserService(user repositories.IUserRepository) IUserService {
	return &userService{user}
}

// create new order
func (repo *userService) CreateOrder(userId uuid.UUID, order models.Orders) (*models.Orders, *dto.ErrorResponse) {
	// call create order repository
	return repo.IUserRepository.CreateOrder(userId, order)
}

// get products based on the filters
func (repo *userService) GetProducts(filters map[string]interface{}) (*[]models.Products, int64, *dto.ErrorResponse) {
	// call get products repository
	return repo.IUserRepository.GetProducts(filters)
}

// get a single product by id
func (repo *userService) GetProduct(id string) (*models.Products, *dto.ErrorResponse) {
	// parse the product_id
	productId, err := uuid.Parse(id)
	if err != nil {
		return nil, &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: err.Error()}
	}

	// call get product repository
	return repo.IUserRepository.GetProduct(productId)
}

// get the user's orders based on provided filters
func (repo *userService) GetOrders(userId uuid.UUID, filters map[string]interface{}) (*[]models.Orders, int64, *dto.ErrorResponse) {
	// call get orders repository
	return repo.IUserRepository.GetOrders(userId, filters)
}

// get a single order by id
func (repo *userService) GetOrder(userId uuid.UUID, id string) (*models.Orders, *dto.ErrorResponse) {
	// parse the order_id
	orderId, err := helper.PasreUuid(id)
	if err != nil {
		return nil, &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: err.Error()}
	}

	// call get product repository
	return repo.IUserRepository.GetOrder(userId, orderId)
}

// cancel order of user
func (repo *userService) UpdateOrder(userId uuid.UUID, id string, order *models.Orders) *dto.ErrorResponse {
	if order.Status != constants.Cancelled {
		loggers.WarnLog.Println("insufficient permission to update specific status")
		return &dto.ErrorResponse{
			Status: fiber.StatusForbidden,
			Error:  "insufficient permission to update specific status"}
	}

	// parse the order_id
	orderId, err := uuid.Parse(id)
	if err != nil {
		return &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: err.Error()}
	}

	// call update products repository
	return repo.IUserRepository.UpdateOrder(userId, orderId)
}

// update user details
func (repo *userService) UpdateUser(userId uuid.UUID, user *models.Users) *dto.ErrorResponse {
	user.UserID = userId

	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Phone == "" || user.Password == "" {
		loggers.WarnLog.Println("Required fields should not be empty")
		return &dto.ErrorResponse{
			Status: fiber.StatusBadRequest,
			Error:  "Required fields should not be empty"}
	}

	// validate user details
	if err := validation.ValidateUser(*user); err != nil {
		loggers.WarnLog.Println(err.Error())
		return &dto.ErrorResponse{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		}
	}

	// generate hashed pass for updated password
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

	// call update user repository
	return repo.IUserRepository.UpdateUser(user)
}
