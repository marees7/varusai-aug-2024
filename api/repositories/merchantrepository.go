package repositories

import (
	"shopping-site/pkg/models"
	"shopping-site/utils/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IMerchantRepository interface {
	AddProductRepository(*models.Products) *dto.ErrorResponse
	RemoveProductRepository(uuid.UUID) *dto.ErrorResponse
	UpdateProductRepository(*models.Products) *dto.ErrorResponse
	UpdateMerchantRepository(*models.Users) *dto.ErrorResponse
	UpdateOrderStatusRepository(uuid.UUID, uuid.UUID, string) *dto.ErrorResponse
	GetProductsRepository(map[string]string, uuid.UUID) (*[]models.Products, *dto.ErrorResponse)
	GetProductRepository(uuid.UUID, uuid.UUID) (*models.Products, *dto.ErrorResponse)
	GetOrdersRepository(uuid.UUID) (*models.Orders, *dto.ErrorResponse)
}

type merchantRepository struct {
	*gorm.DB
}

func CommenceMerchantRepository(db *gorm.DB) IMerchantRepository {
	return &merchantRepository{db}
}

func (db *merchantRepository) AddProductRepository(product *models.Products) *dto.ErrorResponse {
	record := db.Where("product_name = ? AND user_id = ?", product.ProductName, product.UserId).First(product)
	if record.RowsAffected > 0 {
		return &dto.ErrorResponse{Status: fiber.StatusConflict,
			Error: "product already exists on your listing"}
	}

	record = db.Create(product)
	if record.Error != nil {
		return &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: record.Error.Error()}
	}

	return nil
}

func (db *merchantRepository) RemoveProductRepository(productId uuid.UUID) *dto.ErrorResponse {
	var product models.Products

	record := db.Where("product_id = ?", productId).First(&product)
	if record.RowsAffected == 0 {
		return &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "product does not exists"}
	}

	record = db.Where("product_id = ?", product.ProductId).Delete(&product)
	if record.Error != nil {
		return &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: record.Error.Error()}
	} else if record.RowsAffected == 0 {
		return &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: "something went wrong"}
	}

	return nil
}

func (db *merchantRepository) UpdateProductRepository(product *models.Products) *dto.ErrorResponse {
	var productExcist models.Products

	record := db.Where("product_id = ? AND user_id = ? ", product.ProductId, product.UserId).First(&productExcist)
	if record.RowsAffected == 0 {
		return &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "product not found on your listing"}
	}

	record = db.Where("product_id = ?", product.ProductId).Updates(models.Products{ProductName: product.ProductName, Price: product.Price})
	if record.Error != nil {
		return &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: record.Error.Error()}
	} else if record.RowsAffected == 0 {
		return &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: "something went wrong"}
	}

	return nil
}

func (db *merchantRepository) UpdateOrderStatusRepository(orderId uuid.UUID, userId uuid.UUID, orderStatus string) *dto.ErrorResponse {
	var orderExcist models.Orders

	record := db.Where("order_id= ? AND user_id", orderId, userId).First(&orderExcist)
	if record.RowsAffected == 0 {
		return &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "product not found on your listing"}
	}

	record = db.Where("order_id = ?", orderId).Updates(models.Orders{Status: orderStatus})
	if record.Error != nil {
		return &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: record.Error.Error()}
	} else if record.RowsAffected == 0 {
		return &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: "something went wrong"}
	}

	return nil
}

func (db *merchantRepository) UpdateMerchantRepository(user *models.Users) *dto.ErrorResponse {
	var userExcist models.Users

	record := db.Where("user_id = ? ", user.UserId).First(&userExcist)
	if record.RowsAffected == 0 {
		return &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "user not found"}
	}

	record = db.Where("user_id = ?", user.UserId).Updates(models.Users{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Password:  user.Password})
	if record.Error != nil {
		return &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: record.Error.Error()}
	} else if record.RowsAffected == 0 {
		return &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: "something went wrong"}
	}

	for _, data := range user.Address {
		record = db.Where("address_id = ?", data.AddressId).Updates(models.Addresses{
			DoorNo:  data.DoorNo,
			Street:  data.Street,
			City:    data.City,
			State:   data.State,
			ZipCode: data.ZipCode})
		if record.Error != nil {
			return &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
				Error: record.Error.Error()}
		} else if record.RowsAffected == 0 {
			return &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
				Error: "something went wrong"}
		}
	}

	return nil
}

func (db *merchantRepository) GetProductsRepository(filter map[string]string, userId uuid.UUID) (*[]models.Products, *dto.ErrorResponse) {
	var products []models.Products

	categoryName := filter["category_name"]
	brandName := filter["brand_name"]

	record := db.Raw(`SELECT * FROM getProducts_fn($1,$2,$3)`, userId, brandName, categoryName).Find(&products)
	if record.Error != nil {
		return nil, &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: record.Error.Error()}
	}

	return &products, nil
}

func (db *merchantRepository) GetProductRepository(userId uuid.UUID, productId uuid.UUID) (*models.Products, *dto.ErrorResponse) {
	var product models.Products

	record := db.Where("product_id= ? AND user_id= ?", productId, userId).First(&product)
	if record.RowsAffected == 0 {
		return nil, &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "product not found on your listing"}
	}

	return &product, nil
}

func (db *merchantRepository) GetOrdersRepository(userId uuid.UUID) (*models.Orders, *dto.ErrorResponse) {
	var orders models.Orders

	record := db.Where("user_id= ?", userId).First(&orders)
	if record.RowsAffected == 0 {
		return nil, &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "no orders avilable"}
	}

	return &orders, nil
}
