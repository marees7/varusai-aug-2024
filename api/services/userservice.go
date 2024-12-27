package services

import (
	"net/http"
	"shopping-site/api/repositories"
	"shopping-site/api/validation"
	"shopping-site/pkg/loggers"
	"shopping-site/pkg/models"
	"shopping-site/utils/dto"

	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	UpdateUserService(uuid.UUID, *models.Users) *dto.ErrorResponse
	PlaceOrderService(uuid.UUID, models.Orders) (*models.Orders, *dto.ErrorResponse)
	CancelOrderService(uuid.UUID, string) *dto.ErrorResponse
	GetOrdersService(uuid.UUID) (*[]models.Orders, *dto.ErrorResponse)
	GetProductsService(map[string]string, uuid.UUID) (*[]models.Products, *dto.ErrorResponse)
	GetProductService(uuid.UUID, string) (*models.Products, *dto.ErrorResponse)
	FilterProductsService(map[string]string) (*[]models.Products, *dto.ErrorResponse)
}

type userService struct {
	repositories.IUserRepository
}

func CommenceUserService(user repositories.IUserRepository) IUserService {
	return &userService{user}
}

func (repo *userService) UpdateUserService(userIdCtx uuid.UUID, user *models.Users) *dto.ErrorResponse {
	user.UserId = userIdCtx

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
			if data.AddressId == uuid.Nil || data.DoorNo == "" || data.Street == "" || data.City == "" || data.State == "" || data.ZipCode == 0 {
				loggers.WarnLog.Println("Required Address fields should not be empty")
				return &dto.ErrorResponse{
					Status: fiber.StatusBadRequest,
					Error:  "Required Address fields should not be empty"}
			}
		}
	}
	return repo.UpdateUserRepository(user)
}

func (repo *userService) PlaceOrderService(userIdCtx uuid.UUID, order models.Orders) (*models.Orders, *dto.ErrorResponse) {
	UserId := userIdCtx

	return repo.PlaceOrderRepository(UserId, order)
}

func (repo *userService) CancelOrderService(userIdCtx uuid.UUID, id string) *dto.ErrorResponse {
	UserId := userIdCtx

	orderId, err := uuid.Parse(id)
	if err != nil {
		return &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: err.Error()}
	}
	return repo.CancelOrderRepository(UserId, orderId)
}

func (repo *userService) GetOrdersService(userIdCtx uuid.UUID) (*[]models.Orders, *dto.ErrorResponse) {
	UserId := userIdCtx

	return repo.GetOrdersRepository(UserId)
}

func (repo *userService) GetProductsService(filters map[string]string, userIdCtx uuid.UUID) (*[]models.Products, *dto.ErrorResponse) {
	UserId := userIdCtx

	return repo.GetProductsRepository(filters, UserId)
}

func (repo *userService) GetProductService(userIdCtx uuid.UUID, id string) (*models.Products, *dto.ErrorResponse) {
	UserId := userIdCtx

	productId, err := uuid.Parse(id)
	if err != nil {
		return nil, &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: err.Error()}
	}
	return repo.GetProductRepository(UserId, productId)
}

func (repo *userService) FilterProductsService(filters map[string]string) (*[]models.Products, *dto.ErrorResponse) {
	return repo.FilterProductsRepository(filters)
}
