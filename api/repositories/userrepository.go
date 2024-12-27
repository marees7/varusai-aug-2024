package repositories

import (
	"shopping-site/pkg/loggers"
	"shopping-site/pkg/models"
	"shopping-site/utils/constants"
	"shopping-site/utils/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	UpdateUserRepository(*models.Users) *dto.ErrorResponse
	PlaceOrderRepository(uuid.UUID, models.Orders) (*models.Orders, *dto.ErrorResponse)
	CancelOrderRepository(uuid.UUID, uuid.UUID) *dto.ErrorResponse
	GetOrdersRepository(uuid.UUID) (*[]models.Orders, *dto.ErrorResponse)
	GetProductsRepository(map[string]string, uuid.UUID) (*[]models.Products, *dto.ErrorResponse)
	GetProductRepository(uuid.UUID, uuid.UUID) (*models.Products, *dto.ErrorResponse)
	FilterProductsRepository(map[string]string) (*[]models.Products, *dto.ErrorResponse)
}

type userRepository struct {
	*gorm.DB
}

func CommenceUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (db *userRepository) PlaceOrderRepository(userId uuid.UUID, order models.Orders) (*models.Orders, *dto.ErrorResponse) {
	var (
		userDetails    models.Users
		orderItems     []models.OrderedItems
		addressDetails models.Addresses
		totalAmount    float64
	)

	record := db.Where("address_id= ? AND user_id= ?", order.AddressId, userId).First(&addressDetails)
	if record.Error != nil {
		loggers.WarnLog.Println("specified address not avilable on user profile")
		return nil, &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: record.Error.Error()}
	}

	for _, item := range order.Products {
		var productDetails models.Products

		record = db.Where("product_id= ?", item.ProductId).First(&productDetails)
		if record.Error != nil {
			loggers.ErrorLog.Println("error while getting product details")
			return nil, &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
				Error: record.Error.Error()}
		}
		itemAmount := productDetails.Price * float64(item.Quantity)
		totalAmount += itemAmount

		item.ProductName = productDetails.ProductName
		item.Price = productDetails.Price

		orderItems = append(orderItems, item)
	}

	record = db.Where("user_id= ?", userId).First(&userDetails)
	if record.Error != nil {
		return nil, &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: record.Error.Error()}
	}

	order.UserId = userId
	order.Name = userDetails.FirstName + " " + userDetails.LastName
	order.Email = userDetails.Email
	order.Phone = userDetails.Phone
	order.Status = constants.Placed
	order.TotalAmount = totalAmount
	order.Products = orderItems

	record = db.Create(&order)
	if record.Error != nil {
		return nil, &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: record.Error.Error()}
	}

	return &order, nil
}

func (db *userRepository) UpdateUserRepository(user *models.Users) *dto.ErrorResponse {
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

func (db *userRepository) CancelOrderRepository(userId uuid.UUID, orderId uuid.UUID) *dto.ErrorResponse {
	var order models.Orders
	record := db.Where("order_id = ? AND user_id= ? ", orderId, userId).First(&order)
	if record.RowsAffected == 0 {
		return &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "order not avilable"}
	}

	record = db.Model(&order).Where("order_id = ?", orderId).Update("status", constants.Cancelled)
	if record.RowsAffected == 0 {
		return &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: record.Error.Error()}
	}

	return nil
}

func (db *userRepository) GetOrdersRepository(userId uuid.UUID) (*[]models.Orders, *dto.ErrorResponse) {
	var orders []models.Orders

	record := db.Preload("Products").Where("user_id= ?", userId).Find(&orders)
	if record.Error != nil {
		return nil, &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: record.Error.Error()}
	}

	return &orders, nil
}

func (db *userRepository) GetProductsRepository(filter map[string]string, userId uuid.UUID) (*[]models.Products, *dto.ErrorResponse) {
	var products []models.Products

	categoryName := filter["category_name"]
	brandName := filter["brand_name"]

	record := db.Raw(`SELECT * FROM getProductsUser_fn($1,$2)`, brandName, categoryName).Find(&products)
	if record.Error != nil {
		return nil, &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: record.Error.Error()}
	}

	return &products, nil
}

func (db *userRepository) GetProductRepository(userId uuid.UUID, productId uuid.UUID) (*models.Products, *dto.ErrorResponse) {
	var product models.Products

	record := db.Where("product_id = ?", productId).First(&product)
	if record.Error != nil {
		return nil, &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: record.Error.Error()}
	}

	return &product, nil
}

func (db *userRepository) FilterProductsRepository(filter map[string]string) (*[]models.Products, *dto.ErrorResponse) {
	var products []models.Products

	price := filter["price"]
	rating := filter["rating"]
	if price == "" && rating == "" {
		return nil, &dto.ErrorResponse{
			Status: fiber.StatusBadRequest,
			Error:  "filter keywords should not be empty",
		}
	}

	record := db.Raw(`SELECT * FROM filterProductsUser_fn($1,$2)`, price, rating).Find(&products)
	if record.Error != nil {
		return nil, &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: record.Error.Error()}
	}

	return &products, nil
}
