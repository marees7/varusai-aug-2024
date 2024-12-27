package services

import (
	"net/http"
	"shopping-site/api/repositories"
	"shopping-site/api/validation"
	"shopping-site/pkg/loggers"
	"shopping-site/pkg/models"
	"shopping-site/utils/constants"
	"shopping-site/utils/dto"

	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IMerchantService interface {
	AddProductService(uuid.UUID, *models.Products) *dto.ErrorResponse
	RemoveProductService(string) *dto.ErrorResponse
	UpdateProductService(uuid.UUID, *models.Products) *dto.ErrorResponse
	UpdateMerchantService(uuid.UUID, *models.Users) *dto.ErrorResponse
	GetProductsService(map[string]string, uuid.UUID) (*[]models.Products, *dto.ErrorResponse)
	UpdateOrderStatusService(uuid.UUID, string, string) *dto.ErrorResponse
	GetProductService(uuid.UUID, string) (*models.Products, *dto.ErrorResponse)
	GetOrdersService(uuid.UUID) (*models.Orders, *dto.ErrorResponse)
}

type merchantService struct {
	repositories.IMerchantRepository
}

func CommenceMerchantService(merchant repositories.IMerchantRepository) IMerchantService {
	return &merchantService{merchant}
}

func (repo *merchantService) AddProductService(userIdCtx uuid.UUID, product *models.Products) *dto.ErrorResponse {
	product.UserId = userIdCtx

	return repo.AddProductRepository(product)
}

func (repo *merchantService) RemoveProductService(id string) *dto.ErrorResponse {
	productId, err := uuid.Parse(id)
	if err != nil {
		loggers.ErrorLog.Println(err)
		return &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: err.Error()}
	}
	return repo.RemoveProductRepository(productId)
}

func (repo *merchantService) UpdateProductService(userIdCtx uuid.UUID, product *models.Products) *dto.ErrorResponse {
	product.UserId = userIdCtx

	if (product.ProductId) == uuid.Nil || product.ProductName == "" || product.Price == 0 {
		loggers.WarnLog.Println("Required fields should not be empty")
		return &dto.ErrorResponse{
			Status: fiber.StatusBadRequest,
			Error:  "Required fields should not be empty"}
	}

	return repo.UpdateProductRepository(product)
}

func (repo *merchantService) UpdateMerchantService(userIdCtx uuid.UUID, user *models.Users) *dto.ErrorResponse {
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
	return repo.UpdateMerchantRepository(user)
}

func (repo *merchantService) UpdateOrderStatusService(userIdCtx uuid.UUID, id string, orderStatus string) *dto.ErrorResponse {
	if orderStatus != constants.Shipped {
		loggers.WarnLog.Println("insufficient permission to update specific status")
		return &dto.ErrorResponse{
			Status: fiber.StatusForbidden,
			Error:  "insufficient permission to update specific status"}
	}

	userId := userIdCtx
	orderId, err := uuid.Parse(id)
	if err != nil {
		loggers.ErrorLog.Println(err)
		return &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: err.Error()}
	}

	return repo.UpdateOrderStatusRepository(orderId, userId, orderStatus)
}

func (repo *merchantService) GetProductsService(filters map[string]string, userIdCtx uuid.UUID) (*[]models.Products, *dto.ErrorResponse) {
	UserId := userIdCtx

	return repo.GetProductsRepository(filters, UserId)
}

func (repo *merchantService) GetProductService(userIdCtx uuid.UUID, id string) (*models.Products, *dto.ErrorResponse) {
	UserId := userIdCtx

	productId, err := uuid.Parse(id)
	if err != nil {
		loggers.ErrorLog.Println(err)
		return nil, &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: err.Error()}
	}

	return repo.GetProductRepository(UserId, productId)
}

func (repo *merchantService) GetOrdersService(userIdCtx uuid.UUID) (*models.Orders, *dto.ErrorResponse) {
	UserId := userIdCtx

	return repo.GetOrdersRepository(UserId)
}
