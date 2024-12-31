package repositories

import (
	"errors"
	"shopping-site/pkg/models"
	"shopping-site/utils/constants"
	"shopping-site/utils/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IMerchantRepository interface {
	CreateProduct(*models.Products) *dto.ErrorResponse
	DeleteProduct(uuid.UUID) *dto.ErrorResponse
	UpdateProduct(*models.Products) *dto.ErrorResponse
	UpdateMerchant(*models.Users) *dto.ErrorResponse
	UpdateOrderStatus(uuid.UUID, uuid.UUID, string) *dto.ErrorResponse
	GetProducts(map[string]interface{}, uuid.UUID) (*[]models.Products, int64, *dto.ErrorResponse)
	GetProduct(uuid.UUID, uuid.UUID) (*models.Products, *dto.ErrorResponse)
	GetOrders(uuid.UUID, map[string]interface{}) (*[]models.OrderedItems, int64, *dto.ErrorResponse)
	GetOrder(uuid.UUID, uuid.UUID) (*[]models.OrderedItems, *dto.ErrorResponse)
}

type merchantRepository struct {
	*gorm.DB
}

func CommenceMerchantRepository(db *gorm.DB) IMerchantRepository {
	return &merchantRepository{db}
}

func (db *merchantRepository) CreateProduct(product *models.Products) *dto.ErrorResponse {
	record := db.Where("product_name = ? AND user_id = ?", product.ProductName, product.UserID).First(product)
	if record.RowsAffected > 0 {
		return &dto.ErrorResponse{Status: fiber.StatusConflict,
			Error: "product already exists on your listing"}
	}

	record = db.Create(product)
	if record.Error != nil {
		return &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: record.Error.Error()}
	}

	return nil
}

func (db *merchantRepository) DeleteProduct(productId uuid.UUID) *dto.ErrorResponse {
	var product models.Products

	record := db.Where("product_id = ?", productId).First(&product)
	if errors.Is(record.Error, gorm.ErrRecordNotFound) {
		return &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "product does not exists"}
	}

	record = db.Where("product_id = ?", product.ProductID).Delete(&product)
	if record.Error != nil {
		return &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: record.Error.Error()}
	} else if record.RowsAffected == 0 {
		return &dto.ErrorResponse{Status: fiber.StatusNotModified,
			Error: "product deletion failed"}
	}

	return nil
}

func (db *merchantRepository) UpdateProduct(product *models.Products) *dto.ErrorResponse {
	var productExcist models.Products

	record := db.Where("product_id = ? AND user_id = ? ", product.ProductID, product.UserID).First(&productExcist)
	if errors.Is(record.Error, gorm.ErrRecordNotFound) {
		return &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "product not found on your listing"}
	}

	record = db.Where("product_id = ?", product.ProductID).Updates(models.Products{ProductName: product.ProductName, Price: product.Price})
	if record.Error != nil {
		return &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: record.Error.Error()}
	} else if record.RowsAffected == 0 {
		return &dto.ErrorResponse{Status: fiber.StatusNotModified,
			Error: "product updation failed"}
	}

	return nil
}

func (db *merchantRepository) UpdateOrderStatus(orderItemId uuid.UUID, userId uuid.UUID, orderStatus string) *dto.ErrorResponse {
	var itemExcist models.OrderedItems

	record := db.Where("ordered_items_id= ? AND merchant_id = ?", orderItemId, userId).First(&itemExcist)
	if errors.Is(record.Error, gorm.ErrRecordNotFound) {
		return &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "order item not found on your listing"}
	}

	record = db.Model(itemExcist).Where("ordered_items_id = ?", orderItemId).Update("status", constants.Shipped)
	if record.Error != nil {
		return &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: record.Error.Error()}
	} else if record.RowsAffected == 0 {
		return &dto.ErrorResponse{Status: fiber.StatusNotModified,
			Error: "order_item updation failed"}
	}

	return nil
}

