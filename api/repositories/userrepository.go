package repositories

import (
	"errors"
	"shopping-site/pkg/loggers"
	"shopping-site/pkg/models"
	"shopping-site/utils/constants"
	"shopping-site/utils/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	UpdateUser(*models.Users) *dto.ErrorResponse
	PlaceOrder(uuid.UUID, models.Orders) (*models.Orders, *dto.ErrorResponse)
	UpdateOrder(uuid.UUID, uuid.UUID) *dto.ErrorResponse
	GetOrders(uuid.UUID, map[string]interface{}) (*[]models.Orders, int64, *dto.ErrorResponse)
	GetOrder(uuid.UUID, uuid.UUID) (*models.Orders, *dto.ErrorResponse)
	GetProducts(map[string]interface{}) (*[]models.Products, int64, *dto.ErrorResponse)
	GetProduct(uuid.UUID) (*models.Products, *dto.ErrorResponse)
}

type userRepository struct {
	*gorm.DB
}

func CommenceUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (db *userRepository) PlaceOrder(userId uuid.UUID, order models.Orders) (*models.Orders, *dto.ErrorResponse) {
	var (
		userDetails    models.Users
		orderItems     []models.OrderedItems
		addressDetails models.Addresses
		totalAmount    float64
	)

	record := db.Where("address_id= ? AND user_id= ?", order.AddressID, userId).First(&addressDetails)
	if record.Error != nil {
		loggers.WarnLog.Println("specified address not avilable on user profile")
		return nil, &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: "specified address not avilable on user profile"}
	}

	for _, item := range order.Products {
		var productDetails models.Products

		record = db.Where("product_id= ?", item.ProductID).First(&productDetails)
		if record.Error != nil {
			loggers.ErrorLog.Println("error while getting product details")
			return nil, &dto.ErrorResponse{Status: fiber.StatusBadRequest,
				Error: record.Error.Error()}
		}
		itemAmount := productDetails.Price * float64(item.Quantity)
		totalAmount += itemAmount

		item.MerchantId = productDetails.UserID
		item.ProductName = productDetails.ProductName
		item.Price = productDetails.Price
		item.AddressID = order.AddressID
		item.CustomerID = userId

		orderItems = append(orderItems, item)
	}

	record = db.Where("user_id= ?", userId).First(&userDetails)
	if record.Error != nil {
		return nil, &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: record.Error.Error()}
	}

	order.UserID = userId
	order.Name = userDetails.FirstName + " " + userDetails.LastName
	order.Email = userDetails.Email
	order.Phone = userDetails.Phone
	order.Status = constants.Inprogress
	order.TotalAmount = totalAmount
	order.Products = orderItems

	record = db.Create(&order)
	if record.Error != nil {
		return nil, &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: record.Error.Error()}
	}

	return &order, nil
}

func (db *userRepository) GetOrders(userId uuid.UUID, filters map[string]interface{}) (*[]models.Orders, int64, *dto.ErrorResponse) {
	var (
		orders []models.Orders
		count  int64
	)

	fromDate := filters["from_date"].(string)
	toDate := filters["to_date"].(string)
	limit := filters["limit"].(int)
	offset := filters["offset"].(int)

	record := db.Preload("Products").Where("user_id = ?", userId).Limit(limit).Offset(offset).Find(&orders).Count(&count)
	if errors.Is(record.Error, gorm.ErrRecordNotFound) {
		return nil, 0, &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "no orders avilable"}
	}

	if fromDate != "" && toDate == "" {
		record.Where("created_at >= ?", fromDate).Find(&orders).Count(&count)
	} else if fromDate == "" && toDate != "" {
		record.Where("created_at <= ?", toDate).Find(&orders).Count(&count)
	} else if fromDate != "" && toDate != "" {
		if fromDate > toDate {
			{
				return nil, 0, &dto.ErrorResponse{Status: fiber.StatusBadRequest,
					Error: "from date exceeds to date"}
			}
		}

		record.Where("created_at BETWEEN ? AND ?", fromDate, toDate).Find(&orders).Count(&count)
	}

	return &orders, count, nil
}

func (db *userRepository) GetOrder(userId uuid.UUID, orderId uuid.UUID) (*models.Orders, *dto.ErrorResponse) {
	var order models.Orders

	record := db.Preload("Products").Where("order_id= ? AND user_id= ?", orderId, userId).Find(&order)
	if errors.Is(record.Error, gorm.ErrRecordNotFound) {
		return nil, &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "order not found on your list"}
	}

	return &order, nil
}

func (db *userRepository) GetProducts(filter map[string]interface{}) (*[]models.Products, int64, *dto.ErrorResponse) {
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

	record := db.Raw(`SELECT * FROM getProducts_fn($1,$2,$3,$4,$5,$6,$7)`, nil, brandName, categoryName, price, rating, limit, offset).Find(&products)
	if record.Error != nil {
		return nil, 0, &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: record.Error.Error()}
	}

	record = db.Raw(`SELECT count (*) FROM getProducts_fn($1,$2,$3,$4,$5,$6,$7)`, nil, brandName, categoryName, price, rating, limit, offset).Find(&count)
	if record.Error != nil {
		return nil, 0, &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: record.Error.Error()}
	}

	return &products, count, nil
}

func (db *userRepository) GetProduct(productId uuid.UUID) (*models.Products, *dto.ErrorResponse) {
	var product models.Products

	record := db.Where("product_id = ?", productId).First(&product)
	if record.Error != nil {
		return nil, &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: record.Error.Error()}
	}

	return &product, nil
}

func (db *userRepository) UpdateUser(user *models.Users) *dto.ErrorResponse {
	var userExcist models.Users

	record := db.Where("user_id = ? ", user.UserID).First(&userExcist)
	if errors.Is(record.Error, gorm.ErrRecordNotFound) {
		return &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "user not found"}
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
			Error: "user updation failed"}
	}

	for _, data := range user.Address {
		var addressExcist models.Addresses

		record = db.Where("address_id = ? AND user_id = ?", data.AddressID, user.UserID).First(&addressExcist)
		if errors.Is(record.Error, gorm.ErrRecordNotFound) {
			return &dto.ErrorResponse{Status: fiber.StatusBadRequest,
				Error: "specified address not avilable"}
		}

		record = db.Where("address_id = ? AND user_id = ?", data.AddressID, user.UserID).Updates(models.Addresses{
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
				Error: "address updation failed"}
		}
	}

	return nil
}

func (db *userRepository) UpdateOrder(userId uuid.UUID, orderId uuid.UUID) *dto.ErrorResponse {
	var order models.Orders
	record := db.Where("order_id = ? AND user_id= ? ", orderId, userId).First(&order)
	if record.RowsAffected == 0 {
		return &dto.ErrorResponse{Status: fiber.StatusNotFound,
			Error: "order not avilable"}
	}

	record = db.Model(&order).Where("order_id = ?", orderId).Update("status", constants.Cancelled)
	if record.RowsAffected == 0 {
		return &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: record.Error.Error()}
	}

	return nil
}
