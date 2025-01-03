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

type IMerchantService interface {
	CreateProduct(uuid.UUID, *models.Products) *dto.ErrorResponse
	GetCategories() (*[]models.Categories, *dto.ErrorResponse)
	GetBrands() (*[]models.Brands, *dto.ErrorResponse)
	GetProduct(uuid.UUID, string) (*models.Products, *dto.ErrorResponse)
	GetProducts(map[string]interface{}, uuid.UUID) (*[]models.Products, int64, *dto.ErrorResponse)
	GetOrders(uuid.UUID, map[string]interface{}) (*[]models.OrderedItems, int64, *dto.ErrorResponse)
	GetOrder(uuid.UUID, string) (*[]models.OrderedItems, *dto.ErrorResponse)
	UpdateProduct(uuid.UUID, *models.Products) *dto.ErrorResponse
	UpdateOrderStatus(uuid.UUID, string, string) *dto.ErrorResponse
	UpdateMerchant(uuid.UUID, *models.Users) *dto.ErrorResponse
	DeleteProduct(string) *dto.ErrorResponse
}

type merchantService struct {
	repositories.IMerchantRepository
}

func CommenceMerchantService(merchant repositories.IMerchantRepository) IMerchantService {
	return &merchantService{merchant}
}

// create products
func (repo *merchantService) CreateProduct(userId uuid.UUID, product *models.Products) *dto.ErrorResponse {
	product.UserID = userId
	// call create product repository
	return repo.IMerchantRepository.CreateProduct(product)
}

// get the avilable categories
func (repo *merchantService) GetCategories() (*[]models.Categories, *dto.ErrorResponse) {
	// call get categories repository
	return repo.IMerchantRepository.GetCategories()
}

// get the avilable brands
func (repo *merchantService) GetBrands() (*[]models.Brands, *dto.ErrorResponse) {
	// call get brand repository
	return repo.IMerchantRepository.GetBrands()
}

// get all the products for the merchant
func (repo *merchantService) GetProducts(filters map[string]interface{}, userId uuid.UUID) (*[]models.Products, int64, *dto.ErrorResponse) {
	// call get products repository
	return repo.IMerchantRepository.GetProducts(filters, userId)
}

// get a single product by id
func (repo *merchantService) GetProduct(userId uuid.UUID, id string) (*models.Products, *dto.ErrorResponse) {
	// parsing the product_id
	productId, err := helper.PasreUuid(id)
	if err != nil {
		return nil, &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: err.Error()}
	}

	// call get product repository
	return repo.IMerchantRepository.GetProduct(userId, productId)
}

// get all the orders for the merchant
func (repo *merchantService) GetOrders(userId uuid.UUID, filters map[string]interface{}) (*[]models.OrderedItems, int64, *dto.ErrorResponse) {
	// call get orders repository
	return repo.IMerchantRepository.GetOrders(userId, filters)
}

// get a single order by id
func (repo *merchantService) GetOrder(userId uuid.UUID, id string) (*[]models.OrderedItems, *dto.ErrorResponse) {
	// parsing order_id
	orderId, err := helper.PasreUuid(id)
	if err != nil {
		return nil, &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: err.Error()}
	}

	// call get order repository
	return repo.IMerchantRepository.GetOrder(userId, orderId)
}

// update product details
func (repo *merchantService) UpdateProduct(userId uuid.UUID, product *models.Products) *dto.ErrorResponse {
	product.UserID = userId

	if (product.ProductID) == uuid.Nil || product.ProductName == "" || product.Price == 0 {
		loggers.WarnLog.Println("product update required fields should not be empty")
		return &dto.ErrorResponse{
			Status: fiber.StatusBadRequest,
			Error:  "product update required fields should not be empty"}
	}

	// call update product repository
	return repo.IMerchantRepository.UpdateProduct(product)
}

// update order status
func (repo *merchantService) UpdateOrderStatus(userId uuid.UUID, id string, orderStatus string) *dto.ErrorResponse {
	if orderStatus != constants.Shipped {
		loggers.WarnLog.Println("insufficient permission to update specific status")
		return &dto.ErrorResponse{
			Status: fiber.StatusForbidden,
			Error:  "insufficient permission to update specific status"}
	}
	// parsing order_id
	orderId, err := helper.PasreUuid(id)
	if err != nil {
		loggers.ErrorLog.Println(err)
		return &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: err.Error()}
	}

	//call update order status  repository
	return repo.IMerchantRepository.UpdateOrderStatus(orderId, userId, orderStatus)
}

// update merchant details
func (repo *merchantService) UpdateMerchant(userId uuid.UUID, user *models.Users) *dto.ErrorResponse {
	user.UserID = userId

	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Phone == "" || user.Password == "" {
		loggers.WarnLog.Println("merchant update required fields should not be empty")
		return &dto.ErrorResponse{
			Status: fiber.StatusBadRequest,
			Error:  "merchant update required fields should not be empty"}
	}
	// validate merchant details
	if err := validation.ValidateUser(*user); err != nil {
		loggers.WarnLog.Println(err.Error())
		return &dto.ErrorResponse{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		}
	}

	//generate hash pass for updated password
	hashedPin, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		loggers.ErrorLog.Println("Password hasing error")
		return &dto.ErrorResponse{Status: http.StatusBadRequest, Error: "Password hasing error"}
	}

	user.Password = string(hashedPin)

	if len(user.Address) != 0 {
		for _, data := range user.Address {
			if data.AddressID == uuid.Nil || data.DoorNo == "" || data.Street == "" || data.City == "" || data.State == "" || data.ZipCode == 0 {
				loggers.WarnLog.Println("merchant required Address fields should not be empty")
				return &dto.ErrorResponse{
					Status: fiber.StatusBadRequest,
					Error:  "merchant required address fields should not be empty"}
			}
		}
	}

	// call update merchant repository
	return repo.IMerchantRepository.UpdateMerchant(user)
}

// delete product of the merchant
func (repo *merchantService) DeleteProduct(id string) *dto.ErrorResponse {
	productId, err := helper.PasreUuid(id)
	if err != nil {
		return &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: err.Error()}
	}

	// call delete product repository
	return repo.IMerchantRepository.DeleteProduct(productId)
}