func (db *merchantRepository) UpdateMerchant(user *models.Users) *dto.ErrorResponse {
	var userExcist models.Users

	record := db.Where("user_id = ? ", user.UserID).First(&userExcist)
	if errors.Is(record.Error, gorm.ErrRecordNotFound) {
		return &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "merchant not found"}
	}

	record = db.Where("user_id = ?", user.UserID).Updates(models.Users{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Password:  user.Password})
	if record.Error != nil {
		return &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: record.Error.Error()}
	} else if record.RowsAffected == 0 {
		return &dto.ErrorResponse{Status: fiber.StatusNotModified,
			Error: "merchant updation failed"}
	}

	for _, data := range user.Address {
		var checkAddress models.Addresses

		record := db.Where("user_id = ? AND address_id = ?", user.UserID, data.AddressID).First(&checkAddress)
		if errors.Is(record.Error, gorm.ErrRecordNotFound) {
			return &dto.ErrorResponse{Status: fiber.StatusNotFound,
				Error: "specified merchant address not avilable"}
		}

		record = db.Where("address_id = ?", data.AddressID).Updates(models.Addresses{
			DoorNo:  data.DoorNo,
			Street:  data.Street,
			City:    data.City,
			State:   data.State,
			ZipCode: data.ZipCode})
		if record.Error != nil {
			return &dto.ErrorResponse{Status: fiber.StatusBadRequest,
				Error: record.Error.Error()}
		} else if record.RowsAffected == 0 {
			return &dto.ErrorResponse{Status: fiber.StatusNotModified,
				Error: "merchant address updation failed"}
		}
	}

	return nil
}

func (db *merchantRepository) GetProducts(filter map[string]interface{}, userId uuid.UUID) (*[]models.Products, int64, *dto.ErrorResponse) {
	var (
		products []models.Products
		count    int64
	)

	categoryName := filter["category_name"].(string)
	brandName := filter["brand_name"].(string)
	price := filter["price"].(float64)
	rating := filter["rating"].(float32)
	limit := filter["limit"].(int)
	offset := filter["offset"].(int)

	record := db.Raw(`SELECT * FROM getProducts_fn($1,$2,$3,$4,$5,$6,$7)`, userId, brandName, categoryName, price, rating, limit, offset).Find(&products)
	if record.Error != nil {
		return nil, 0, &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: record.Error.Error()}
	}

	record = db.Raw(`SELECT count (*) FROM getProducts_fn($1,$2,$3,$4,$5,$6,$7)`, userId, brandName, categoryName, price, rating, limit, offset).Find(&count)
	if record.Error != nil {
		return nil, 0, &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: record.Error.Error()}
	}

	return &products, count, nil
}

func (db *merchantRepository) GetProduct(userId uuid.UUID, productId uuid.UUID) (*models.Products, *dto.ErrorResponse) {
	var product models.Products

	record := db.Where("product_id= ? AND user_id= ?", productId, userId).First(&product)
	if record.RowsAffected == 0 {
		return nil, &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "product not found on your listing"}
	}

	return &product, nil
}

func (db *merchantRepository) GetOrders(userId uuid.UUID, filters map[string]interface{}) (*[]models.OrderedItems, int64, *dto.ErrorResponse) {
	var (
		orderItems []models.OrderedItems
		count      int64
	)

	fromDate := filters["from_date"].(string)
	toDate := filters["to_date"].(string)
	limit := filters["limit"].(int)
	offset := filters["offset"].(int)

	record := db.Where("merchant_id = ?", userId).Limit(limit).Offset(offset).Find(&orderItems).Count(&count)
	if errors.Is(record.Error, gorm.ErrRecordNotFound) {
		return nil, 0, &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "no orders avilable"}
	}

	if fromDate != "" && toDate == "" {
		record.Where("created_at >= ?", fromDate).Find(&orderItems).Count(&count)
	} else if fromDate == "" && toDate != "" {
		record.Where("created_at <= ?", toDate).Find(&orderItems).Count(&count)
	} else if fromDate != "" && toDate != "" {
		if fromDate > toDate {
			{
				return nil, 0, &dto.ErrorResponse{Status: fiber.StatusBadRequest,
					Error: "from date exceeds to date"}
			}
		}

		record.Where("created_at BETWEEN ? AND ?", fromDate, toDate).Find(&orderItems).Count(&count)
	}

	return &orderItems, count, nil
}

func (db *merchantRepository) GetOrder(userId uuid.UUID, orderId uuid.UUID) (*[]models.OrderedItems, *dto.ErrorResponse) {
	var orderItem []models.OrderedItems

	record := db.Where("order_id= ? AND merchant_id= ?", orderId, userId).Find(&orderItem)
	if errors.Is(record.Error, gorm.ErrRecordNotFound) {
		return nil, &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "order not found on your list"}
	}

	return &orderItem, nil
}
